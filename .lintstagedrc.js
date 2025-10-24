// .lintstagedrc.js
const path = require('path');

module.exports = {
  '*.{md,mdx,json,yml,yaml}': 'prettier --write',
  'docker-compose.yml': 'prettier --write',

  'frontend/**/*.{js,jsx,ts,tsx,css,scss}': (files) => {
    const relToFrontend = files.map((f) => path.relative('frontend', f));
    return [
      `prettier --write ${files.join(' ')}`,
      `eslint --fix --cwd frontend ${relToFrontend.join(' ')}`,
    ];
  },

  'backend/**/*.go': (files) => {
    // ルート相対 → backend 相対
    const rel = files.map((f) => path.relative('backend', f));

    // 所属ディレクトリをユニークリスト化
    const dirs = Array.from(new Set(rel.map((f) => path.dirname(f))));

    return [
      // 1) フォーマット（厳格）
      `go run mvdan.cc/gofumpt@v0.6.0 -l -w ${files.join(' ')}`,
      // 2) import 整理（並び替え・未使用削除）
      `go run golang.org/x/tools/cmd/goimports@latest -w ${files.join(' ')}`,
      // 3) Lint（型解決のため backend をカレントにしてディレクトリ単位）
      `bash -lc "cd backend && golangci-lint run --fix ${dirs.join(' ')}"`,
    ];
  },
};
