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
    // ルート相対 → backend 相対に直す（golangci-lint は backend で実行）
    const relToBackend = files.map((f) => path.relative('backend', f)).join(' ');
    return [
      // gofumpt はルート実行でOK（ファイルパスを渡す）
      `go run mvdan.cc/gofumpt@v0.6.0 -l -w ${files.join(' ')}`,
      // モジュールルート backend で実行（型解決OK）
      `bash -lc "cd backend && golangci-lint run --fix ${relToBackend}"`,
    ];
  },
};
