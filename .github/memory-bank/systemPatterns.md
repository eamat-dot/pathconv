# System Patterns: pathconv

## Architecture Overview

pathconvは3層のシンプルなアーキテクチャで構成：

```
┌─────────────────────────────────────┐
│         CLI Layer (main.go)         │
│  - Argument parsing (urfave/cli)    │
│  - Flag management                  │
│  - Output formatting                │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│   Business Logic (internal/paths)   │
│  - Path normalization               │
│  - Format conversion                │
│  - Path type detection              │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│    Standard Library (os, strings)   │
│  - Path operations (filepath)       │
│  - String manipulation              │
│  - Environment variables            │
└─────────────────────────────────────┘
```

## Key Technical Decisions

### 1. Input Normalization Pipeline
**決定**: すべての入力を内部的にWindows形式に正規化してから変換

**理由**:
- 各変換関数がWindows形式のみを期待できる（シンプル化）
- 入力形式の追加が容易（正規化関数のみ拡張）
- テストが書きやすい

**実装**:
```go
func Convert(inputPath string, format Format) string {
    // 1. 入力を正規化（様々な形式 → Windows形式）
    path := normalizeInputPath(inputPath)
    
    // 2. ~ 展開（必要な場合）
    if strings.HasPrefix(path, "~") {
        path = expandTilde(path)
    }
    
    // 3. 相対/絶対判定
    relative := isRelative(path)
    
    // 4. 指定形式に変換
    result := convertToFormat(path, format, relative)
    
    // 5. 必要に応じてクオート
    return QuoteIfNeeded(result)
}
```

**トレードオフ**:
- ✅ コードの単純化、保守性向上
- ✅ 新形式追加が容易
- ⚠️ 入力形式情報が失われる（問題なし：出力に影響しない）

### 2. Format Type System
**決定**: enumパターンで形式を型安全に定義

**実装**:
```go
type Format int

const (
    Windows Format = iota
    Home
    DOpus
    VSCode
    GitBash
    WSL
    Linux
    Escape
    UNC
    URL
)
```

**理由**:
- コンパイル時型チェック
- switchでの網羅性チェック
- パフォーマンス（文字列比較不要）

### 3. Converter Function Pattern
各形式への変換は独立した純粋関数：

```go
func convertToWindows(path string, relative bool) string
func convertToHome(path string, home string, relative bool) string
func convertToWSL(path string, relative bool) string
// ...
```

**メリット**:
- 単体テストが容易
- 並列実行可能（将来的に）
- 責任が明確

### 4. Flag Order Preservation
**決定**: ユーザーが指定したフラグの順序を保持して出力

**実装**:
```go
// os.Args をスキャンしてフラグの出現順を記録
for _, arg := range os.Args[1:] {
    if format, ok := flagMap[flagName]; ok {
        if !seen[format] {
            formats = append(formats, format)
            seen[format] = true
        }
    }
}
```

**理由**:
- ユーザーの意図を尊重
- パイプライン処理で便利
- 予測可能な動作

## Design Patterns in Use

### 1. Strategy Pattern
形式変換をStrategyとして実装：
- `Format` enum = Strategy識別子
- `convertToXXX()` 関数群 = 具体的Strategy

### 2. Facade Pattern
`Convert()` 関数が複雑な変換ロジックをシンプルなインターフェースで提供：
```go
result := Convert(inputPath, Windows)
```

### 3. Template Method Pattern
変換の基本フローは共通、各形式固有の処理は個別関数で実装：
```
正規化 → チルダ展開 → 相対判定 → 形式変換 → クオート
   ↑         ↑           ↑          ↑          ↑
 共通      共通       共通       形式固有    共通
```

### 4. Null Object Pattern
ホームディレクトリが取得できない場合、空文字列を返して後続処理を継続：
```go
func getHome() string {
    if home := os.Getenv("HOME"); home != "" {
        return home
    }
    if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
        return userprofile
    }
    return "" // Null object
}
```

## Component Relationships

