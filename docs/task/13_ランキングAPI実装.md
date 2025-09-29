# 13. ランキングAPI実装

## 目的 / 背景

- `docs/plan_0907.md` Sprint3: `GET /api/v1/leaderboard`。
- MVPで他ユーザー比較可視化を提供。

## 完了条件

- ルート: `GET /api/v1/leaderboard?period=day|week|all&limit=50`。
- レスポンス: `[{ user_id, username, best_score, attempts, last_at }]`。
- リポジトリ呼出で集計されたデータをそのまま返す。
- period 未指定時デフォルト: `day`。
- バリデーション: 不正 period -> 400。

## スコープ

- 含む: ハンドラ・ルート追加 / パラメータ検証。
- 含まない: キャッシュ, ページング。

## 作業手順

- [ ] ハンドラ `handlers/leaderboard_handler.go`
- [ ] ルート `routes/leaderboard_routes.go`
- [ ] リポジトリ呼出
- [ ] テスト (リポジトリモック)

## 依存関係 / リスク

- 依存: スコア保存リポジトリ層実装
- リスク: SQL集計コスト → limit 固定で軽減

## 見積り / 担当 / 期日

- 見積り: 1人日
- 担当: 未定
- 期日: Sprint3 (2025-09-25)

## 参考リンク

- `../plan_0907.md`
