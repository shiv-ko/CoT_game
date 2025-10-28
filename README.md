# CoT Game - プロンプトバトルWebアプリ

[![CI](https://github.com/shiv-ko/CoT_game/actions/workflows/ci.yaml/badge.svg)](https://github.com/shiv-ko/CoT_game/actions/workflows/ci.yaml)

---

## 🧭 基本構成

### 1. プロジェクト概要（Overview）

**CoT Game** は、ユーザーがプロンプトを工夫して AI に問題を解かせる対戦型 Web アプリケーションです。算数や論理パズルなど様々な問題に対して、最適なプロンプトを作成し、AI の回答精度を競います。

本リポジトリはモノレポ構成で、フロントエンド（Next.js）、バックエンド（Go）、データベース（PostgreSQL）、ドキュメントから構成されています。

**主な特徴:**

- AI（Gemini API）を活用した自動評価システム
- プロンプトの工夫による AI の回答精度向上を競うゲーム性
- リアルタイムでのスコア評価とランキング機能
- Docker による簡単な環境構築

---

### 2. 機能一覧（Features）

#### 🎯 コア機能

- **問題バンク**: 様々な難易度の問題を閲覧・選択
  - 難易度別フィルタリング
  - 問題の詳細表示
- **プロンプト送信と AI 評価**: ユーザーが作成したプロンプトを AI に送信し、自動評価
  - Gemini API との連携
  - 数値抽出と正解判定
  - スコア算出（100点満点）
  - レスポンスタイムの測定
- **スコア管理**: 解答履歴の保存と管理
  - AI モデル情報の記録
  - 抽出された数値と評価詳細の保存
  - 提出時刻とレイテンシの記録

#### 📊 今後実装予定の機能

- **ランキングシステム**: ユーザー間のスコア競争
  - 期間別ランキング（日次・週次・全期間）
  - ベストスコアとチャレンジ回数の表示
- **履歴表示**: 個人の解答履歴の確認
  - 問題別スコア履歴
  - プロンプトの振り返り

#### 🛡️ セキュリティ・運用機能

- **レートリミット**: API 乱用防止（実装予定）
- **構造化ログ**: リクエスト追跡とパフォーマンス分析（実装予定）
- **CI/CD**: 自動テスト・リント・ビルド検証

---

### 3. デモ / スクリーンショット（Optional）

<!-- スクリーンショットは今後追加予定 -->

_現在、スクリーンショットはありません。今後追加予定です。_

---

### 4. 使用技術（Tech Stack）

#### フロントエンド

- **[Next.js](https://nextjs.org/)** - React ベースのフルスタックフレームワーク
- **[TypeScript](https://www.typescriptlang.org/)** - 型安全な JavaScript
- **[Tailwind CSS](https://tailwindcss.com/)** - ユーティリティファーストの CSS フレームワーク
- **[ESLint](https://eslint.org/)** + **[Prettier](https://prettier.io/)** - コード品質とフォーマット管理

#### バックエンド

- **[Go](https://go.dev/)** - 高速で効率的なバックエンド言語
- **[Gin](https://gin-gonic.com/)** - Go の軽量 Web フレームワーク
- **[database/sql](https://pkg.go.dev/database/sql)** - Go 標準の SQL インターフェース
- **[Gemini API](https://ai.google.dev/)** - Google の生成 AI API

#### データベース

- **[PostgreSQL 14](https://www.postgresql.org/)** - オープンソースのリレーショナルデータベース

#### インフラ・開発環境

- **[Docker](https://www.docker.com/)** / **[Docker Compose](https://docs.docker.com/compose/)** - コンテナ化と開発環境の構築
- **[GitHub Actions](https://github.com/features/actions)** - CI/CD パイプライン
- **[Husky](https://typicode.github.io/husky/)** + **[lint-staged](https://github.com/okonet/lint-staged)** - Git フック管理

---

## 🚀 クイックスタート

### 前提条件

- Docker と Docker Compose がインストールされていること
- Gemini API キー（バックエンドの AI 機能を使用する場合）

### 起動手順

1. **リポジトリのクローン**

```bash
git clone https://github.com/shiv-ko/CoT_game.git
cd CoT_game
```

2. **環境変数の設定**

```bash
# バックエンド用の .env ファイルを作成
cp backend/.env.example backend/.env
# Gemini API キーを設定
# GEMINI_API_KEY=your_api_key_here

# データベース用の .env ファイルを作成
cp db/.env.example db/.env
```

3. **Docker コンテナの起動**

```bash
docker-compose up
```

4. **アクセス**

- フロントエンド: http://localhost:3000
- バックエンド API: http://localhost:8080
- PostgreSQL: localhost:5432

---

## 📁 ディレクトリ構成

```
.
├── backend/              # Go バックエンド
│   ├── handlers/         # HTTP ハンドラ
│   ├── internal/         # 内部パッケージ
│   │   ├── ai/          # AI クライアント（Gemini）
│   │   ├── eval/        # 評価ロジック
│   │   └── repository/  # データアクセス層
│   ├── models/          # データモデル
│   └── routes/          # ルーティング定義
├── frontend/            # Next.js フロントエンド
│   └── src/
│       └── app/         # App Router ページ
├── db/                  # データベース設定
│   ├── init.sql         # 初期スキーマ
│   └── migrations/      # マイグレーションファイル
├── docs/                # ドキュメント
└── docker-compose.yml   # Docker 構成
```

---

## 🛠️ 開発

### リント・フォーマット

```bash
# すべてのコードをフォーマット
npm run format

# フロントエンドのリント
npm run lint:fe

# バックエンドのリント
npm run lint:go

# 両方実行
npm run lint
```

### テスト

```bash
# バックエンドのテスト
cd backend
go test ./...

# 特定パッケージのテスト
go test ./internal/eval/...
```

### データベースマイグレーション

マイグレーションファイルは `db/migrations/` に配置されています。

**適用手順:**

```bash
# PostgreSQL コンテナに接続
docker-compose exec db psql -U postgres -d cot_game

# マイグレーションファイルを実行
\i /docker-entrypoint-initdb.d/migrations/20251024_add_scores_columns.sql

# 確認
\d scores
```

**利用可能なマイグレーション:**

- `20251024_add_scores_columns.sql`: `scores` テーブルに AI モデル情報、抽出数値、レイテンシ、評価詳細の各カラムを追加

---

## 📚 ドキュメント

- **[開発計画](docs/plan_0907.md)** - MVP に向けたスプリント計画
- **[タスクリスト](docs/task.md)** - 実装タスクの一覧と進捗
- **[進捗記録](docs/progresses/)** - 日次の開発進捗ログ
- **[セットアップガイド](setup.md)** - 詳細なセットアップ手順

---

## 🤝 コントリビューション

プルリクエストを歓迎します！以下のガイドラインに従ってください:

1. **コミットメッセージ**: [コミット規約](docs/instruction/commit-rule-jp.md) に従う
2. **ブランチ名**: [ブランチ・PR 命名規則](docs/instruction/branch-pr-name.md) に従う
3. **コーディング規約**: [コーディングガイドライン](docs/instruction/guideline_code.md) を参照

---

## 📝 ライセンス

ISC

---

## 📧 お問い合わせ

質問や提案がある場合は、Issue を作成してください。
