---
agent: 'agent'
description: 'issueを実装して自動でpull requestを作成するためのプロンプトです'
---

> **参照**: このプロンプトは `.github/instructions/ai-assisted-workflow.instructions.md` の統一ルールに準拠しています。チャットボックスからの Issue・PR 指示でも同じルールが適用されます。

## Steps（Memory Bank 統合版）
以下の手順に従って、issue 実装から Draft PR 作成まで進めます。各ステップ完了時に Memory Bank のタスク（`.github/memory-bank/tasks/TASK###-issue-implementation.md`）を更新し、進捗ログと Subtasks を最新化してください。

**Planning（ステップ1～4） → ユーザー承認 → Action（ステップ5～9）の順に作業します。**

### Planning フェーズ

1. タスク開始・初期化（Subtask 1.1）
   - `manage_todo_list` で「Issue Implementation」タスクを作成（ID: `TASK###-issue-implementation`、status: in-progress）。
   - `.github/memory-bank/tasks/TASK###-issue-implementation.md` を作成/更新し、Original Request と Subtasks の雛形を追加：
     - 1.1 タスク開始・初期化
     - 1.2 Issue の取得と内容確認
     - 1.3 ブランチ名決定と実装計画策定
     - 1.4 ユーザー承認の取得
     - 1.5 ブランチ作成と切り替え
     - 1.6 実装とテスト（細かくコミット）
     - 1.7 最終チェック（task all）
     - 1.8 ブランチ push と Draft PR 作成
     - 1.9 タスク完了処理
   - `tasks/_index.md` にタスクを In Progress として追記。

2. Issue の取得と内容確認（Subtask 1.2）
   - 以下のツールをこの順序で試す：
     1. MCP_DOCKER ゲートウェイ経由: `issue_read`（最優先）
     2. GitHub Pull Requests 拡張: `github-pull-request_issue_fetch`
     3. 失敗時のみ `gh issue view ${issue_number}`
   - Issue 本文とコメントを確認し、要件・制約・期待動作を整理。
   - Memory Bank の Progress Log に Issue 概要と要件を記録。

3. ブランチ名決定と実装計画策定（Subtask 1.3）
   - Issue タイプに応じてブランチ名を決定：
     - バグ修正: `fix_${issue_number}/[説明的なブランチ名]`
     - 機能追加: `feature_${issue_number}/[説明的なブランチ名]`
     - ドキュメント更新: `docs_${issue_number}/[説明的なブランチ名]`
     - リファクタリング: `refactor_${issue_number}/[説明的なブランチ名]`
     - テスト追加: `test_${issue_number}/[説明的なブランチ名]`
   - `codebases` tool が使える場合は修正すべきコードを検索してから実装計画を立案。
   - 追加・修正・削除するファイルを明示し、実装計画を `.github/memory-bank/tasks/TASK###-issue-implementation.md` の「Implementation Plan」に記録。
   - Memory Bank に影響範囲・関連ファイル・必要テストを記録。

4. ユーザー承認の取得（Subtask 1.4）
   - ブランチ名と実装計画をユーザーに提示し、承認を求める。
   - 承認が得られた場合は Subtask 1.4 を completed に更新し、Action フェーズへ進む。
   - 承認が得られなかった場合は Subtask 1.3 に戻り、ブランチ名と実装計画を再検討。Memory Bank にフィードバック内容を記録。

### Action フェーズ

5. ブランチ作成と切り替え（Subtask 1.5）
   - `git status` がクリーンであることを確認。
   - クリーンではない場合は現状を報告し、ユーザーの指示を待つ。
   - クリーンであれば `git checkout -b [ブランチ名]` でブランチ作成・切り替え。
   - Memory Bank にブランチ名と作成時刻を記録。

6. 実装とテスト（Subtask 1.6）
   - 実装計画に基づき、小さく検証可能な単位で実装。必要に応じてテスト追加/更新。
   - 各実装単位で以下を実行：
     1. コード変更
     2. 「Check commands for after fix code」を参照し、`task all` で lint/test/build を実行して合格確認
     3. 変更をステージング（例: `git add -u`、新規ファイルは個別指定）
     4. Conventional Commits 準拠でコミット（`.github/instructions/commit-message.instructions.md` 参照）
     5. **コミットメッセージまたは本文に「AI による自動実装」等の AI 起源を明記**（例: body に `Co-authored-by: GitHub Copilot <noreply@github.com>` や `[AI による自動実装]` を追記）
   - 各コミット後に Memory Bank の Progress Log と Subtasks を更新。

7. 最終チェック（Subtask 1.7）
   - すべての実装が完了したら `task all` で最終チェック。
   - 合格を確認し、Memory Bank にチェック結果を記録。

8. ブランチ push と Draft PR 作成（Subtask 1.8）
   - 現在のブランチを push（例: `git push origin [ブランチ名]`）。
   - 以下の方法でドラフトモードの Pull Request を作成：
     1. **最優先**: MCP_DOCKER ゲートウェイ経由の `create_pull_request`（ドラフトオプション指定）
     2. 次点: `github-pull-request_copilot-coding-agent`
     3. 失敗時のみ `gh pr create --draft`
   - PR のフォーマットは「Pull Request Format」セクション参照。
   - **PR の description に以下を追記し、日本語レビューを要求**：
     ```
     **Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。
     ```
   - Memory Bank に PR URL・作成時刻を記録。

9. タスク完了処理（Subtask 1.9）
   - `manage_todo_list` でタスクを completed に更新し、完了日時を記録。
   - `.github/memory-bank/tasks/TASK###-issue-implementation.md` の Progress Log を締め、実装内容・対応 Issue 番号・PR 番号を追記。
   - チャットで「✅ タスク完了」を明示し、Issue 番号・PR 番号を列記。

## Check Commands for After Fix Code
実装後に pull request を作成する前には必ず check all を行ってください。
- check all（lint / test / build）: `task all`

## Pull Request Format
- **タイトル**: 実装内容を簡潔にまとめる
- **本文の先頭行**: issue タイプに応じた commit type を使用
  - バグ修正: `fix: #${issue番号} AI による自動 PR`
  - 機能追加: `feat: #${issue番号} AI による自動 PR`
  - ドキュメント: `docs: #${issue番号} AI による自動 PR`
  - リファクタリング: `refactor: #${issue番号} AI による自動 PR`
  - テスト追加: `test: #${issue番号} AI による自動 PR`
- **本文の2行目**: 日本語レビュー要求
  ```
  **Copilot Reviewer へ**: このプルリクエストのレビューは日本語でお願いします。
  ```
- **本文の内容**:
  - 実装内容の説明
  - 実装にあたって参考にした情報や、特に注意した点
  - 実装にあたっての質問や懸念点
- **ブランチ**: `main` を指定
