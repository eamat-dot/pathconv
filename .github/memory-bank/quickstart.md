# Quick Start Guide: pathconv

## 開発者向けクイックスタート

### 1分でビルド・実行
```bash
# クローン
git clone <repository-url>
cd pathconv

# ビルド
go build -o ./bin/pathconv .

# 実行
./bin/pathconv -W /mnt/c/hoge
# → C:\hoge
```

### 5分で開発開始
```bash
# タスクランナー使用
task build   # ビルド
task test    # テスト実行
task lint    # 静的解析

# または直接
go test ./...
go vet ./...
```

## プロジェクト構造理解

```
pathconv/
├── main.go                          # CLIエントリポイント（フラグ処理）
├── internal/paths/
│   ├── format.go                    # Format enum定義
│   ├── convert.go                   # 変換ロジック（コア）
│   ├── convert_test.go              # 単体テスト
│   ├── normalize_test.go            # 正規化テスト
│   └── allformats_test.go           # 統合テスト
├── bin/                             # ビルド成果物
├── .github/
│   ├── memory-bank/                 # プロジェクト知識ベース（重要！）
│   │   ├── projectbrief.md          # プロジェクト概要
│   │   ├── activeContext.md         # 現在の状態
│   │   ├── implementation-history.md # 実装の詳細経緯
│   │   └── ...
│   └── workflows/                   # CI/CD設定
└── README.md                        # ユーザー向けドキュメント
```

## コア概念（3つだけ）

### 1. 入力正規化
すべての入力 → Windows形式 → 各出力形式
```
/mnt/c/hoge  ─┐
/c/hoge      ─┤
c:/hoge      ─┼→ normalizeInputPath() → C:\hoge → convertToXXX() → 出力
C:\hoge      ─┤
file:///C:/  ─┘
```

### 2. Format enum
型安全な形式定義:
```go
type Format int
const (
    Windows Format = iota
    Home
    // ... 10種類
)
```

### 3. 純粋関数
すべての変換は副作用なし:
```go
// 入力→出力、ファイルI/O・グローバル変数なし
func convertToXXX(path string, relative bool) string
```

## よくある作業

### 新しい出力形式を追加
1. `internal/paths/format.go` に Format 追加
2. `internal/paths/convert.go` に `convertToXXX()` 実装
3. `Convert()` の switch に追加
4. `main.go` にフラグ追加
5. テスト追加
6. README更新

### 新しい入力形式を追加
1. `internal/paths/convert.go` の `normalizeInputPath()` に判定ロジック追加
2. `isRelative()` で新形式を考慮（必要なら）
3. `TestNormalizeInputPath` にテスト追加
4. README の「入力形式対応」セクション更新

### バグ修正
1. 既存テストで再現確認
2. 必要なら新規テスト追加
3. 修正実装
4. `task test` で全テスト通過確認
5. `task lint` で静的解析クリア

## トラブルシューティング

### テスト失敗時
```bash
# 詳細出力
go test -v ./internal/paths

# 特定テストのみ
go test -v ./internal/paths -run TestConvert_Windows

# カバレッジ確認
go test -cover ./...
```

### ビルドエラー
```bash
# 依存関係再取得
go mod download
go mod tidy

# キャッシュクリア
go clean -cache
```

### フォーマット・lint
```bash
# フォーマット
go fmt ./...

# vet
go vet ./...

# golangci-lint（推奨）
golangci-lint run
```

## 重要ファイル（読むべき順）

1. **README.md** - ユーザー向け使い方
2. **.github/memory-bank/projectbrief.md** - プロジェクト概要
3. **.github/memory-bank/implementation-history.md** - 実装の詳細経緯（このチャットの要約）
4. **internal/paths/convert.go** - コアロジック
5. **main.go** - CLI実装

## CI/CD

GitHub Actions で自動実行:
- `go fmt` チェック
- `go vet` 実行
- `go test` 全テスト
- `go build` ビルド確認

プッシュ前に `task lint` と `task test` を実行推奨。

## リリース手順

1. バージョンタグ作成
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Releases でバイナリ添付（手動、またはGoreleaser導入予定）

3. go install 経由でのインストールも可能:
   ```bash
   go install github.com/eamat-dot/pathconv@latest
   ```

## ヘルプ・質問

- **実装の経緯**: `.github/memory-bank/implementation-history.md` 参照
- **設計判断**: `.github/memory-bank/systemPatterns.md` 参照
- **現在の状態**: `.github/memory-bank/activeContext.md` 参照

Memory Bank に全てが記録されています。
