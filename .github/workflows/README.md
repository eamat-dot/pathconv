# GitHub Actions CI

基本的な CI パイプラインを実行します。

## セットアップ

ワークフローファイルを `.github/workflows/` にプッシュすれば自動有効化されます。

```bash
git add .github/workflows/
git commit -m "ci: add basic GitHub Actions CI"
git push
```

## 実行

### 自動実行
- `main` / `develop` ブランチへのプッシュ時
- PR 作成時

### 実行内容
1. **go fmt** - コード形式チェック
2. **go vet** - 静的解析
3. **go test -v** - テスト実行
4. **go build** - ビルド

## ローカルテスト

CI 実行前にローカルで検証：

```bash
task lint   # go fmt + go vet
task test   # go test
task build  # ビルド
```

または全て実行：

```bash
task all
```

## ワークフロー確認

GitHub リポジトリ → **Actions** タブで実行状況確認
