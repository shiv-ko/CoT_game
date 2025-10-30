import Link from 'next/link';
import React from 'react';

import { Question } from '../../../../types/question';

import styles from './QuestionsTable.module.css';
import { renderStars } from './utils';

interface QuestionsTableProps {
  questions: Question[];
}

/**
 * 問題一覧テーブルコンポーネント
 * セキュリティとゲーム性の観点から、問題文は表示しません。
 * ユーザーは問題IDと難易度のみを見て、挑戦する問題を選びます。
 */
const QuestionsTable: React.FC<QuestionsTableProps> = ({ questions }) => {
  if (questions.length === 0) {
    return <p className={styles.emptyMessage}>選択した難易度の問題がありません</p>;
  }

  return (
    <div className={styles.container}>
      <p className={styles.count}>全 {questions.length} 問</p>
      <table className={styles.table}>
        <thead className={styles.tableHeader}>
          <tr>
            <th>難易度</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody className={styles.tableBody}>
          {questions.map((question) => (
            <tr key={question.id}>
              <td className={styles.levelCell}>{renderStars(question.level)}</td>
              <td className={styles.actionCell}>
                <Link href={`/questions/${question.id}`} className={styles.solveButton}>
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
