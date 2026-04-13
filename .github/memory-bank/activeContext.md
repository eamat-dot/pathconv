# Active Context: pathconv

## Current Work Focus
**状態**: v1.0基本実装完了、本番利用可能

### 直近完了した作業（2025-12-27）

#### 1. 複数入力形式対応の実装
**問題**: 当初は特定の入力形式（Windows, 相対パス）のみサポート
**解決**: 入力パス正規化機能を実装し、6種類の入力形式に対応

実装内容：
- `normalizeInputPath()`: 入力形式を自動判定してWindows形式に正規化
  - WSL形式 (`/mnt/c/...`)
  - Linux形式 (`/c/...`)
  - Windows形式 (`C:/...`, `C:\...`, `c:/...`)
  - UNC拡張長パス (`\\?\C:\...`)
  - file URL (`file:///C:/...`, `file://server/...`)
- `normalizeWindowsPath()`: ドライブレター大文字化とスラッシュ統一
- `isRelative()`: file:// やドライブレター検出の改善

#### 2. GitBash形式の修正
**問題**: GitBash形式が `C:/hoge` のまま（Linux形式と同じ `/c/hoge` になるべき）
**解決**: `convertToGitBash()` でドライブレターを `/c` 形式に変換

#### 3. フラグ競合の解決
**問題**: `-h` フラグがhelpと競合（urfave/cli v3の仕様）
**解決**: Home形式を `-H` に変更
- `main.go` のフラグ定義を更新
- `flagMap` のマッピングを修正
- README を更新

#### 4. 広範なテスト追加
- `TestNormalizeInputPath`: 16パターンの入力形式正規化テスト
- `TestConvert_MultipleInputFormats`: 6入力 × 10出力 = 60パターンの変換テスト
- すべてのテストが通過

#### 5. ドキュメント整備
- README に「入力形式対応」セクション追加
- 各入力形式の出力例を明記
- フラグ一覧を更新（`-H` への変更を反映）

## Recent Changes

### コード変更
1. **internal/paths/convert.go**
   - `isRelative()`: ドライブレター・file:// 判定を追加
   - `normalizeInputPath()`: 入力パス正規化ロジック実装
   - `normalizeWindowsPath()`: ドライブレター大文字化
   - `convertToGitBash()`: Linux形式への変換追加

2. **internal/paths/normalize_test.go**（新規作成）
   - 入力正規化のユニットテスト
   - 複数入力形式から全出力形式への統合テスト

3. **internal/paths/convert_test.go**
   - 既存テストの期待値を更新（大文字化対応）

4. **main.go**
   - Home形式のフラグを `-h` → `-H` に変更
   - `flagMap` を更新

5. **README.md**
   - 入力形式対応セクション追加
   - フラグ一覧を更新
   - 使用例を拡充

### 動作確認済み
すべての入力形式から全出力形式への変換が正しく動作：
```bash
# 入力: /mnt/c/hoge, /c/hoge, c:/hoge, C:\hoge, \\?\C:\hoge, file:///C:/hoge
# 出力: すべて正しく各形式に変換される
```

## Next Steps

### 優先度: 高
現時点でコア機能は完成。次のステップは実運用フィードバック待ち。

### 優先度: 中（将来検討）
1. **追加フォーマットのサポート**
   - Cygwin形式 (`/cygdrive/c/...`)
   - MSYS2形式
   - その他のシェル環境

2. **設定ファイルサポート**
   - `.pathconvrc` でデフォルトフォーマット指定
   - カスタムエイリアス定義

3. **出力フォーマットのカスタマイズ**
   - テンプレート機能
   - JSON/YAML出力

### 優先度: 低
1. **対話モード**
   - fzfライクなインタラクティブ選択
   
2. **クリップボード統合**
   - 自動コピー機能

3. **シェル統合**
   - Zsh/Bash補完スクリプト
   - PowerShell モジュール

## Active Decisions

### 設計判断
1. **入力正規化を Convert() の最初に実行**
   - 理由: すべての変換ロジックが統一された入力を期待できる
   - トレードオフ: 若干の処理オーバーヘッド（実用上問題なし）

2. **ドライブレターは常に大文字化**
   - 理由: Windows標準、一貫性
   - トレードオフ: ユーザー入力の大文字小文字が保持されない（問題なし）

3. **GitBash形式は Linux形式と同じ `/c/` 表記**
   - 理由: Git Bash内部では `/c/` が標準
   - 代替案: `$USERPROFILE` 変数展開は別途対応済み

4. **file:// 形式の優先判定**
   - 理由: `isRelative()` より前にチェックしないと誤判定される
   - 実装: `normalizeInputPath()` の最初でチェック

### 保留中の決定
なし（すべての主要決定は完了）

## Known Issues
なし（すべてのテストが通過）

## Performance Notes
- すべての処理はメモリ内で完結
- ファイルI/O不要
- レスポンスは即座（< 10ms）

## Current Blockers
なし

## Dependencies
- Go 1.23+
- urfave/cli v3
- 標準ライブラリのみ（外部依存最小）