### Dependency Graph
```
main.go
  └─> internal/paths/
        ├─> format.go      (型定義)
        ├─> convert.go     (変換ロジック)
        └─> (test files)   (テスト)
```

### Data Flow
```
User Input (CLI)
    ↓
Argument Parsing (urfave/cli)
    ↓
Path Detection & Normalization
    ↓
Format Conversion
    ↓
Quote Processing
    ↓
Output (stdout)
```

## Code Organization

### ディレクトリ構造
```
pathconv/
├── main.go                    # CLIエントリーポイント
├── internal/
│   └── paths/
│       ├── format.go          # Format enum定義
│       ├── convert.go         # 変換ロジック
│       ├── convert_test.go    # 単体テスト
│       ├── normalize_test.go  # 正規化テスト
│       └── allformats_test.go # 統合テスト
├── bin/                       # ビルド出力
├── Taskfile.yml              # タスク定義
├── go.mod                    # 依存管理
└── README.md                 # ドキュメント
```

### 責任分離
- **main.go**: CLI、フラグ処理、出力フォーマット
- **internal/paths/format.go**: Format型定義、Stringメソッド
- **internal/paths/convert.go**: パス変換ロジック、正規化、判定

## Error Handling Strategy

### 原則
1. **Fail Fast**: 無効な入力は早期に検出
2. **Graceful Degradation**: 可能な限り処理を継続
3. **Clear Messages**: ユーザーフレンドリーなエラーメッセージ

### 実装
```go
// 引数多すぎ
if cmd.Args().Len() > 1 {
    return fmt.Errorf("引数が多すぎます。パスは0個または1個のみ指定できます")
}

// カレントディレクトリ取得失敗
wd, err := os.Getwd()
if err != nil {
    return fmt.Errorf("カレントディレクトリの取得に失敗しました: %w", err)
}
```

## Performance Considerations

### 最適化ポイント
1. **メモリアロケーション削減**
   - 文字列操作は必要最小限
   - スライス容量を事前確保（flagsリスト）

2. **正規表現の効率的使用**
   - ドライブレター変換のみに使用
   - コンパイル済みパターンの再利用は不要（処理時間短い）

3. **I/O排除**
   - ファイル読み書き一切なし
   - すべてメモリ内処理

### パフォーマンス特性
- 実行時間: < 10ms（典型的ケース）
- メモリ使用: < 1MB
- I/O: 0回（環境変数読み取りのみ）

## Testing Strategy

### テストピラミッド
```
        ┌─────────┐
        │  E2E    │  CLI実行テスト（手動）
        └─────────┘
       ┌───────────┐
       │Integration│  複数入力→全出力テスト
       └───────────┘
      ┌─────────────┐
      │    Unit     │  個別変換関数テスト
      └─────────────┘
```

### テストカバレッジ目標
- **Unit**: 各変換関数を個別テスト
- **Integration**: 入力形式 × 出力形式の組み合わせ
- **Edge Cases**: 相対パス、空パス、特殊文字

### テストファイル構成
- `convert_test.go`: 個別形式変換のユニットテスト
- `normalize_test.go`: 入力正規化と統合テスト
- `allformats_test.go`: 全形式出力の統合テスト

## Security Considerations

### パス安全性
- パスの存在確認なし（pure変換）
- ファイルシステムアクセスなし
- パストラバーサル攻撃のリスクなし

### 入力バリデーション
- 長さ制限なし（OSの制限に依存）
- 特殊文字のサニタイズなし（そのまま変換）
- 理由: 変換ツールであり、実行ツールではない

## Future Extensibility

### 新形式追加の手順
1. `format.go` に新しいFormat定数を追加
2. `convert.go` に `convertToNewFormat()` 関数を実装
3. `Convert()` のswitch文に追加
4. テストケースを追加
5. `main.go` にフラグを追加

### プラグイン可能性
現状は単一バイナリだが、将来的には：
- カスタムフォーマッタのプラグイン機構
- Lua/JavaScript でのカスタムルール定義
