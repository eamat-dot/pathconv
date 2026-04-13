# Project Brief: pathconv

## Overview
pathconvは、様々なパス形式を相互変換するコマンドラインツールです。Windows、WSL、Linux、VS Code、Git Bash、Directory Opusなど、異なる環境間でパスをコピー&ペーストする際の煩わしさを解消します。

## Core Problem
開発者は日常的に複数の環境（Windows PowerShell、WSL、Git Bash、設定ファイルなど）を行き来します。各環境では異なるパス形式が必要となり、手動で変換するのは：
- 時間がかかる
- エラーが発生しやすい
- 作業フローを中断する

## Solution
pathconvは1つのコマンドで、入力パスを自動判定し、必要な形式に変換します。

## Key Requirements

### 機能要件
1. **多様な入力形式のサポート**
   - WSL形式: `/mnt/c/Users/...`
   - Linux形式: `/c/Users/...`
   - Windows形式: `C:\Users\...` または `C:/Users/...`
   - UNC拡張長パス: `\\?\C:\Users\...`
   - file URL形式: `file:///C:/Users/...`

2. **多様な出力形式のサポート**
   - Windows形式 (`C:\path`)
   - Home形式 (`~/path`)
   - DOpus形式 (`/profile/path`)
   - VSCode形式 (`${env:USERPROFILE}/path`)
   - GitBash形式 (`$USERPROFILE/path` または `/c/path`)
   - WSL形式 (`/mnt/c/path`)
   - Linux形式 (`/c/path`)
   - Escape形式 (`C:\\path`)
   - UNC形式 (`\\?\C:\path`)
   - URL形式 (`file:///C:/path`)

3. **自動正規化**
   - ドライブレターの大文字化（`c:/` → `C:\`）
   - スラッシュ方向の統一
   - 相対パスと絶対パスの適切な処理

4. **使いやすさ**
   - 引数なしでカレントディレクトリを表示
   - フラグで特定形式のみ出力
   - フラグの指定順で出力順序を制御

### 非機能要件
1. **パフォーマンス**: 即座にレスポンス（ファイルI/O不要）
2. **信頼性**: 広範なテストカバレッジ
3. **保守性**: クリーンなコード構造、明確な責任分離
4. **移植性**: Go製でクロスプラットフォーム

## Success Criteria
- [x] 6種類以上の入力形式を自動判定
- [x] 10種類の出力形式をサポート
- [x] 全テストが通過
- [x] CLI動作が直感的
- [x] ドキュメントが明確

## Out of Scope（現バージョン）
- パスの存在チェック
- ネットワークパスの詳細なバリデーション
- 設定ファイルによるカスタマイズ
- 対話モード

## Target Users
- マルチ環境で作業する開発者
- WSLユーザー
- 設定ファイルを頻繁に編集する人
- パスを頻繁にコピー&ペーストする人

## Project Timeline
- **開始**: 2025年12月27日
- **初期実装**: 2025年12月27日（同日完了）
- **現在フェーズ**: v1.0 - 基本機能完成

## Key Decisions
1. **Go言語を選択**: シングルバイナリ、クロスプラットフォーム、高速
2. **urfave/cli v3を採用**: モダンなCLIフレームワーク
3. **入力パスは常に正規化**: 一貫した内部処理
4. **フラグ指定順を保持**: ユーザーの意図を尊重
5. **Home形式を `-H` に変更**: `-h` とhelp の競合を回避
