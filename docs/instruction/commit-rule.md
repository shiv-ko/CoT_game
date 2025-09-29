# Commit Rules (Instructions for the AI: Codex CLI)

This document is the authoritative instruction the AI assistant (Codex CLI) must follow when creating commits in this repository. All guidance below is expressed as mandatory rules.

---

### 0) Scope and Identity

- Author identity: use the local git config of the repository (already set to the human owner).
- Co-authoring: every commit created by the AI MUST include Codex CLI as a co-author.
  - Co-author trailer to append to every commit body:
    - `Co-authored-by: Codex CLI <codex@example.com>`

  - Do not duplicate the trailer if it already exists.

---

### 1) Message Style (Conventional Commits / Angular style)

- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.
- Format: `type(scope?): subject`
- Subject: imperative, English, 50 chars or less (including type/scope). No period.
- Body: English, wrap at \~72 chars. Explain the “why/what”, not implementation minutiae.
- Footers (when applicable):
  - `BREAKING CHANGE: ...`
  - `Co-authored-by: Codex CLI <codex@example.com>` (always present for AI commits)

#### Examples

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

### 2) How the AI should commit (mandatory)

1. Run `npm run format` to format the codebase.
2. Stage only the intended files (no unrelated changes).
3. Compose a subject in English, imperative, <= 50 chars.
4. Add a short body when it clarifies motivation or scope.
5. Always append the co-author trailer line.
6. If Husky hooks block progress for minor config/doc changes, `--no-verify` may be used with explicit approval.

#### Command template

```bash
npm run format
git add <files>
git commit -m "<type>: <subject>" \
           -m "<optional body>" \
           -m "Co-authored-by: Codex CLI <codex@example.com>"
```

---

### 3) Commit Template (optional but recommended)

Create a commit message template at the repo root and configure git:

```bash
cat > .gitmessage <<'MSG'
<type>(<optional scope>): <subject in English, <=50 chars>

<body, wrap at ~72 chars>

Co-authored-by: Codex CLI <codex@example.com>
MSG

git config --local commit.template .gitmessage
```

---

### 5) Language requirements

- Commit messages: English only.
- Code comments: Japanese is acceptable for in-code documentation per project policy; commit messages remain English.

---

### 6) Quality checklist for the AI before committing

- [ ] Only the intended files are staged
- [ ] Subject is imperative, English, <= 50 chars
- [ ] Body (if present) explains why/what
- [ ] Includes `Co-authored-by: Codex CLI <codex@example.com>`
- [ ] No secrets or tokens are present in the diff
