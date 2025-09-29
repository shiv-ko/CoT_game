# AI連携クライアント実装

## 目的 / 背景

- `docs/plan_0907.md` Sprint1: Gemini 連携クライアント層の実装。
- `docs/app.md` AI API 連携構成の中核。`POST /solve` 他バックエンド処理の前提。

## 完了条件

- `internal/ai` (予定) 配下に Gemini クライアント（インターフェース + 実装）ファイルを作成。
- 環境変数 `GEMINI_API_KEY` を読み込み、未設定時は起動エラーまたはハンドラで 500 + 明示メッセージ。
- タイムアウト（<= 15s）とリトライ（ネットワーク一時失敗時 1〜2回）実装。
- レスポンス構造体に `raw_text` を含め取得できる。
- 単体テスト（モックHTTPサーバー）で成功/4xx/5xx/タイムアウトの判定が通る。

## スコープ

- 含む: Gemini呼出（REST/SDKいずれか）、認証ヘッダ、HTTPクライアント設定。
- 含まない: OpenAI切替（別タスク: モデル切替抽象化）、利用料金集計。

## 作業手順

- [ ] `backend/internal/ai/` ディレクトリ作成
- [ ] `Client` インターフェース定義 (`Generate(ctx, prompt) (text, latency, error)`)
- [ ] Gemini 実装構造体 + コンストラクタ (`NewGeminiClient(apiKey string, httpClient *http.Client)`)
- [ ] タイムアウト付き `http.Client` 作成
- [ ] リトライロジック（`context.DeadlineExceeded` / 5xx）
- [ ] エラーハンドリング（分類: auth/network/api）
- [ ] 単体テスト (httptest.Server)
- [ ] README or inline comment に利用例追記

## 依存関係 / リスク

- リスク: API Key 未設定 → 明示的Fail Fast
- リスク: 外部API仕様変更 → エラーフォーマット吸収層を薄く持つ

## 見積り / 担当 / 期日

- 見積り: 1.5人日
- 担当: 未定
- 期日: Sprint1 前半 (2025-09-12 目安)

## 参考リンク

- `../app.md`
- `../plan_0907.md`
