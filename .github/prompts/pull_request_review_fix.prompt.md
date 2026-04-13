---
agent: 'agent'
description: 'pull_requestのレビューコメントを確認し、必要な修正を行うプロンプトです'
---

> **参照**: このプロンプトは `.github/instructions/ai-assisted-workflow.instructions.md` の統一ルールに準拠しています。チャットボックスからの PR 修正指示でも同じルールが適用されます。

## Steps（Memory Bank 統合版）
以下の手順に従って、pull_request のレビューコメント対応を進めます。各ステップ完了時に Memory Bank のタスク（`.github/memory-bank/tasks/TASK###-pr-review-fix.md`）を更新し、進捗ログと Subtasks を最新化してください。

1. タスク開始・初期化
   - `manage_todo_list` で「PR Review Fix」タスクを作成（ID: `TASK###-pr-review-fix`、status: in-progress）。
   - `.github/memory-bank/tasks/TASK###-pr-review-fix.md` を作成/更新し、Original Request と Subtasks の雛形を追加：
     - 1.1 レビューコメントの取得と分類
     - 1.2 コメント分析と修正計画の記録
     - 1.3 修正の実装（必要に応じてテスト追加）
     - 1.4 ステージング・コミット（Conventional Commits 準拠）
     - 1.5 ブランチの push・PR 更新
     - 1.6 コメントの解決と記録（番号付き）
     - 1.7 タスク完了処理（ステータス更新・最終ログ）
   - `tasks/_index.md` にタスクを In Progress として追記。

2. レビューコメントの取得と分類（Subtask 1.1）
   - 以下のツールをこの順序で試す：
     1. MCP_DOCKER ゲートウェイ経由: `pull_request_read`（最優先）
     2. GitHub Pull Requests 拡張: `github-pull-request_activePullRequest` または `github-pull-request_openPullRequest`
     3. 失敗時のみ `gh pr view`
   - resolved 済みコメントは除外。行が重複する場合は eamat-dot のコメントを優先。
   - 「必須修正（セキュリティ/バグ/機能）」「推奨修正（パフォーマンス/可読性/保守性）」に分類。
   - Memory Bank の Progress Log に取得結果と分類方針を記録。

3. コメント分析と修正計画（Subtask 1.2）
   - 各コメントの対応要否、破壊的変更の可能性、影響範囲、関連ファイル/モジュール、必要テストを判断。
   - 具体的な修正内容と実施順序をタスクファイルの「Implementation Plan」に記録。
   - ユーザー承認が必要な場合は計画を提示し、承認後に Subtask ステータス更新。

4. 修正の実装（Subtask 1.3）
   - 計画に基づき、小さく検証可能な単位で修正を実施。必要に応じてテストを追加/更新。
   - 各修正後に Memory Bank の Progress Log/ Subtasks を更新。
   - コミット前に「Check commands for after fix code」を参照し、`task all` で lint/ test/ build を実行して合格を確認。

5. ステージング・コミット（Subtask 1.4）
   - 変更ファイルのみステージング（例: `git add -u`、新規ファイルは個別指定）。
   - コミットメッセージは `.github/instructions/commit-message.instructions.md` に準拠（日本語・絵文字・50文字制限等）。
   - **コミットメッセージまたは本文に「AI による自動 fix」等の AI 起源を明記**（例: body に `Co-authored-by: GitHub Copilot <noreply@github.com>` や `[AI による自動修正]` を追記）。
   - Memory Bank にコミット概要を記録。

6. push・PR 更新（Subtask 1.5）
   - 現在のブランチを push し、PR を更新。必要なら CI 結果を確認。
   - **PR の説明（description）に以下の行を追記し、日本語レビューを要求**：
     ```
     **Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。
     ```
   - Memory Bank に push 時刻・PR URL（わかれば）を記録。

7. コメント解決（Subtask 1.6）
   - 修正対応済みのスレッドを PR ページで「Resolve conversation」。
   - もしくは `mcp_mcp_docker_add_comment_to_pending_review` でコメント追加時に解決。
   - **コメント本文に「AI による自動 fix」や「AI により修正しました」等の AI 起源を明記**。
   - チャットで「✅ コメント番号〇番: 修正完了・解決済み」を列記し、Memory Bank に同番号で記録。

8. タスク完了処理（Subtask 1.7）
   - `manage_todo_list` でタスクを completed に更新し、完了日時を記録。
   - `.github/memory-bank/tasks/TASK###-pr-review-fix.md` の Progress Log を締め、対応コメント番号一覧を追記。
   - チャットで「✅ タスク完了」を明示し、対応コメント番号を列記。

## Check commands for after fix code
実装後に変更をコミットする前には必ずcheck allを行ってください。
- check all（lint / test / build）: `task all`

