# Technical Context: pathconv

## Technology Stack

### Core Technologies
- **Language**: Go 1.23+
- **CLI Framework**: urfave/cli v3.0.0-beta1
- **Build Tool**: Taskfile (task v3.x)
- **Version Control**: Git

### Why Go?
1. **Single Binary**: 依存関係なしで配布可能
2. **Cross-Platform**: Windows/Linux/macOSで同一コード
3. **Fast Compilation**: 開発サイクルが高速
4. **Standard Library**: パス操作に必要な機能が充実
5. **Type Safety**: コンパイル時エラー検出

### Why urfave/cli v3?
1. **Modern API**: v2より洗練されたインターフェース
2. **Flag Handling**: 豊富なフラグタイプサポート
3. **Help Generation**: 自動生成されるヘルプ
4. **Aliases**: 短縮形フラグのサポート
5. **Community**: 活発なメンテナンス

## Development Setup

### 必須ツール
```bash
# Go 1.23以上
go version  # go version go1.23.x

# Task runner
task --version  # Task version: v3.x

# Optional: golangci-lint（推奨）
golangci-lint --version
```

### 環境構築
```bash
# リポジトリクローン
git clone https://github.com/eamat-dot/pathconv
cd pathconv

# 依存関係インストール
go mod download

# ビルド
task build
# または
go build -o ./bin/pathconv .

# テスト実行
task test
# または
go test ./...
```

### 開発ワークフロー
```bash
# 1. コード変更
# 2. フォーマット
go fmt ./...

# 3. ベット
go vet ./...

# 4. テスト
task test

# 5. ビルド
task build

# 6. 動作確認
./bin/pathconv -W /mnt/c/test
```

## Dependencies

### Direct Dependencies
```go
module github.com/eamat-dot/pathconv

go 1.23

require github.com/urfave/cli/v3 v3.0.0-beta1
```

### Why Minimal Dependencies?
- **信頼性**: 依存関係の脆弱性リスク最小化
- **保守性**: アップデート負担の軽減
- **バイナリサイズ**: 小さなバイナリ（< 5MB）
- **セキュリティ**: 攻撃面の縮小

### Standard Library Usage
- `os`: 環境変数、カレントディレクトリ取得
- `path/filepath`: パス操作、絶対パス判定
- `strings`: 文字列操作、置換、分割
- `regexp`: ドライブレター変換の正規表現
- `fmt`: 文字列フォーマット、エラー生成
- `testing`: テストフレームワーク

## Build Configuration

### Taskfile.yml
プロジェクトのビルド・テスト・リントタスクを定義：

```yaml
version: '3'

vars:
  BINARY_NAME: pathconv
  BUILD_DIR: ./bin

tasks:
  default:
    cmds:
      - task: build

  build:
    desc: Build the binary
    cmds:
      - go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}} .

  run:
    desc: Run the application
    cmds:
      - go run . {{.CLI_ARGS}}

  test:
    desc: Run tests
    cmds:
      - go test ./...

  test-coverage:
    desc: Run tests with coverage
    cmds:
      - go test -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out

  lint:
    desc: Run linters
    cmds:
      - go fmt ./...
      - go vet ./...
      - golangci-lint run

  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf {{.BUILD_DIR}}
      - rm -f coverage.out

  deps:
    desc: Download dependencies
    cmds:
      - go mod download
      - go mod tidy
```

### CI/CD Configuration
`.github/workflows/build-and-test.yml`:
- **Trigger**: Push/PR to main
- **Go Version**: 1.23.4
- **Steps**:
  1. Checkout
  2. Setup Go
  3. Format check (`gofmt`)
  4. Vet (`go vet`)
  5. Lint (`golangci-lint`)
  6. Test (`go test`)
  7. Build (`go build`)

## Technical Constraints

### Platform Constraints
1. **Windows-Centric**: 主にWindows環境での利用を想定
   - WSL、Git Bash、PowerShell対応が優先
   - macOS/Linuxでも動作するが、Windows形式への変換が中心

2. **Path Length**: OS依存
   - Windows: 260文字（従来）、32,767文字（UNC拡張長パス）
   - Linux/macOS: 4096文字（PATH_MAX）
   - 制約: 本ツールでは検証しない（OS任せ）

