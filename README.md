# pathconv

カレントディレクトリまたは指定されたパスを様々な形式で表示する CLI ツール。

様々な入力形式（WSL パス、Linux パス、Windows パス、UNC パス、file URL など）を自動判定・正規化して、Windows、WSL、Linux、VS Code、Git Bash、Directory Opus など、複数の環境で使えるパス形式に変換できます。

## このリポジトリについて

個人用に作ったツールを置いています。
AIコーディングやCIまわりの練習も兼ねています。

## インストール

<!-- 
```bash
go install github.com/eamat-dot/pathconv@latest
```

 -->

ソースからビルド：

```bash
git clone https://github.com/eamat-dot/pathconv
cd pathconv
go build -o ./bin/pathconv .
```

## 使い方

### 基本的な使用方法

引数なしで実行すると、カレントディレクトリを全形式で表示：

```bash
$ pathconv
C:\Users\user\Documents\project
~/Documents/project
/profile/Documents/project
${env:USERPROFILE}/Documents/project
$USERPROFILE/Documents/project
/mnt/c/Users/user/Documents/project
/c/Users/user/Documents/project
C:\\Users\\user\\Documents\\project
\\?\C:\Users\user\Documents\project
file:///C:/Users/user/Documents/project
```

### 特定の形式だけを表示

フラグを使って必要な形式だけを出力できます。フラグは指定順で出力されます。

```bash
$ pathconv -W -w
C:\Users\user\Documents\project
/mnt/c/Users/user/Documents/project
```

### 様々な入力形式に対応

入力パスの形式がどんな形式でも、正しく認識・変換されます：

```bash
# WSL形式の入力
$ pathconv -W "/mnt/c/hoge"
C:\hoge

# Linux形式の入力
$ pathconv -W "/c/hoge"
C:\hoge

# Windows形式（スラッシュ）の入力
$ pathconv -W "c:/hoge"
C:\hoge

# Windows形式（バックスラッシュ）の入力
$ pathconv -W "C:\hoge"
C:\hoge

# UNC拡張長パス形式の入力
$ pathconv -W "\\?\C:\hoge"
C:\hoge

# file URL形式の入力
$ pathconv -W "file:///C:/hoge"
C:\hoge
```

### 相対パスの変換

相対パスも変換できます（実在チェックはしません）：

```bash
$ pathconv -W -w ./sub/path
.\sub\path
./sub/path
```

### チルダ展開

`~` で始まるパスは HOME（USERPROFILE）に展開されます：

```bash
$ pathconv -W ~/Documents
C:\Users\user\Documents
```

## オプション

| フラグ | 長い形式 | 説明 | 出力例 |
|--------|----------|------|--------|
| `-W` | `--windows` | Windows 形式 | `C:\Users\user\...` |
| `-H` | `--home` | Home 形式（~ 展開） | `~/Documents/...` |
| `-d` | `--dopus` | Directory Opus 形式 | `/profile/Documents/...` |
| `-v` | `--vscode` | VS Code 形式 | `${env:USERPROFILE}/Documents/...` |
| `-g` | `--gitbash` | Git Bash 形式 | `$USERPROFILE/Documents/...` |
| `-w` | `--wsl` | WSL 形式 | `/mnt/c/Users/user/...` |
| `-l` | `--linux` | Linux 形式 | `/c/Users/user/...` |
| `-e` | `--escape` | Windows エスケープ形式 | `C:\\Users\\user\\...` |
| `-u` | `--unc` | UNC 拡張長パス形式 | `\\?\C:\Users\user\...` |
| `-U` | `--url` | file URL 形式 | `file:///C:/Users/user/...` |

## 入力形式対応

以下のいずれかの形式でパスを入力できます。プログラムが自動的に判定・正規化します：

- **WSL 形式**: `/mnt/c/Users/...`
- **Linux 形式**: `/c/Users/...`
- **Windows 形式**: `C:\Users\...` または `C:/Users/...`（大文字小文字と区切り文字を自動正規化）
- **UNC 拡張長パス**: `\\?\C:\Users\...`
- **file URL 形式**: `file:///C:/Users/...` または `file://server/share/...`

## 出力順序

- フラグなし: 全形式をデフォルト順（W, H, d, v, g, w, l, e, u, U）で出力
- フラグあり: **指定した順番**で出力

```bash
# -w が先、-W が後
$ pathconv -w -W
/mnt/c/Users/user/Documents
C:\Users\user\Documents
```

## 開発

### ビルド

```bash
task build
```

または

```bash
go build -o ./bin/pathconv .
```

### テスト

```bash
task test
```

または

```bash
go test ./...
```

### リント

```bash
task lint
```
