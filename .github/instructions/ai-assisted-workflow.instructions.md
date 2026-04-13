---
description: 'AI アシスト型ワークフローのルール。Issue実装・PR レビュー・チャットボックス指示など、すべての AI 支援作業で統一されたルールを適用。'
applyTo: '**'
---

# AI アシスト型ワークフロー統一ルール

このドキュメントは、GitHub Copilot やその他の AI エージェントが本リポジトリで Issue 実装・PR レビュー・チャットボックス指示を実行する際に遵守すべき統一ルールです。`.github/prompts/` の各プロンプトもこのルールに準拠します。

## ブランチ名の付け方

すべてのブランチは以下のパターンで命名してください：

```
${prefix}_${issue-or-PR_number}/[説明的なブランチ名]
```

### Prefix の種類

| Prefix | 用途 | 例 |
|--------|------|-----|
| `fix` | バグ修正 | `fix_42/rakuten-author-parsing-bug` |
| `feature` | 機能追加 | `feature_15/google-api-optimization` |
| `docs` | ドキュメント更新 | `docs_8/setup-guide-enhancement` |
| `refactor` | リファクタリング | `refactor_20/types-normalization` |
| `test` | テスト追加・修正 | `test_33/kobo-integration-tests` |
| `chore` | ビルド・CI・依存関係など | `chore_99/github-actions-upgrade` |

### 説明的なブランチ名の付け方
- 英数字とハイフンのみ使用（スペース・アンダースコア・日本語は非推奨）
- 簡潔・明確（最大50文字程度、短い方が望ましい）
- 変更の「何を」に焦点（「why」は不要）

### 例
- ✅ `fix_42/handle-utf8-author-names`
- ✅ `feature_15/add-cache-layer-api-responses`
- ❌ `fix/42-handle-utf8-author-names` （アンダースコア位置違い）
- ❌ `fix_42/修正` （日本語）

## AI 起源の明記

AI がコード変更を実施した場合、必ず「AI による修正・実装」であることを明記してください。

### 1. コミット時

コミットメッセージの body に以下のいずれかを追記：

```
Co-authored-by: GitHub Copilot <noreply@github.com>
```

または本文に（日本語で OK）：

```
[AI による自動実装]
```

### 2. PR・Issue コメント時

修正やフィードバックを含むコメントの冒頭に以下を記載：

```
**AI による自動修正：** (修正内容の説明)
```

または

```
🤖 AI により修正・実装しました：(修正内容の説明)
```

### 3. PR Description 時

PR 説明の冒頭に commit type と AI 明示を含める：

```
fix: #${issue番号} AI による自動 PR

**Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。

(以下、通常の PR 説明)
```

## Copilot Reviewer・Copilot Coding Agent への言語要求

AI エージェント（特に Copilot Reviewer や Copilot Coding Agent）にレビュー・作業を依頼する際は、**必ず日本語対応を明示**してください。

### 方法1: PR Description に記載

```markdown
**Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。
```

### 方法2: Issue Comment に記載

```markdown
@copilot-reviewer このタスクのレビューコメントは日本語で返してください。
```

### 方法3: PR Comment に記載

修正提案時：

```markdown
**Copilot Coding Agent へ**: 修正内容は日本語で説明し、日本語コメント付きコード を返してください。
```

## Memory Bank との統合

AI アシスト型の作業は Memory Bank のタスク管理システムと統合します。

### 実施方法

1. **タスク開始時**
   - `manage_todo_list` で新規タスク作成（ID: `TASK###-[タスク名]`、status: `in-progress`）
   - `.github/memory-bank/tasks/TASK###-[タスク名].md` を作成

2. **進捗ログ更新**
   - 各ステップ完了後、該当タスクファイルの Progress Log を更新
   - Subtasks の状態を最新化

3. **完了処理**
   - タスクを `completed` に更新
   - 対応 Issue・PR 番号を記録
   - 最終ログを記述

詳細は `.github/memory-bank/memory-bank.instructions.md` を参照してください。

## チェックコマンド

実装後に commit・PR 作成前には必ず以下を実行：

```bash
task all  # lint / test / build の全チェック
```

## チャットボックス指示での適用例

### Example 1: Issue 実装指示

```
Issue #42 を実装してください。バグ修正なので fix_42/[説明的なブランチ名] でブランチを切ってください。
AI による実装であることをコミットに明記し、PR 作成時には日本語レビュー要求を忘れずに。
```

AI エージェントは以下を実行：
- ブランチ: `fix_42/handle-utf8-author-names`
- コミット body: `Co-authored-by: GitHub Copilot <noreply@github.com>`
- PR description: `**Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。`

### Example 2: PR レビュー修正指示

```
PR #10 のレビューコメント #3 に対応して修正してください。修正をコメントするときは「AI による自動修正」と明記してください。
```

AI エージェントは以下を実行：
- 修正実装
- コメント: `**AI による自動修正：** ...（修正内容の説明）`
- Resolve conversation を実行

## 参考ファイル

- `.github/prompts/auto_impl_issue.prompt.md` … Issue 実装プロンプト
- `.github/prompts/pull_request_review_fix.prompt.md` … PR レビュー修正プロンプト
- `.github/instructions/commit-message.instructions.md` … Conventional Commits ルール
- `.github/memory-bank/memory-bank.instructions.md` … Memory Bank タスク管理