### Runtime Constraints
1. **環境変数依存**:
   - `HOME` (Unix系)
   - `USERPROFILE` (Windows)
   - 取得失敗時はチルダ展開を省略

2. **文字エンコーディング**:
   - UTF-8前提
   - Go内部はUTF-8で処理
   - ターミナルのエンコーディングはOS依存

### Performance Constraints
- **メモリ**: 1MB以下（典型的ケース）
- **実行時間**: 10ms以下
- **I/O**: なし（環境変数読み取りを除く）

## Code Quality Standards

### Formatting
- **gofmt**: 標準フォーマッタ使用
- **Line Length**: 制限なし（Go慣習に従う）
- **Imports**: goimportsで自動整理

### Naming Conventions
- **Functions**: camelCase（exported: PascalCase）
- **Variables**: camelCase
- **Constants**: PascalCase（exported）
- **Files**: snake_case.go

### Error Handling
```go
// Good: ラップしてコンテキスト追加
if err != nil {
    return fmt.Errorf("カレントディレクトリの取得に失敗: %w", err)
}

// Good: Graceful degradation
home := getHome()
if home == "" {
    // ホームが取得できなくても処理続行
}
```

### Testing
- **Unit Test**: 個別関数のテスト
- **Table-Driven**: 複数ケースをテーブルで定義
- **Naming**: `TestFunctionName_Scenario`

Example:
```go
func TestConvert_Windows(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {`C:\Users\test`, `C:\Users\test`},
        {`./sub/path`, `.\sub\path`},
    }
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := Convert(tt.input, Windows)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

## Development Tools

### Recommended VSCode Extensions
- **Go** (golang.go): 公式Go拡張
- **Go Test Explorer** (ethan-reesor.vscode-go-test-adapter)
- **Task** (task.vscode-task)

### Debugging
```bash
# Delve debugger
dlv debug -- -W /mnt/c/test

# または VSCode launch.json
{
    "name": "Debug pathconv",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}",
    "args": ["-W", "/mnt/c/test"]
}
```

### Profiling
```bash
# CPU profiling
go test -cpuprofile cpu.prof ./internal/paths
go tool pprof cpu.prof

# Memory profiling
go test -memprofile mem.prof ./internal/paths
go tool pprof mem.prof
```

## Deployment

### Distribution
1. **Binary Release**: GitHub Releasesで各プラットフォーム用バイナリ配布
2. **Go Install**: `go install github.com/eamat-dot/pathconv@latest`
3. **Package Managers**（将来）:
   - Homebrew (macOS/Linux)
   - Scoop (Windows)
   - Chocolatey (Windows)

### Build Targets
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o pathconv.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o pathconv

# macOS
GOOS=darwin GOARCH=amd64 go build -o pathconv
GOOS=darwin GOARCH=arm64 go build -o pathconv-arm64
```

## Troubleshooting

### Common Issues

#### Issue: `go: cannot find main module`
**Solution**: `go mod init` または正しいディレクトリで実行

#### Issue: `undefined: cli.Command`
**Solution**: `go mod download` で依存関係をインストール

#### Issue: テストが通らない
**Solution**: 
1. `go fmt ./...` でフォーマット
2. `go vet ./...` で静的解析
3. ログを確認して具体的なエラーを特定

#### Issue: ビルドしたバイナリが動かない
**Solution**:
1. 正しいGOOS/GOARCHでビルドしたか確認
2. 実行権限を確認 (`chmod +x pathconv`)
3. パスが通っているか確認

## Performance Benchmarks

### 目標
- 単一変換: < 1μs
- 10形式出力: < 10μs
- 起動からexit: < 10ms

### ベンチマーク実行
```bash
# ベンチマーク実行
go test -bench=. -benchmem ./internal/paths

# 結果例（目標）
BenchmarkConvert-8      1000000     500 ns/op    64 B/op    2 allocs/op
```

## Future Technical Debt
現時点では技術的負債は最小限：
- ✅ テストカバレッジ高い
- ✅ 依存関係最小
- ✅ コード品質良好
- ⚠️ ベンチマーク未実装（優先度低）
