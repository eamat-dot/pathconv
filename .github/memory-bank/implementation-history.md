# Implementation History: pathconv

## 実装の経緯と重要な判断

このドキュメントは、pathconv実装時の詳細な経緯、遭遇した問題、解決策、および将来の開発者が理解すべき重要な判断を記録しています。

## 実装フェーズ

### Phase 1: 基本実装（2025-12-27 初期）
**目標**: 単一のWindows形式入力から複数形式への変換

**実装内容**:
- 10種類の出力形式対応
- 相対パス・絶対パスの判定
- チルダ展開機能
- フラグ指定順での出力制御

**遭遇した課題**:
1. **Windows形式のスラッシュ問題**
   - 入力: `c:/test`
   - 期待: `c:\test` (バックスラッシュ)
   - 実際: `c:/test` (スラッシュそのまま)
   - 解決: `convertToWindows()` で常に `/` → `\` に変換

### Phase 2: 複数入力形式対応（2025-12-27 午後）
**目標**: 様々な入力形式を自動判定・変換

#### 設計判断: 正規化パイプライン

**判断**: すべての入力をWindows形式に正規化してから各形式に変換
```
入力 → normalizeInputPath() → Windows形式 → convertToXXX() → 出力
```

**理由**:
1. **単一責任**: 各変換関数はWindows形式のみを期待
2. **拡張性**: 新しい入力形式は `normalizeInputPath()` だけ修正
3. **テスト容易性**: 正規化と変換を独立してテスト可能
4. **保守性**: 変換ロジックが単純化

**代替案を採用しなかった理由**:
- ❌ 各変換関数で全入力形式をサポート → コード重複、保守困難
- ❌ 入力形式ごとに変換関数を分ける → 組み合わせ爆発 (6入力 × 10出力 = 60関数)

#### 実装の詳細

**1. 入力形式の自動判定**

最初の実装では `isRelative()` で相対パス判定後、絶対パス形式をチェックしていたが、問題発生：

```go
// ❌ 問題のあった実装
func normalizeInputPath(path string) string {
    if isRelative(path) {  // file:/// が相対と判定される！
        return strings.ReplaceAll(path, "/", "\\")
    }
    
    if strings.HasPrefix(path, "file:///") {
        // ここに到達しない
    }
}
```

**問題**: `file:///C:/hoge` が相対パスと判定され、`file:\\\C:\hoge` という破損したパスになる

**原因**: `isRelative()` で file:// をチェックしていなかった

**解決**:
```go
// ✅ 修正後
func normalizeInputPath(path string) string {
    // file:// を最初にチェック（相対判定より前）
    if strings.HasPrefix(path, "file:///") {
        urlPath := path[8:]
        return normalizeWindowsPath(urlPath)
    }
    
    // UNC拡張長パスも早めにチェック
    if strings.HasPrefix(path, "\\\\?\\") {
        return normalizeWindowsPath(path[4:])
    }
    
    // その後、相対パス判定
    if isRelative(path) {
        return strings.ReplaceAll(path, "/", "\\")
    }
    
    // 残りの絶対パス形式
    // ...
}
```

**教訓**: 特殊形式のチェックは相対パス判定より前に行う

**2. ドライブレター大文字化**

**問題**: `c:/hoge` と `C:/hoge` が異なる出力になる
- テストで `c:/hoge` → 期待 `c:\hoge`、実際 `C:\hoge`

**判断**: ドライブレターは常に大文字に統一

**理由**:
- Windows APIは大文字を返す (`GetFullPathName` など)
- 大文字が慣例的
- 一貫性の向上

**実装**:
```go
func normalizeWindowsPath(path string) string {
    path = strings.ReplaceAll(path, "/", "\\")
    
    if len(path) >= 2 && path[1] == ':' {
        path = strings.ToUpper(string(path[0])) + path[1:]
    }
    
    return path
}
```

**3. GitBash形式の修正**

**問題**: GitBash形式が `/c/hoge` ではなく `C:/hoge` を出力

**原因**: `convertToGitBash()` がドライブレター変換を実装していなかった

**修正**:
```go
func convertToGitBash(path string, home string, relative bool) string {
    result := path
    if home != "" && strings.HasPrefix(path, home) {
        result = "$USERPROFILE" + path[len(home):]
    }
    result = strings.ReplaceAll(result, "\\", "/")
    
    // 追加: ドライブレター変換
    if !relative {
        re := regexp.MustCompile(`^([A-Za-z]):`)
        result = re.ReplaceAllStringFunc(result, func(match string) string {
            drive := strings.ToLower(string(match[0]))
            return "/" + drive
        })
    }
    
    return result
}
```

**4. フラグ競合の解決**

**問題**: `-h` フラグが help と競合（urfave/cli v3の仕様）

**症状**:
```bash
$ pathconv -h "/c/hoge"
エラー: No help topic for '/c/hoge'
```

**原因**: urfave/cli v3 は `-h` を自動的に help に割り当てる

