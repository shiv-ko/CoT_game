import React from 'react';

/**
 * 問題が存在しない状態を表示するコンポーネント
 */
const EmptyState: React.FC = () => {
  return (
    <main style={{ padding: '2rem' }}>
      <h1>問題一覧</h1>
      <p>問題がありません</p>
    </main>
  );
};

export default EmptyState;
