# .env.example 整備

## 目的 / 背景

- 環境変数の明示でオンボーディング容易化 (`docs/progresses/progress_0908.md` 課題)。

## 完了条件

- ルートに `.env.example` 作成。
- 項目: `DATABASE_URL`, `GEMINI_API_KEY`, `RATE_LIMIT_PER_MINUTE`, `API_BASE_URL` (FE), `NODE_ENV`。
- README にコピー手順追記（後続 ドキュメント更新タスク）。

## スコープ

- 含む: .env.example 作成。
- 含まない: 秘密値生成スクリプト。

## 作業手順

- [ ] 変数一覧整理
- [ ] `.env.example` 追加
- [ ] 既存 `.gitignore` 確認

## 依存関係 / リスク

- 依存: AI連携クライアント実装 (必要キー確定)
- リスク: 変数増殖 → コメントで用途明示

## 見積り / 担当 / 期日

- 見積り: 0.25人日
- 担当: 未定
- 期日: Sprint2 (2025-09-18)

## 参考リンク

- `../progresses/progress_0908.md`
