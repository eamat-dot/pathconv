# TASK001 - pathconv Go実装タスク化

**Status:** Completed  
**Added:** 2025-12-26  
**Updated:** 2025-12-27

## Original Request
- メモリバンクを初期化し、`pathconv-implementation.md` の内容をタスクに追加する。

## Thought Process
- nyagos `pathconv` を Go + urfave/cli v3 で単一コマンドとして移植する計画を、実行可能なタスク単位に分割する。
- まずは仕様確定と環境セットアップ、次に変換ロジックとCLI実装、最後にビルド/CI名称の整合とテストを進める。

## Implementation Plan
- [ ] 環境セットアップ: `go mod init`（未作成なら）と urfave/cli v3 追加
- [ ] 変換ロジック: `internal/paths` にフォーマット変換関数群と列挙を実装（windows/wsl/linux/escape/home/vscode/gitbash/dopus/unc/url）
- [ ] CLI実装: `main.go` で urfave/cli v3 を用いフラグ/引数処理と出力を組み立て（複数フラグは指定順で出力）
- [ ] フラグ改名: `-p, --pwd` を `-W, --windows` に変更（ヘルプ・README 反映）
- [ ] 相対パス対応: 存在確認せず相対のまま扱い、各形式ごとの出力例に合わせる（W/g/w/e）
- [ ] バイナリ名整合: Taskfile と CI の `BINARY_NAME` を `pathconv` に更新
- [ ] テスト: 変換関数と CLI の基本挙動をカバーするユニットテスト追加

## Progress Tracking

**Overall Status:** Completed - 100%

### Subtasks
| ID | Description | Status | Updated | Notes |
|----|-------------|--------|---------|-------|
| 1.1 | `go mod init` と urfave/cli v3 追加 | Complete | 2025-12-27 | go.mod作成、cli v3.6.1追加 |
| 1.2 | 変換関数群を `internal/paths` に実装 | Complete | 2025-12-27 | format.go, convert.go実装 |
| 1.3 | CLI 実装（フラグ/引数/出力） | Complete | 2025-12-27 | main.go実装、指定順出力対応 |
| 1.4 | Taskfile/CI のバイナリ名を `pathconv` に更新 | Complete | 2025-12-27 | Taskfile.yml, build-and-test.yml更新 |
| 1.5 | テスト追加（変換関数・CLI） | Complete | 2025-12-27 | convert_test.go, allformats_test.go追加 |

## Progress Log
### 2025-12-27
- go mod init と urfave/cli v3 依存追加完了
- internal/paths パッケージ実装（format.go で10形式定義、convert.go で変換ロジック実装）
- main.go で CLI 実装（フラグ指定順の出力を os.Args と cmd.Bool() の組み合わせで実現）
- Taskfile.yml と .github/workflows/build-and-test.yml のバイナリ名を pathconv に更新
- ユニットテスト追加（convert_test.go, allformats_test.go）
- 全テスト通過、ビルド成功、実動作確認完了（デフォルト出力、フラグ指定順、相対パス対応）
- README.md 作成（使い方、オプション、出力例を記載）
- タスク完了
