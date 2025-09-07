# プロジェクト概要

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

## 補足
- 機密情報は `.env` に置き、ルートの `.gitignore` で除外済みです。
- 詳細仕様は `docs/app.md` を参照してください。
- 進行中の開発計画は `docs/plan_0907.md` に記載しています。

