## Commit Rules (Instructions for the AI: Codex CLI)

This document is the authoritative instruction the AI assistant (Codex CLI) must follow when creating commits in this repository. All guidance below is expressed as mandatory rules.

---

### 0) Scope and Identity

- Author identity: use the local git config of the repository (already set to the human owner).
- Co‑authoring: every commit created by the AI MUST include Codex CLI as a co‑author.
  - Co‑author trailer to append to every commit body:
    - `Co-authored-by: Codex CLI <codex@example.com>`
  - Do not duplicate the trailer if it already exists.

---

### 1) Message Style (Conventional Commits / Angular style)

- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.
- Format: `type(scope?): subject`
- Subject: imperative, English, 50 chars or less (including type/scope). No period.
- Body: English, wrap at ~72 chars. Explain the “why/what”, not implementation minutiae.
- Footers (when applicable):
  - `BREAKING CHANGE: ...`
  - `Co-authored-by: Codex CLI <codex@example.com>` (always present for AI commits)

Examples

```
feat(backend): add /solve endpoint

Implements request validation and wiring to evaluator. Stores results
into scores with latency metrics.

Co-authored-by: Codex CLI <codex@example.com>
```

```
docs: add README with directory structure (2025-09-07)

Co-authored-by: Codex CLI <codex@example.com>
```

---

### 2) How the AI should commit (step-by-step)

1. Stage only the intended files (no unrelated changes).
2. Compose a subject in English (<= 50 chars, imperative).
3. Add a short body when it clarifies motivation or scope.
4. Always append the co‑author trailer line shown above.
5. If Husky hooks would block progress for minor config/doc changes, use `--no-verify` after explicit approval.

Command template

```bash
git add <files>
git commit -m "<type>: <subject>" \
           -m "<optional body>" \
           -m "Co-authored-by: Codex CLI <codex@example.com>"
```

---

### 3) Commit Template (optional but recommended)

Create a commit message template at the repo root and configure git.

```bash
cat > .gitmessage <<'MSG'
<type>(<optional scope>): <subject in English, <=50 chars>

<body, wrap at ~72 chars>

Co-authored-by: Codex CLI <codex@example.com>
MSG

git config --local commit.template .gitmessage
```

---

### 4) Bulk co‑authoring for existing history (advanced)

Warning: rewrites history; coordinate with collaborators. Force‑push required.

```bash
git branch backup/coauthor-$(date +%Y%m%d)
export FILTER_BRANCH_SQUELCH_WARNING=1
git filter-branch -f --msg-filter '
cat
printf "\nCo-authored-by: Codex CLI <codex@example.com>\n"
' --tag-name-filter cat -- --branches --tags
git push --force-with-lease origin main --tags
```

De‑duplication: if you need to avoid duplicate trailers, add a small sed/grep guard to the `--msg-filter` script.

---

### 5) Language requirements

- Commit messages: English only.
- Code comments: Japanese is acceptable for in‑code documentation per project policy; however messages remain English.

---

### 6) Quality checklist for the AI before committing

- [ ] Only the intended files are staged
- [ ] Subject is imperative, English, <= 50 chars
- [ ] Body (if present) explains why/what
- [ ] Includes `Co-authored-by: Codex CLI <codex@example.com>`
- [ ] No secrets or tokens are present in the diff

---

## 共同コミット（Co-authored-by）

複数人による共同作業を明示するため、コミット本文の末尾に共同作者トレーラーを追加します。

### 基本ルール

- 形式: `Co-authored-by: 氏名 <メールアドレス>`
- 複数人可。1行につき1名、空行を挟まずに連続で記載
- 件名（サブジェクト）と本文の間、および本文とトレーラーの間に空行を1つ入れる
- 本ドキュメントのコミットメッセージ規約（Angularプレフィックス、日本語、50文字以内）を遵守

### 新規コミットで共同作者を付与する（推奨）

```bash
git commit -m "feat: ランキングAPIを追加" \
           -m "実装詳細を本文に記載" \
           -m "Co-authored-by: Codex CLI <codex@example.com>"
```

- 共同作者が複数なら `-m` を追加して複数行記載します。

### コミットテンプレートを使う

1. ルートに `.gitmessage` を作成し、末尾に共同作者行を入れておく
2. `git config --local commit.template .gitmessage`

テンプレート例（先頭1行は件名。2行目は空行。以降本文。末尾にトレーラー）:

```
<type>: <短く要点を日本語で（50文字以内）>

本文（任意）

Co-authored-by: Codex CLI <codex@example.com>
```
