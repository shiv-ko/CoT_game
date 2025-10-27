# 23. 問題管理API実装

## 目的 / 背景

- 現状、問題(questions)は手動SQLでしか追加できない。
- 管理者が問題を作成・更新・削除できるAPIエンドポイントを提供する。
- 将来的には管理画面UIからも利用可能にする基盤を作る。

## 完了条件

- 問題のCRUD操作ができるRESTful APIエンドポイントが実装されている。
- 管理者認証・認可のミドルウェアが実装されている。
- 問題のバリデーション（必須項目チェック、正答形式チェック）が実装されている。
- 単体テスト・結合テストが実装されている。

## スコープ

- 含む:
  - 問題CRUD APIエンドポイント実装
  - 管理者権限チェックミドルウェア
  - リクエストバリデーション
  - テストコード
- 含まない:
  - 管理画面UI実装
  - 問題の一括インポート/エクスポート機能
  - 問題のバージョン管理機能

## 作業手順

### 1. リポジトリ層拡張

- [ ] `internal/repository/questions_repo.go` を作成
  - `QuestionsRepository` インターフェース定義
    - `Create(ctx, question) error`
    - `FindByID(ctx, id) (*Question, error)`
    - `FindAll(ctx, filters) ([]Question, error)`
    - `Update(ctx, id, question) error`
    - `Delete(ctx, id) error`
  - PostgreSQL実装 `postgresQuestionsRepo` を実装
- [ ] `internal/repository/questions_repo_test.go` を作成
  - 各メソッドの結合テスト実装

### 2. ミドルウェア実装

- [ ] `middleware/admin_auth.go` を作成
  - 管理者権限チェックミドルウェア実装
  - ユーザーロールベースの認可ロジック
  - テストヘルパー（モックユーザー作成）
- [ ] usersテーブルに `role` カラム追加のマイグレーション作成
  - `db/migrations/20251027_add_user_role.sql`
  - デフォルト: `user`, 管理者: `admin`

### 3. ハンドラ実装

- [ ] `handlers/question_handler.go` を拡張
  - `PostQuestion(c *gin.Context)` - 問題作成
  - `GetQuestion(c *gin.Context)` - 問題詳細取得（既存）
  - `GetQuestions(c *gin.Context)` - 問題一覧取得（既存）
  - `PutQuestion(c *gin.Context)` - 問題更新
  - `DeleteQuestion(c *gin.Context)` - 問題削除
- [ ] リクエスト/レスポンス構造体定義
  - バリデーションタグ追加（`binding:"required"` など）
- [ ] エラーハンドリング統一

### 4. ルーティング設定

- [ ] `routes/question_routes.go` を拡張
  - 管理者用ルートグループ作成 `/api/v1/admin/questions`
  - AdminAuthミドルウェア適用
  - 各エンドポイントをマッピング
    - `POST /api/v1/admin/questions` - 作成
    - `GET /api/v1/admin/questions/:id` - 詳細
    - `GET /api/v1/admin/questions` - 一覧
    - `PUT /api/v1/admin/questions/:id` - 更新
    - `DELETE /api/v1/admin/questions/:id` - 削除

### 5. テスト実装

- [ ] `handlers/question_handler_test.go` を拡張
  - 管理者権限での各操作の成功ケース
  - 非管理者での操作の失敗ケース（403 Forbidden）
  - バリデーションエラーケース（400 Bad Request）
  - 存在しない問題へのアクセス（404 Not Found）
- [ ] `middleware/admin_auth_test.go` を作成
  - 権限チェックロジックのテスト

### 6. ドキュメント更新

- [ ] `docs/app.md` にAPI仕様を追記
  - エンドポイント一覧
  - リクエスト/レスポンス例
  - 認証・認可要件
- [ ] README.mdに管理者アカウント作成方法を追記

## API仕様例

### 問題作成

```http
POST /api/v1/admin/questions
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "level": 1,
  "problem_statement": "2+3を計算してください。",
  "correct_answer": "5"
}
```

### 問題更新

```http
PUT /api/v1/admin/questions/1
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "level": 2,
  "problem_statement": "10×5を計算してください。",
  "correct_answer": "50"
}
```

### 問題削除

```http
DELETE /api/v1/admin/questions/1
Authorization: Bearer <admin_token>
```

## データモデル拡張

### users テーブル

```sql
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user' NOT NULL;
-- 'user' | 'admin'
```

### 管理者ユーザー作成例

```sql
INSERT INTO users (username, password_hash, role)
VALUES ('admin', '$2a$10$...', 'admin');
```

## 依存関係

- タスク#1 (AI連携) - 完了済み
- タスク#17 (統一エラーレスポンス) - 推奨（エラー処理の一貫性）
- 認証機能 - 既存のauth実装を利用

## 見積もり

- リポジトリ層: 2-3時間
- ミドルウェア: 1-2時間
- ハンドラ: 2-3時間
- テスト: 2-3時間
- ドキュメント: 1時間
- **合計: 8-12時間**

## 備考

- 当面は手動SQLで問題を追加するため、優先度は中〜低。
- タスク#9-11 (フロント実装) の後に実装すると、管理画面UIとセットで提供できる。
- 問題の一括インポート機能（CSV/JSON）は別タスクとして分離推奨。
