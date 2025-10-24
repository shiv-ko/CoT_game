# 開発・運用コマンド集

## 1. Docker（全体起動・DBマイグレーション）

### サービス一括起動

```bash
docker-compose up --build
```

- フロントエンド（Next.js）、バックエンド（Go）、DB（PostgreSQL）がまとめて起動します。
- コード変更はホットリロードで即時反映。

### サービス停止

```bash
docker-compose down
```

### DBマイグレーション（手動）

```bash
# PostgreSQLコンテナに入る
docker-compose exec db psql -U postgres -d cot_game

# マイグレーションファイルを適用
\i /docker-entrypoint-initdb.d/migrations/20251024_add_scores_columns.sql

# テーブル構造確認
\d scores
```

---

## 2. バックエンド（Go APIサーバー）

### ローカル起動（手動）

```bash
cd backend
go run main.go
```

- `.env` で `DATABASE_URL` などを設定しておくこと

### テスト実行

```bash
cd backend
go test ./...
```

### Lint & フォーマット

```bash
cd backend
golangci-lint run ./...
go run mvdan.cc/gofumpt@v0.6.0 -l -w .
go run golang.org/x/tools/cmd/goimports@latest -w .
```

---

## 3. フロントエンド（Next.js）

### 開発サーバー起動

```bash
cd frontend
npm install
npm run dev
```

- ブラウザで [http://localhost:3000](http://localhost:3000) にアクセス

### 本番ビルド

```bash
cd frontend
npm run build
npm start
```

### Lint

```bash
cd frontend
npm run lint
```

---
