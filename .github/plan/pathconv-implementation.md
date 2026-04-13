# pathconv Go 実装計画

## ゴール
- `pathconv` エイリアス（Lua/nyagos）の挙動を Go 製単一コマンドとして再現し、引数で任意パスも変換できるようにする。
- urfave/cli v3 を用いた単一コマンド（サブコマンドなし）。
- デフォルトはカレントディレクトリを基準に全形式を出力。引数でパスが与えられた場合はそのパスを基準に出力。

## 期待する入出力仕様
  - `-W, --windows` : Windows 形式だけを出力（旧 `-p, --pwd` から改名）
  - `-h, --home`: home 形式だけを出力
  - `-d, --dopus`: DOpus 形式だけを出力
  - `-v, --vscode`: VSCode 形式だけを出力
  - `-g, --gitbash`: Git Bash 形式だけを出力
  - `-w, --wsl`: WSL 形式だけを出力
  - `-l, --linux`: Linux 形式だけを出力
  - `-e, --escape`: WinEscape 形式だけを出力
- **出力**: 1 行 1 形式の文字列。パスにスペースが含まれる場合は二重引用符で囲む（nyagos踏襲）。

## 変換ロジック（nyagosの移植方針）
- **入力パスの正規化**:
  - 引数なし: `os.Getwd()` を使用。
  - 引数あり: 次のルールで正規化する。
    - `~/...` で始まる場合: HOME（`USERPROFILE` 等）に展開し、フルパス系（windows/wsl/linux/escape）では絶対パス化。
    - 相対パス（例: `./sub/path`, `../x`）: 存在確認を行わず、相対のまま扱う（絶対パス化しない）。
    - 絶対パス（例: `C:\...`, `/mnt/c/...`）: そのまま受け取り、各形式に応じて変換。
  - パス区切りは環境依存のまま受け取り、変換時に必要に応じてスラッシュ変換。
- **HOME 判定**: `HOME` または `USERPROFILE` を使用。HOME 配下でないパスが入力された場合は home/dopus/vscode/gitbash でも置換せずフルパスを返す。環境変数が取得できないレアケースでもエラーにはせずフルパス返却。
-- **各形式の定義**（元 Lua 実装と同等）：
  - windows: そのまま（Windows区切り維持）。相対入力時は `.` を残し `\\` ではなく `\` を使用（例: `./sub/path` → `\.\sub\path` ではなく ` .\sub\path`）。
  - home: HOME プレフィックスを `~` に置換し、区切りを `/` に。`C:\Users\me\Doc` → `~/Doc`。
  - WSL: 区切りを `/` にし、`^([A-Za-z]):` を `/mnt/<lower>` に置換。
  - Linux: 区切りを `/` にし、`^([A-Za-z]):` を `/<lower>` に置換。
  - WinEscape: `\` を `\\` にエスケープ。
  - VSCode: HOME を `${env:USERPROFILE}` に置換し、区切りを `/` に。
  - GitBash: HOME を `$USERPROFILE` に置換し、区切りを `/` に。
  - DOpus: HOME を `/profile` に置換し、区切りを `/` に。
  - UNC: 絶対 Windows パスは拡張長プレフィックス `\\?\` を付与（例: `C:\path` → `\\?\C:\path`）。UNC ネットワークパス（`\\server\share\...`）はそのまま。相対入力時は UNC 非対応のため Windows 形式に準拠（例: `./sub/path` → ` .\sub\path`）。
  - URL: `file://` URL に変換。`C:\path` → `file:///C:/path`、`\\server\share\path` → `file://server/share/path`。スペース等は URL エンコード（` ` → `%20`）。相対入力時は `./sub/path`（必要に応じて `%20` などをエンコード）。
- **クオート**: パスにスペースが含まれる場合は `"` で囲む。
- **出力順序**: デフォルト（フラグなし）は固定順 [windows, home, dopus, vscode, gitbash, wsl, linux, escape]。複数フラグ指定時は指定順で出力。

## CLI 設計（urfave/cli v3）
- 単一コマンド `pathconv`。
- フラグ: 上記 10 種類（`BoolFlag`）。複数指定された場合は指定されたものだけを**指定順で出力**する。
  - 相対入力の例（存在確認なし）:
    - `-W` → ` .\sub\path`
    - `-g` → ` ./sub/path`
    - `-w` → ` ./sub/path`
    - `-e` → ` .\\sub\\path`
    - `-U` → ` ./sub/path`（スペースは `%20`）
- 引数: 0..1 パス。2 個以上はエラーとしてヘルプ表示。
- ヘルプ: nyagosの説明文を移植（日本語可）。

## 実装ステップ案
1) `go mod init`（まだ存在しない場合）と urfave/cli v3 を追加。
2) `cmd` は不要。`main.go` で CLI を組み立て、変換ロジックは `internal/paths` に切り出す。
   - `internal/paths/convert.go` にフォーマット変換関数群。
   - `internal/paths/format.go` にフォーマット種別の列挙や変換テーブル（テストしやすい形）。
3) 変換ユーティリティの単体テストを追加（Windows/Unix で挙動が変わる箇所はスラッシュ変換前提でテストケースを分ける）。
4) `main.go` でフラグを解析し、対象フォーマットのリストを組み立てて出力。
5) Taskfile の `BINARY_NAME` を `pathconv` に更新し、ビルド/CI のバイナリ名を合わせる。

## テスト方針
- `TestQuoteIfNeeded`: スペースあり/なしでのクオート挙動。
- 各フォーマット関数の単体テスト（入力パス + HOME を与え、期待文字列を比較）。
- CLI レベルの軽いテスト: `go test ./...` に含める（フラグ組み合わせを最小限で）。
- CI: 既存の GitHub Actions と Taskfile の `task lint/test` を活用。

## オープンな検討事項
- なし（`~` 展開方針 / HOME 未定義時の扱い / フラグ指定順を決定済み）。
