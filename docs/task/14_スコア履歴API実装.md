# 14. スコア履歴API実装

## 目的 / 背景

- `docs/plan_0907.md` Sprint3: `GET /api/v1/scores`。
- ユーザー自身の履歴閲覧 UI のデータソース。

## 完了条件

- ルート: `GET /api/v1/scores?user_id=&from=&to=&limit=`。
- フィルタ: user_id 指定可（未指定は全体 or 拒否 → MVPでは必須にする方針でセキュリティ簡素化）。
- レスポンス: `[{ id, question_id, score, prompt_len, created_at }]`。
- 400: user_id 未指定。

## スコープ

- 含む: クエリパラメータ解析, リポジトリ呼出。
- 含まない: 認証/認可（後続タスクの可能性）。

## 作業手順

- [ ] ハンドラ `handlers/scores_history_handler.go`
- [ ] ルート `routes/scores_routes.go`
- [ ] リポジトリメソッド利用
- [ ] テスト (user_id 必須, limit 挙動)

## 依存関係 / リスク

- 依存: スコア保存リポジトリ層実装
- リスク: データ量増加 → MVP段階で limit デフォルト 50

## 見積り / 担当 / 期日

- 見積り: 0.5人日
- 担当: 未定
- 期日: Sprint3 (2025-09-25)

## 参考リンク

- `../plan_0907.md`
