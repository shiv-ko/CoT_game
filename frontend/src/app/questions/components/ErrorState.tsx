import React from 'react';

interface ErrorStateProps {
  error: string;
  onRetry: () => void;
}

/**
 * エラー状態を表示するコンポーネント
 */
const ErrorState: React.FC<ErrorStateProps> = ({ error, onRetry }) => {
  return (
    <main style={{ padding: '2rem' }}>
      <h1>問題一覧</h1>
      <div style={{ color: 'red', marginBottom: '1rem' }}>
        <p>エラー: {error}</p>
      </div>
      <button onClick={onRetry}>再試行</button>
    </main>
  );
};

export default ErrorState;
