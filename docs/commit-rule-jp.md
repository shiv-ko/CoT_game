# コミットルール（AI: Codex CLI向け指示）

このドキュメントは、AIアシスタント（Codex CLI）がこのリポジトリでコミットを作成する際に従うべき権威ある指示です。以下のガイドラインはすべて必須ルールです。

---

### 0) スコープと著者情報

- 著者情報: リポジトリのローカルgit設定（すでに人間の所有者に設定済み）を使用してください。
- 共同著者: AIが作成するすべてのコミットにはCodex CLIを共同著者として含めてください。
  - コミット本文の末尾に追加する共同著者トレーラー:
    - `Co-authored-by: Codex CLI <codex@example.com>`

  - すでにトレーラーが存在する場合は重複しないようにしてください。

---

### 1) メッセージスタイル（Conventional Commits / Angularスタイル）

- タイプ: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- フォーマット: `type(scope?): subject`
- サブジェクト: 命令形、英語、50文字以内（type/scope含む）。ピリオド不要。
- 本文: 英語、約72文字で改行。実装詳細ではなく「なぜ/何を」を説明。
- フッター（該当時）:
  - `BREAKING CHANGE: ...`
  - `Co-authored-by: Codex CLI <codex@example.com>`（AIコミットには必ず）

#### 例

```
feat(backend): /solveエンドポイントを追加

リクエストのバリデーションと評価者への接続を実装。
結果をスコアに保存し、レイテンシ指標も記録。

Co-authored-by: Codex CLI <codex@example.com>
```

```
docs: ディレクトリ構成付きREADMEを追加 (2025-09-07)

Co-authored-by: Codex CLI <codex@example.com>
```

---

### 2) AIによるコミット方法（必須）

1. `npm run format`でコードベースを整形する。
2. 意図したファイルのみステージする（無関係な変更は含めない）。
3. 英語で命令形、50文字以内のサブジェクトを作成する。
4. 動機や範囲が明確になる場合は短い本文を追加する。
5. 必ず共同著者トレーラー行を追加する。
6. Huskyフックが軽微な設定/ドキュメント変更で進行を妨げる場合、明示的な承認があれば`--no-verify`を使用可能。

#### コマンドテンプレート

```bash
npm run format
git add <files>
git commit -m "<type>: <subject>" \
           -m "<optional body>" \
           -m "Co-authored-by: Codex CLI <codex@example.com>"
```

---

### 3) コミットテンプレート（任意だが推奨）

リポジトリルートにコミットメッセージテンプレートを作成し、gitに設定：

```bash
cat > .gitmessage <<'MSG'
<type>(<optional scope>): <subject in English, <=50 chars>

<body, wrap at ~72 chars>

Co-authored-by: Codex CLI <codex@example.com>
MSG

git config --local commit.template .gitmessage
```

---

### 5) 言語要件

- コミットメッセージ: 英語のみ。
- コードコメント: プロジェクト方針により日本語も許容されるが、コミットメッセージは英語のみ。

---

### 6) AIがコミット前に確認すべき品質チェックリスト

- [ ] 意図したファイルのみステージされている
- [ ] サブジェクトが命令形・英語・50文字以内
- [ ] 本文（ある場合）は「なぜ/何を」を説明
- [ ] `Co-authored-by: Codex CLI <codex@example.com>`を含む
- [ ] 差分に秘密情報やトークンが含まれていない
