## Git Commit Messages

When generating git commit messages:

- **ALWAYS use a prefix** (feat, fix, docs, style, refactor, perf, test, chore)
- Keep the main message under 50 characters INCLUDING the prefix
- Use imperative mood after the prefix
- Be specific but concise
- No periods at the end
- Format: `prefix: message`
- **Language**: Always use Japanese for commit messages

### Prefix Rules (Angular Convention):

- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation only changes
- `style:` Changes that do not affect the meaning of the code (formatting, missing semi-colons, etc)
- `refactor:` A code change that neither fixes a bug nor adds a feature
- `perf:` A code change that improves performance
- `test:` Adding missing or correcting existing tests
- `chore:` Changes to the build process or auxiliary tools and libraries

### Code Comments Language:

Use Japanese comments

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
1) ルートに `.gitmessage` を作成し、末尾に共同作者行を入れておく
2) `git config --local commit.template .gitmessage`

テンプレート例（先頭1行は件名。2行目は空行。以降本文。末尾にトレーラー）:
```
<type>: <短く要点を日本語で（50文字以内）>

本文（任意）

Co-authored-by: Codex CLI <codex@example.com>
```

### 既存履歴に一括で共同作者を付与する（高度）
注意: 履歴を書き換えるため、共有ブランチでは周知・合意の上で実施し、完了後は force-push が必要です。

1) 念のためバックアップブランチを作成
```bash
git branch backup/coauthor-$(date +%Y%m%d)
```
2) 共同作者トレーラーを全コミット本文末尾に付与（filter-branch 版）
```bash
export FILTER_BRANCH_SQUELCH_WARNING=1
git filter-branch -f --msg-filter '
cat
printf "\nCo-authored-by: Codex CLI <codex@example.com>\n"
' --tag-name-filter cat -- --branches --tags
```
3) 強制プッシュ（安全策として --force-with-lease を使用）
```bash
git push --force-with-lease origin main --tags
```

補足:
- `git filter-repo` が利用可能なら、そちらの使用を推奨します（高速・安全）。導入後、同様に本文末尾へトレーラーを付与してください。
- トレーラーは既に存在する場合は重複させないこと（フィルタで重複行を除去する処理を追加）。

### 共同作者情報のガイド
- 氏名・メールは各自の GitHub 公開メール（または noreply メール）を推奨
- 共同作者の同意を得た上で付与
- 履歴書き換え後は全共同者に `git fetch --all --prune` と再ベース操作の周知を行う
