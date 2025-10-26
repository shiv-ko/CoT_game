# プロジェクト概要

[![CI](https://github.com/shiv-ko/CoT_game/actions/workflows/ci.yaml/badge.svg)](https://github.com/shiv-ko/CoT_game/actions/workflows/ci.yaml)

本リポジトリは「プロンプトバトルWebアプリ」のモノレポ構成です。フロントエンド（Next.js）／バックエンド（Go）／DB（PostgreSQL）／ドキュメントで構成されています。

## ディレクトリ構成（2025-09-07 現在）

```

.
.env
README.md
backend/
  Dockerfile
  go.mod
  go.sum
  handlers/
    question_handler.go
  main.go
  models/
    question.go
  routes/
    question_routes.go

db/
  init.sql
  migrations/
    20251024_add_scores_columns.sql

docker-compose.yml

docs/
  app.md
  guideline_code.md
  plan_0907.md
  progresses/
    progress_0810.md
    progress_instruction.md

frontend/
  .gitignore
  Dockerfile
  README.md
  components/
    SignupForm.tsx
  eslint.config.mjs
  next-env.d.ts
  next.config.ts
  package-lock.json
  package.json
  postcss.config.mjs
  public/
    file.svg
    globe.svg
    next.svg
    vercel.svg
    window.svg
  services/
    api.ts
    authService.ts
  src/
    app/
      favicon.ico
      globals.css
      layout.tsx
      page.tsx
  tsconfig.json
  types/
    auth.ts
```

## データベースマイグレーション

`db/migrations/` ディレクトリにマイグレーションファイルを管理しています。

### 初回セットアップ

初回の `docker-compose up` 時に `db/init.sql` が自動実行され、基本的なテーブル構造が作成されます。

### マイグレーションの適用手順

マイグレーションファイルは手動で適用する必要があります。以下の手順で実行してください:

```bash
# 1. PostgreSQL コンテナに接続
docker-compose exec db psql -U postgres -d cot_game

# 2. マイグレーションファイルを実行
\i /docker-entrypoint-initdb.d/migrations/20251024_add_scores_columns.sql

# 3. 正常に適用されたか確認
\d scores
```

### ロールバック方法

問題が発生した場合は、各マイグレーションファイル内にコメントで記載されているロールバック SQL を実行してください。

例（`20251024_add_scores_columns.sql` のロールバック）:

```sql
ALTER TABLE scores
DROP COLUMN model_vendor,
DROP COLUMN model_name,
DROP COLUMN answer_number,
DROP COLUMN latency_ms,
DROP COLUMN evaluation_detail;
```

### 利用可能なマイグレーション

- `20251024_add_scores_columns.sql`: `scores` テーブルに AI モデル情報、抽出数値、レイテンシ、評価詳細の各カラムを追加

## 補足

- 機密情報は `.env` に置き、ルートの `.gitignore` で除外済みです。
- 詳細仕様は `docs/app.md` を参照してください。
- 進行中の開発計画は `docs/plan_0907.md` に記載しています。
