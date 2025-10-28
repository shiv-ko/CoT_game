# セットアップガイド

このドキュメントは、プロンプトバトルWebアプリ（CoT_game）の開発環境をセットアップするための手順をまとめたものです。初学者でもスムーズに環境構築からアプリの起動までを行えるように、具体的なコマンドとその意味を記載しています。

---

## 📋 目次

1. [前提条件](#前提条件)
2. [リポジトリのクローン](#リポジトリのクローン)
3. [環境変数の設定](#環境変数の設定)
4. [Dockerを使った起動（推奨）](#dockerを使った起動推奨)
5. [ローカル開発環境での起動](#ローカル開発環境での起動)
6. [動作確認](#動作確認)
7. [開発用コマンド](#開発用コマンド)
8. [トラブルシューティング](#トラブルシューティング)

---

## 前提条件

開発を始める前に、以下のツールをインストールしてください。

### 必須ツール

#### 1. **Git**（バージョン管理ツール）
- **役割**: ソースコードのバージョン管理
- **インストール確認**:
  ```bash
  git --version
  ```
  `git version 2.x.x` と表示されればOK

- **インストール方法**:
  - **macOS**: `brew install git`
  - **Windows**: [Git for Windows](https://git-scm.com/download/win)
  - **Linux**: `sudo apt-get install git` (Ubuntu/Debian) または `sudo yum install git` (CentOS/RHEL)

#### 2. **Docker**（コンテナ仮想化ツール）
- **役割**: フロントエンド、バックエンド、データベースを統一環境で動かす
- **インストール確認**:
  ```bash
  docker --version
  docker-compose --version
  ```
  両方のバージョンが表示されればOK

- **インストール方法**:
  - [Docker Desktop](https://www.docker.com/products/docker-desktop/) をダウンロード・インストール
  - インストール後、Docker Desktopを起動しておく

### ローカル開発環境用（オプション）

Dockerを使わずにローカルで開発する場合は、以下も必要です：

#### 3. **Node.js**（JavaScriptランタイム）
- **役割**: フロントエンド（Next.js）の実行環境
- **推奨バージョン**: v18以上
- **インストール確認**:
  ```bash
  node --version
  npm --version
  ```
  `v18.x.x` 以上と表示されればOK

- **インストール方法**:
  - [Node.js公式サイト](https://nodejs.org/)からLTS版をダウンロード・インストール

#### 4. **Go**（プログラミング言語）
- **役割**: バックエンドAPIサーバーの実行環境
- **推奨バージョン**: 1.23以上
- **インストール確認**:
  ```bash
  go version
  ```
  `go version go1.23.x` 以上と表示されればOK

- **インストール方法**:
  - [Go公式サイト](https://go.dev/dl/)からダウンロード・インストール

#### 5. **PostgreSQL**（データベース）
- **役割**: ユーザー情報、スコア、問題などのデータ保存
- **推奨バージョン**: 14以上
- **インストール方法**:
  - **macOS**: `brew install postgresql@14`
  - **Windows**: [PostgreSQL公式サイト](https://www.postgresql.org/download/windows/)
  - **Linux**: `sudo apt-get install postgresql-14`

---

## リポジトリのクローン

### 1. GitHubからリポジトリをクローンする

```bash
# HTTPSでクローン（推奨）
git clone https://github.com/shiv-ko/CoT_game.git

# または、SSHでクローン（SSH鍵設定済みの場合）
git clone git@github.com:shiv-ko/CoT_game.git
```

**意味**: GitHubからソースコードを自分のPCにダウンロードします。

### 2. プロジェクトディレクトリに移動

```bash
cd CoT_game
```

**意味**: クローンしたプロジェクトのフォルダに移動します。

### 3. ディレクトリ構成を確認

```bash
ls -la
```

以下のようなファイル・フォルダが表示されればOKです：

```
.
├── README.md              # プロジェクト概要
├── setup.md               # このファイル
├── command.md             # 開発・運用コマンド集
├── docker-compose.yml     # Docker統合設定
├── package.json           # ルートのNode.js設定
├── backend/               # Goバックエンド
│   ├── Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── .env.example       # バックエンド環境変数テンプレート
├── frontend/              # Next.jsフロントエンド
│   ├── Dockerfile
│   ├── package.json
│   └── src/
├── db/                    # PostgreSQL初期化SQL
│   ├── init.sql
│   ├── migrations/
│   └── .env.example       # データベース環境変数テンプレート
└── docs/                  # 開発ドキュメント
```

---

## 環境変数の設定

このプロジェクトでは、APIキーやデータベース接続情報などを環境変数で管理しています。

### 1. バックエンド用の環境変数を設定

```bash
# backend/.env.exampleをコピーして.envファイルを作成
cp backend/.env.example backend/.env
```

**意味**: テンプレートファイルをコピーして、実際の設定ファイルを作ります。

### 2. backend/.envを編集

テキストエディタで `backend/.env` を開き、以下の項目を設定します：

```bash
# エディタで開く（例：Visual Studio Code）
code backend/.env

# または、viやnanoなどのコマンドラインエディタを使用
vi backend/.env
```

#### 必須項目:

```env
# データベース接続URL（Docker使用時はこのままでOK）
DATABASE_URL=postgres://postgres:postgres@db:5432/cot_game?sslmode=disable

# Gemini APIキー（後で取得）
GEMINI_API_KEY=__REPLACE_ME__

# 使用するAIモデル
AI_MODEL_NAME=gemini-1.5-flash
```

**注意**: 
- `GEMINI_API_KEY` は後で [Google AI Studio](https://makersuite.google.com/app/apikey) から取得して設定します
- Dockerを使う場合、`DATABASE_URL` のホスト名は `db` にしてください（Docker内部の名前解決のため）
- ローカル環境で動かす場合は `localhost` に変更してください

### 3. データベース用の環境変数を設定

```bash
# db/.env.exampleをコピーして.envファイルを作成
cp db/.env.example db/.env
```

### 4. db/.envを編集

```bash
code db/.env
```

以下のように設定します：

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=cot_game
```

**意味**: PostgreSQLコンテナの初期設定（ユーザー名、パスワード、データベース名）を定義します。

---

## Dockerを使った起動（推奨）

Docker Composeを使うと、フロントエンド・バックエンド・データベースを一括で起動できます。

### 1. Dockerコンテナをビルド・起動

```bash
docker-compose up --build
```

**このコマンドの意味**:
- `docker-compose`: Docker Composeコマンド
- `up`: コンテナを起動
- `--build`: Dockerイメージを最新のコードでビルドしてから起動

**初回は数分かかります**。以下のようなログが表示されます：

```
Creating network "cot_game_default" with the default driver
Creating volume "cot_game_postgres_data" with default driver
Building frontend
Building backend
Creating cot_game_db_1 ... done
Creating cot_game_backend_1 ... done
Creating cot_game_frontend_1 ... done
Attaching to cot_game_db_1, cot_game_backend_1, cot_game_frontend_1
```

### 2. 起動確認

以下のメッセージが表示されればOKです：

```
frontend_1  | Ready on http://localhost:3000
backend_1   | Server is running on port 8080
db_1        | database system is ready to accept connections
```

### 3. バックグラウンドで起動する場合

```bash
docker-compose up -d
```

**意味**: 
- `-d`: デタッチモード（バックグラウンド実行）

ログを確認する場合：
```bash
docker-compose logs -f
```

### 4. データベースマイグレーションの実行

初回起動時は自動でテーブルが作成されますが、追加のマイグレーションは手動で実行します。

```bash
# PostgreSQLコンテナに接続
docker-compose exec db psql -U postgres -d cot_game
```

**意味**:
- `docker-compose exec`: 実行中のコンテナ内でコマンドを実行
- `db`: データベースコンテナの名前
- `psql`: PostgreSQLのコマンドラインツール
- `-U postgres`: postgresユーザーで接続
- `-d cot_game`: cot_gameデータベースを指定

接続できたら、マイグレーションを適用：

```sql
-- マイグレーションファイルを実行
\i /docker-entrypoint-initdb.d/migrations/20251024_add_scores_columns.sql

-- テーブル構造を確認
\d scores

-- PostgreSQLから抜ける
\q
```

### 5. コンテナの停止

```bash
docker-compose down
```

**意味**: 起動中のコンテナを停止・削除します（データは保持されます）。

---

## ローカル開発環境での起動

Dockerを使わずに各サービスを個別に起動する方法です。

### 1. データベースの起動

PostgreSQLサービスを起動します：

```bash
# macOS (Homebrew)
brew services start postgresql@14

# Linux
sudo systemctl start postgresql

# Windows
# PostgreSQLサービスをサービス管理から起動
```

### 2. データベースとユーザーを作成

```bash
# PostgreSQLに接続
psql -U postgres

# データベースを作成
CREATE DATABASE cot_game;

# 抜ける
\q
```

### 3. 初期化SQLを実行

```bash
psql -U postgres -d cot_game -f db/init.sql
```

### 4. バックエンドの起動

```bash
# backendディレクトリに移動
cd backend

# 依存パッケージをダウンロード
go mod download

# サーバーを起動
go run main.go
```

**実行結果**:
```
Server is running on port 8080
```

**意味**:
- `go mod download`: go.modに記載された依存パッケージをダウンロード
- `go run main.go`: Goプログラムをコンパイル・実行

### 5. フロントエンドの起動（新しいターミナルで）

```bash
# frontendディレクトリに移動
cd frontend

# 依存パッケージをインストール
npm install

# 開発サーバーを起動
npm run dev
```

**実行結果**:
```
Ready on http://localhost:3000
```

**意味**:
- `npm install`: package.jsonに記載された依存パッケージをインストール
- `npm run dev`: Next.js開発サーバーを起動（ホットリロード対応）

---

## 動作確認

### 1. フロントエンドの確認

ブラウザで以下のURLにアクセス：

```
http://localhost:3000
```

Next.jsのウェルカムページまたはアプリのホーム画面が表示されればOKです。

### 2. バックエンドAPIの確認

ブラウザまたはcurlコマンドで以下にアクセス：

```bash
curl http://localhost:8080/api/questions
```

または、ブラウザで `http://localhost:8080/api/questions` を開きます。

JSON形式のレスポンスが返ってくればOKです。

### 3. データベース接続の確認

```bash
# Docker使用時
docker-compose exec db psql -U postgres -d cot_game -c "\dt"

# ローカル環境
psql -U postgres -d cot_game -c "\dt"
```

テーブル一覧が表示されればOKです：

```
              List of relations
 Schema |    Name    | Type  |  Owner   
--------+------------+-------+----------
 public | questions  | table | postgres
 public | scores     | table | postgres
 public | users      | table | postgres
```

---

## 開発用コマンド

### コードのフォーマット

```bash
# 全体をフォーマット
npm run format
```

**意味**: Prettierを使ってコード全体を整形します。

### リント（コード品質チェック）

```bash
# フロントエンドのリント
npm run lint:fe

# バックエンドのリント
npm run lint:go

# 両方実行
npm run lint
```

**意味**: ESLint（フロントエンド）とgolangci-lint（バックエンド）でコードの問題をチェックします。

### Goのフォーマット

```bash
cd backend

# gofumptでフォーマット
go run mvdan.cc/gofumpt@v0.6.0 -l -w .

# goimportsでインポート整理
go run golang.org/x/tools/cmd/goimports@latest -w .
```

### テストの実行

```bash
# バックエンドのテスト
cd backend
go test ./...

# フロントエンドのテスト（実装されている場合）
cd frontend
npm test
```

### Dockerコンテナのログ確認

```bash
# 全サービスのログを表示
docker-compose logs -f

# 特定のサービスのログのみ表示
docker-compose logs -f frontend
docker-compose logs -f backend
docker-compose logs -f db
```

**意味**: 
- `logs`: コンテナのログを表示
- `-f`: フォローモード（リアルタイムで新しいログを表示）

### Dockerコンテナの再起動

```bash
# 全サービスを再起動
docker-compose restart

# 特定のサービスのみ再起動
docker-compose restart backend
```

---

## トラブルシューティング

### 問題1: `docker-compose up` でエラーが出る

#### エラー: "port is already allocated"

**原因**: ポート番号が他のプロセスで使用されている

**解決方法**:
```bash
# ポート使用状況を確認
lsof -i :3000   # フロントエンド
lsof -i :8080   # バックエンド
lsof -i :5432   # PostgreSQL

# プロセスを停止するか、docker-compose.ymlのポート番号を変更
```

#### エラー: "permission denied"

**原因**: Dockerデーモンが起動していない、または権限がない

**解決方法**:
```bash
# Docker Desktopが起動しているか確認
# Linuxの場合、ユーザーをdockerグループに追加
sudo usermod -aG docker $USER
# ログアウト・ログインして再試行
```

### 問題2: データベースに接続できない

#### エラー: "connection refused"

**解決方法**:
```bash
# コンテナの状態を確認
docker-compose ps

# dbコンテナが起動していない場合
docker-compose up -d db

# ログを確認
docker-compose logs db
```

### 問題3: フロントエンドが表示されない

#### ブラウザに何も表示されない

**解決方法**:
```bash
# フロントエンドのログを確認
docker-compose logs frontend

# node_modulesを再インストール
docker-compose exec frontend rm -rf node_modules
docker-compose exec frontend npm install
docker-compose restart frontend
```

### 問題4: バックエンドAPIにアクセスできない

**解決方法**:
```bash
# バックエンドのログを確認
docker-compose logs backend

# .envファイルの設定を確認
cat backend/.env

# DATABASE_URLが正しいか確認
# Docker使用時は "db" をホスト名に使用
```

### 問題5: マイグレーションが失敗する

**解決方法**:
```bash
# データベースコンテナに入る
docker-compose exec db psql -U postgres -d cot_game

# 既存のテーブルを確認
\dt

# 必要に応じてテーブルを削除して再作成
DROP TABLE scores CASCADE;

# init.sqlを再実行（Docker再起動で自動実行されます）
docker-compose down
docker-compose up -d
```

### 問題6: `npm install` でエラーが出る

**解決方法**:
```bash
# npmのキャッシュをクリア
npm cache clean --force

# node_modulesとpackage-lock.jsonを削除して再インストール
rm -rf node_modules package-lock.json
npm install
```

### 問題7: Goのビルドでエラーが出る

**解決方法**:
```bash
cd backend

# go.modとgo.sumを整合
go mod tidy

# 依存パッケージを再ダウンロード
go mod download

# ビルドを試す
go build -o main .
```

---

## 次のステップ

セットアップが完了したら、以下のドキュメントを参照して開発を進めてください：

- **[README.md](./README.md)**: プロジェクト全体の概要
- **[docs/app.md](./docs/app.md)**: アプリの仕様と機能説明
- **[command.md](./command.md)**: よく使う開発・運用コマンド集
- **[docs/instruction/guideline_code.md](./docs/instruction/guideline_code.md)**: コーディング規約

---

## サポート

- 不明点があれば、GitHubのIssueで質問してください
- チーム内のSlackチャンネルでも質問できます

Happy Coding! 🚀
