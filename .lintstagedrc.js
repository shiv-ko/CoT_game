// .lintstagedrc.js
const path = require('path');

module.exports = {
  '*.{md,mdx,json,yml,yaml}': 'prettier --write',
  'docker-compose.yml': 'prettier --write',

  'frontend/**/*.{js,jsx,ts,tsx,css,scss}': files => {
    const relToFrontend = files.map(f => path.relative('frontend', f));
    return [
      `prettier --write ${files.join(' ')}`,
      `eslint --fix --cwd frontend ${relToFrontend.join(' ')}`
    ];
  },

  'backend/**/*.go': files => {
    // ルート相対 → backend 相対
    const rel = files.map(f => path.relative('backend', f));

    // 所属ディレクトリをユニークリスト化
    const dirs = Array.from(new Set(rel.map(f => path.dirname(f))));

    return [
      // gofumpt はファイル指定でOK
      `go run mvdan.cc/gofumpt@v0.6.0 -l -w ${files.join(' ')}`,

      // golangci-lint は「ディレクトリ」を渡す（複数OK）
      // 例: golangci-lint run --fix handlers routes
      `bash -lc "cd backend && golangci-lint run --fix ${dirs.join(' ')}"`
    ];
  }
};