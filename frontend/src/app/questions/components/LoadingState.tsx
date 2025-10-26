import React from 'react';

/**
 * ローディング状態を表示するコンポーネント
 */
const LoadingState: React.FC = () => {
  return (
    <main style={{ padding: '2rem' }}>
      <h1>問題一覧</h1>
      <p>読み込み中...</p>
    </main>
  );
};

export default LoadingState;
