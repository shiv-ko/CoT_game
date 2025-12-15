CI導入手順（Docs / Lighthouse / Vulnerability）

## 0. 前提（この手順が想定する構成）

- リポジトリがモノレポで、例として以下の構成を想定：
  - `frontend/` … Next.js
  - `backend/` … Go API
  - `openapi/` … OpenAPI仕様（`openapi.yaml`）
- CIは **GitHub Actions**
- PostgreSQLは **Lighthouse計測では不要** （ページ表示だけ測るため）。
  DBが必須のページを計測したい場合は後述の拡張を参照。

---

## 1. 自動ドキュメント生成（OpenAPI → 静的HTML）

### 1-1. OpenAPIファイルを用意

`openapi/openapi.yaml` を作成（既にあるならOK）。

推奨：**Spec-first**で「このYAMLが正」として管理する（コメント依存を避けたい場合）。

### 1-2. CIでHTMLを生成してArtifact保存

`.github/workflows/docs-openapi.yml` を作成：

```yaml
name: docs-openapi-html

on:
  pull_request:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read

jobs:
  build-openapi-html:
    runs-on: ubuntu-latest
    defaults:
      run:
        shell: bash

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      # Redocly CLI を npx 経由で使って単一HTMLを生成
      - name: Build OpenAPI HTML (Redocly)
        run: |
          mkdir -p dist/api-docs
          npx -y @redocly/cli build-docs openapi/openapi.yaml -o dist/api-docs/index.html

      - name: Upload artifact (API Docs HTML)
        uses: actions/upload-artifact@v4
        with:
          name: api-docs-html
          path: dist/api-docs
```

**レポート用に取れる証拠**

- Actionsのログ（`Build OpenAPI HTML` が成功）
- `api-docs-html` artifact が生成されているスクショ

---

## 2. Lighthouse パフォーマンス測定（Next.js）

### 2-1. Lighthouse CI設定ファイルを追加

リポジトリ直下に `.lighthouserc.json` を作成（URLは計測したいページに合わせて変更）：

```json
{
  "ci": {
    "collect": {
      "numberOfRuns": 2,
      "startServerCommand": "pnpm --dir frontend start -p 3000",
      "url": ["http://localhost:3000/", "http://localhost:3000/about"]
    },
    "assert": {
      "assertions": {
        "categories:performance": ["warn", { "minScore": 0.6 }],
        "categories:accessibility": ["warn", { "minScore": 0.7 }],
        "categories:seo": ["warn", { "minScore": 0.7 }]
      }
    },
    "upload": {
      "target": "temporary-public-storage"
    }
  }
}
```

- **最初は `warn`** にして、落ちないようにするのが安定（運用開始後に `error` へ）
- スコアは環境ブレがあるので、最初はゆるめでOK

### 2-2. CIワークフロー追加

`.github/workflows/lighthouse.yml` を作成：

```yaml
name: lighthouse

on:
  pull_request:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read

jobs:
  lhci:
    runs-on: ubuntu-latest
    defaults:
      run:
        shell: bash

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Enable Corepack (pnpm)
        run: corepack enable

      - name: Install frontend deps
        run: pnpm --dir frontend install --frozen-lockfile

      - name: Build frontend
        run: pnpm --dir frontend build

      - name: Install Lighthouse CI
        run: npm i -g @lhci/cli@0.14.x

      - name: Run Lighthouse CI (LHCI)
        env:
          # startServerCommand で pnpm を使うので PATH に入っている必要がある
          # corepack enable 済みならOK
          NODE_ENV: production
        run: |
          lhci autorun --config=.lighthouserc.json

      # レポートを成果物として残したい場合（任意）
      - name: Upload Lighthouse results (optional)
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: lighthouse-results
          path: .lighthouseci
```

**レポート用に取れる証拠**

- Actionsログで `lhci autorun` が完走した画面
- `temporary-public-storage` のURL（ログに出る）や `.lighthouseci` artifact のスクショ

---

## 3. 脆弱性診断（Trivyでまとめて簡単に）

### 3-1. CIワークフロー追加

`.github/workflows/security-trivy.yml` を作成：

```yaml
name: security-trivy

on:
  pull_request:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read

jobs:
  trivy-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # リポジトリ全体（Go/Node依存、設定、秘密情報っぽいもの）をまとめてスキャン
      - name: Trivy filesystem scan (vuln + misconfig + secret)
        uses: aquasecurity/trivy-action@0.28.0
        with:
          scan-type: fs
          scan-ref: .
          format: table
          severity: HIGH,CRITICAL
          ignore-unfixed: true
          vuln-type: os,library

      # 失敗させたい場合（課題で「検出できた」を示したいなら最初はfailしないのもアリ）
      # - fail-build: true 相当の制御を厳密にしたい場合は後で調整
```

**運用のコツ**

- 最初は「落とさない（failしない）」で導入 → 影響を確認
- 次に「CRITICALだけ落とす」など段階的に厳しくする

**レポート用に取れる証拠**

- Trivyの出力（検出0でも“チェックが走った”証拠になる）
- もし検出が出たら、その内容（例：依存パッケージの脆弱性）を利点として書ける

---

## 4. 3つのCIが揃ったか確認するチェックリスト

- [ ] PRを作ったら `docs-openapi-html` が走る
- [ ] PRを作ったら `lighthouse` が走る
- [ ] PRを作ったら `security-trivy` が走る
- [ ] Actionsの各ジョブで、成功ログ or Artifact が残る（スクショ可能）

---

## 5. よくある詰まりポイント（先に潰す）

### Lighthouseが落ちる

- `pnpm start` が起動できてない：`frontend build` が成功しているか確認
- 計測URLが存在しない：`/about` 等を実在ページに合わせる
- スコアブレ：最初は `warn` + `minScore`低めにする（0.6〜）

### OpenAPI HTMLが生成できない

- `openapi.yaml` がパス通りに存在するか
- YAMLの構文エラー（ここはCIログで見える）

### Trivyがノイズ多い

- まずは `HIGH,CRITICAL` のみに絞る（上の設定通り）
- `ignore-unfixed: true` で未修正案件を除外し、課題の初期導入を楽にする
