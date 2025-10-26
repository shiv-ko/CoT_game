## タスク別ブランチ / PR 命名案

`docs/task.md` のタスクリスト順に Issue 番号を #1〜#22 で割り当てた想定で整理している。PR タイトルは必ず `#番号: 命令形の英語文` フォーマットとし、ブランチ名はケバブケースで記載。

| Issue  | タスク                            | ブランチ名                           | PR タイトル                                                               |
| ------ | --------------------------------- | ------------------------------------ | ------------------------------------------------------------------------- |
| #1     | AI連携クライアント実装            | `feat/01-ai-client`                  | `#1: Implement AI client`                                                 |
| #2     | 評価ロジック実装                  | `feat/02-eval-logic`                 | `#2: Implement evaluation logic`                                          |
| #3     | evaluatorユニットテスト           | `test/03-evaluator-unit-tests`       | `#3: Add evaluator unit tests`                                            |
| #4     | scores拡張マイグレーション        | `chore/04-scores-migration`          | `#4: Extend scores migration`                                             |
| #5-7   | スコア保存・Solve・結合テスト実装 | `feat/05-07-score-solve-integration` | `#5-7: Implement score repository, solve endpoint, and integration tests` |
| #8     | CIワークフロー整備                | `chore/08-ci-workflow`               | `#8: Establish CI workflow`                                               |
| #9-11  | フロント問題一覧・詳細・結果UI    | `feat/09-11-problem-ui`              | `#9-11: Build problem list, detail form, and result component`            |
| #12    | E2Eハッピーパス                   | `test/12-e2e-happy-path`             | `#12: Add e2e happy path`                                                 |
| #13-14 | ランキング・履歴API実装           | `feat/13-14-ranking-history-api`     | `#13-14: Implement leaderboard and score history api`                     |
| #15-16 | ランキング・履歴UIページ          | `feat/15-16-ranking-history-ui`      | `#15-16: Build leaderboard and personal history page`                     |
| #17    | 統一エラーレスポンス整備          | `chore/17-error-response`            | `#17: Standardize error responses`                                        |
| #18    | 構造化ログ基盤整備                | `chore/18-structured-logging`        | `#18: Establish structured logging`                                       |
| #19    | レートリミットミドルウェア        | `feat/19-rate-limit-middleware`      | `#19: Implement rate limit middleware`                                    |
| #20    | `.env.example` 整備               | `chore/20-env-example`               | `#20: Provide env example`                                                |
| #21    | ドキュメント更新(app.md補強)      | `docs/21-app-docs`                   | `#21: Enrich app documentation`                                           |
| #22    | モデル切替抽象化                  | `feat/22-model-switch`               | `#22: Abstract model switching`                                           |

> 備考: 実際の Issue 番号が異なる場合はタイトルの番号のみ合わせて調整する。ブランチ名フォーマットは変更なく利用可能。
