# Product Context: pathconv

## Why This Project Exists

開発者は日々、異なる環境とツール間を移動しながら作業します：
- Windows PowerShellでコマンド実行
- WSL（Windows Subsystem for Linux）でLinuxツール使用
- Git Bashでバージョン管理
- VS Codeで設定ファイル編集
- Directory Opusでファイル管理

各環境は独自のパス表記を要求します。同じディレクトリでも：
- PowerShell: `C:\Users\tama\Documents\project`
- WSL: `/mnt/c/Users/tama/Documents/project`
- Git Bash: `$USERPROFILE/Documents/project` または `/c/Users/tama/Documents/project`
- VS Code設定: `${env:USERPROFILE}/Documents/project`
- file URL: `file:///C:/Users/tama/Documents/project`

この変換を手動で行うのは：
- **時間の無駄**: 毎回考えて書き直す
- **エラーの原因**: スラッシュの向き、ドライブレター、エスケープを間違える
- **フロー中断**: 本来の作業から気をそらされる

## Problems This Solves

### 主要な問題
1. **パス形式の手動変換**
   - コピーしたパスを別環境で使うために書き換える手間
   - 特に設定ファイルやスクリプトで頻発

2. **エスケープやフォーマットの記憶**
   - `\` を `\\` にエスケープする必要があるか？
   - `/` と `\` どちらを使うべきか？
   - ドライブレターは大文字？小文字？

3. **複数ツール間のパス互換性**
   - WSLからWindowsアプリにパスを渡す
   - Windows CLIからLinuxツールにパスを渡す
   - 設定ファイルに貼り付けられる形式に変換

### 具体的なユースケース

#### ユースケース1: WSLからWindowsツール呼び出し
```bash
# WSLで作業中、現在のディレクトリをWindows Explorerで開きたい
pwd  # /mnt/c/Users/tama/project
# → pathconv を使って Windows形式に変換
explorer.exe $(pathconv -W .)
```

#### ユースケース2: 設定ファイルへのパス記述
```json
// VS Code settings.json に書くパスを取得
// pathconv -v ~/project/configs
// → ${env:USERPROFILE}/project/configs
```

#### ユースケース3: エイリアス・スクリプト作成
```bash
# 複数環境で動くエイリアスを作りたい
# pathconv で各環境用の形式を一度に取得
pathconv ~/tools
# → 全形式が出力されるので、必要なものをコピペ
```

#### ユースケース4: ドキュメント作成
```markdown
# README に複数環境向けの手順を書く
# pathconv -W -w -l /mnt/c/project
Windows: C:\project
WSL: /mnt/c/project
Linux: /c/project
```

## How It Should Work

### 理想的なユーザー体験
1. **即座に変換**: コマンド入力から結果表示まで瞬時
2. **直感的**: フラグ名が覚えやすく、自然
3. **柔軟**: 必要な形式だけを取得可能
4. **正確**: どんな入力形式でも正しく解釈

### コア機能

#### 自動判定
ユーザーは入力形式を意識する必要がない：
```bash
pathconv -W "/mnt/c/hoge"     # WSL → Windows
pathconv -W "/c/hoge"         # Linux → Windows
pathconv -W "c:/hoge"         # Windows(/) → Windows(\)
pathconv -W "C:\hoge"         # Windows(\) → Windows(\)
pathconv -W "file:///C:/hoge" # URL → Windows
```
すべて同じ結果: `C:\hoge`

#### シンプルなインターフェース
```bash
# 引数なし = カレントディレクトリの全形式
pathconv

# フラグで特定形式のみ
pathconv -W        # Windows形式のみ
pathconv -w        # WSL形式のみ
pathconv -W -w     # 両方（指定順）

# パス指定
pathconv -W ~/Documents
pathconv -w /c/temp
```

#### コピー&ペースト最適化
出力はそのままコピーして貼り付け可能：
- クオートは必要な場合のみ（スペース含むパス）
- エスケープ済み（`-e` オプション）
- そのまま設定ファイルに記述可能

### User Goals
1. **高速なワークフロー**: パス変換で作業を中断されない
2. **正確性**: 変換ミスによるエラーを防ぐ
3. **学習コスト最小**: すぐに使い始められる
4. **汎用性**: あらゆるパス変換シナリオに対応

### Non-Goals（意図的にサポートしない）
1. パスの存在確認や作成（pure変換ツール）
2. GUIインターフェース（CLI専用）
3. パス以外の変換（URL、環境変数展開など）
4. 複雑な設定ファイル（シンプルさ重視）

## Value Proposition
**「どんなパス形式でも、1コマンドで必要な形式に変換」**

開発者の時間を節約し、パス変換のストレスから解放するツール。
