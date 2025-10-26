import Link from 'next/link';
import React from 'react';

import { Question } from '../../../../types/question';

import { renderStars } from './utils';

interface QuestionsTableProps {
  questions: Question[];
}

/**
 * 問題一覧テーブルコンポーネント
 */
const QuestionsTable: React.FC<QuestionsTableProps> = ({ questions }) => {
  if (questions.length === 0) {
    return <p>選択した難易度の問題がありません</p>;
  }

  return (
    <div>
      <p>全 {questions.length} 問</p>
      <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
        <thead>
          <tr style={{ borderBottom: '2px solid #ddd' }}>
            <th style={{ padding: '0.75rem', textAlign: 'left' }}>ID</th>
            <th style={{ padding: '0.75rem', textAlign: 'left' }}>難易度</th>
            <th style={{ padding: '0.75rem', textAlign: 'left' }}>問題文</th>
            <th style={{ padding: '0.75rem', textAlign: 'center' }}>操作</th>
          </tr>
        </thead>
        <tbody>
          {questions.map((question) => (
            <tr key={question.id} style={{ borderBottom: '1px solid #eee' }}>
              <td style={{ padding: '0.75rem' }}>{question.id}</td>
              <td style={{ padding: '0.75rem' }}>{renderStars(question.level)}</td>
              <td style={{ padding: '0.75rem' }}>{question.problem_statement}</td>
              <td style={{ padding: '0.75rem', textAlign: 'center' }}>
                <Link
                  href={`/questions/${question.id}`}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#0070f3',
                    color: 'white',
                    textDecoration: 'none',
                    borderRadius: '4px',
                    display: 'inline-block',
                  }}
                >
                  解く
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default QuestionsTable;
