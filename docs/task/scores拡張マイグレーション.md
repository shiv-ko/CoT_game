# scores拡張マイグレーション

## 目的 / 背景

- `docs/plan_0907.md` DB変更案に基づき `scores` テーブルへ拡張カラム追加。
- `POST /solve` 応答フィールド保存に必要。

## 完了条件

- `db/migrations/20250907_add_scores_columns.sql` (案) に相当するマイグレーションファイル作成。
- 追加カラム: `model_vendor TEXT NOT NULL DEFAULT 'gemini'`, `model_name TEXT NULL`, `answer_number NUMERIC NULL`, `latency_ms INT NOT NULL DEFAULT 0`, `evaluation_detail JSONB NULL`。
- 既存行にデフォルト適用されエラーなく適用可能。
- ロールバックSQL（DROP COLUMN セット）別ファイルまたは同ファイルコメント記述。
- ローカルで `docker-compose up --build` 時に自動適用されるか、適用手順 README 追記（後続ドキュメント更新タスク）。

## スコープ

- 含む: 追加 DDL, 既存との後方互換性確認。
- 含まない: インデックス最適化（後続）。

## 作業手順

- [ ] `db/migrations/` ディレクトリ作成
- [ ] 追加カラムDDL作成
- [ ] ロールバックDDLコメント記載
- [ ] 手動適用テスト (psql / migrate ツール想定)記録

## 依存関係 / リスク

- 依存: 評価ロジックで必要フィールド定義確定
- リスク: JSONB 未使用時の空領域 → 後続で NULL 運用

## 見積り / 担当 / 期日

- 見積り: 0.5人日
- 担当: 未定
- 期日: Sprint1 (2025-09-12)

## 参考リンク

- `../plan_0907.md`