**解決**: Home形式を `-H` (大文字) に変更
- `main.go` のフラグ定義: `Aliases: []string{"H"}`
- `flagMap`: `"H": paths.Home`
- README更新

**教訓**: CLI フレームワークの予約フラグを確認する

### Phase 3: テストとドキュメント

**テスト戦略**:
1. **単体テスト**: 各正規化・変換関数
2. **統合テスト**: 入力形式 × 出力形式の全組み合わせ
3. **エッジケース**: 相対パス、ドライブレターなし、UNCパスなど

**テストで発見した問題**:

#### 問題1: Escape形式のテスト期待値
```go
// ❌ 誤った期待値（raw stringの理解不足）
{`C:\Users\test`, `C\\Users\\test`}  // これは C:\Users\test を期待

// ✅ 正しい期待値
{`C:\Users\test`, `C:\\Users\\test`}  // エスケープ後は C:\\Users\\test
```

**原因**: Go の raw string リテラル (`...`) では `\` が1文字として扱われる
**解決**: テスト期待値を `\\` で記述

#### 問題2: 複数入力形式テストの期待値不整合
- `c:/hoge` の正規化後は `C:\hoge` (大文字化)
- テストの期待値を更新

## 重要な実装パターン

### Pattern 1: 純粋関数での変換
すべての変換関数は純粋関数（副作用なし）：
```go
func convertToXXX(path string, relative bool) string {
    // 入力に基づいて計算するだけ
    // ファイルI/O、グローバル状態変更なし
    return result
}
```

**利点**:
- テストが容易
- 並列実行可能
- 予測可能な動作

### Pattern 2: 入力形式の優先順位
`normalizeInputPath()` での判定順序:
1. `file:///` (最優先、特殊形式)
2. `file://`
3. `\\?\` (UNC拡張長パス)
4. 相対パス判定
5. `/mnt/` (WSL)
6. `/c/` (Linux)
7. その他 Windows形式

**理由**: より特殊な形式を先にチェックすることで、誤判定を防ぐ

### Pattern 3: フラグマッピング
フラグ名と Format enum の二重マッピング:
```go
flagMap := map[string]paths.Format{
    "windows": paths.Windows,
    "W":       paths.Windows,  // エイリアス
    // ...
}
```

**理由**: os.Args から直接フラグ順序を取得するため

## 今後の拡張可能性

### 追加可能な入力形式
- Cygwin形式: `/cygdrive/c/...`
- macOS形式: `/Volumes/...`
- Samba/SMB形式: `smb://server/share/...`

**実装方法**: `normalizeInputPath()` に判定ロジック追加

### 追加可能な出力形式
- PowerShell形式: `Microsoft.PowerShell.Core\FileSystem::C:\...`
- Python pathlib形式: `Path('C:/...')`
- JSON/YAML用エスケープ形式

**実装方法**:
1. `Format` enum に追加
2. `convertToXXX()` 関数実装
3. フラグ定義追加
4. テスト追加

### パフォーマンス最適化の余地
現在の実装は十分高速だが、さらに最適化するなら：
- 正規表現のプリコンパイル（現在は毎回コンパイル）
- パス判定の短絡評価の最適化
- メモリアロケーション削減

## トラブルシューティングガイド

### 問題: 期待した形式が出力されない
**チェック項目**:
1. フラグ名を確認（`-h` は使えない、`-H` を使う）
2. 入力パスの形式が想定通りか
3. 相対パスか絶対パスか

### 問題: テストが失敗する
**チェック項目**:
1. raw string リテラルのバックスラッシュ数を確認
2. 大文字小文字の統一（ドライブレターは大文字）
3. 正規化後の形式を確認（入力→Windows→出力）

### 問題: 新しい入力形式を追加したい
**手順**:
1. `normalizeInputPath()` に判定ロジック追加（適切な優先順位で）
2. `isRelative()` で新形式を考慮
3. テスト追加（`TestNormalizeInputPath`）
4. 統合テスト実行（`TestConvert_MultipleInputFormats`）

## 参考資料

### Windows パス仕様
- ドライブレター: `C:`, `D:`, etc.
- UNC: `\\server\share\path`
- UNC拡張長パス: `\\?\C:\very\long\path`
- 最大パス長: 260文字（UNC拡張で32,767文字）

### WSL パス仕様
- Windows C: → `/mnt/c/`
- Windows D: → `/mnt/d/`
- ケース: Windows側は大文字小文字を区別しない

### file:// URL 仕様
- ローカル: `file:///C:/path`
- UNC: `file://server/share/path`
- RFC 8089 準拠

## 開発者へのメモ

このプロジェクトは**シンプルさ**を最優先に設計されています：
- 複雑な状態管理なし
- 依存関係最小限
- 純粋関数での実装
- 明確な責任分離

新機能追加時もこの原則を守ってください。「複雑さ」より「明快さ」を選びましょう。
