'use client';

import Link from 'next/link';
import React, { useEffect, useMemo, useState } from 'react';

import { fetchQuestions } from '../../../services/api';
import { Question } from '../../../types/question';

/**
 * 問題一覧ページコンポーネント
 * 問題の取得、フィルタリング、表示を管理します
 */
const QuestionsPage: React.FC = () => {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [filteredQuestions, setFilteredQuestions] = useState<Question[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedLevel, setSelectedLevel] = useState<number | 'all'>('all');

  /**
   * 問題一覧を取得する
   */
  const loadQuestions = async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await fetchQuestions();
      setQuestions(data);
      setFilteredQuestions(data);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('問題の取得に失敗しました');
      }
    } finally {
      setLoading(false);
    }
  };

  /**
   * 初回マウント時に問題を取得
   */
  useEffect(() => {
    loadQuestions();
  }, []);

  /**
   * レベルでフィルタリング（レベル順にソート）
   */
  useEffect(() => {
    let filtered: Question[];
    if (selectedLevel === 'all') {
      filtered = questions;
    } else {
      filtered = questions.filter((q) => q.level === selectedLevel);
    }
    // レベル順にソート（簡単な順）
    setFilteredQuestions([...filtered].sort((a, b) => a.level - b.level));
  }, [selectedLevel, questions]);

  /**
   * レベルを星で表示
   * @param level 問題のレベル (1-5)
   * @returns 星の文字列 (例: ★★★☆☆)
   */
  const renderStars = (level: number): string => {
    const maxStars = 5;
    const filledStars = '★'.repeat(level);
    const emptyStars = '☆'.repeat(maxStars - level);
    return filledStars + emptyStars;
  };

  /**
   * リトライボタンのハンドラ
   */
  const handleRetry = () => {
    loadQuestions();
  };

  // ユニークなレベルを取得（メモ化）
  const uniqueLevels = useMemo(
    () => Array.from(new Set(questions.map((q) => q.level))).sort((a, b) => a - b),
    [questions],
  );

  // ローディング中
  if (loading) {
    return (
      <main style={{ padding: '2rem' }}>
        <h1>問題一覧</h1>
        <p>読み込み中...</p>
      </main>
    );
  }

  // エラー時
  if (error) {
    return (
      <main style={{ padding: '2rem' }}>
        <h1>問題一覧</h1>
        <div style={{ color: 'red', marginBottom: '1rem' }}>
          <p>エラー: {error}</p>
        </div>
        <button onClick={handleRetry}>再試行</button>
      </main>
    );
  }

  // 問題が空の場合
  if (questions.length === 0) {
    return (
      <main style={{ padding: '2rem' }}>
        <h1>問題一覧</h1>
        <p>問題がありません</p>
      </main>
    );
  }

  return (
    <main style={{ padding: '2rem' }}>
      <h1>問題一覧</h1>

      {/* 難易度フィルタ */}
      <div style={{ marginBottom: '1.5rem' }}>
        <label htmlFor="level-filter" style={{ marginRight: '0.5rem' }}>
          難易度フィルタ:
        </label>
        <select
          id="level-filter"
          value={selectedLevel}
          onChange={(e) =>
            setSelectedLevel(e.target.value === 'all' ? 'all' : Number(e.target.value))
          }
          style={{ padding: '0.5rem' }}
        >
          <option value="all">すべて</option>
          {uniqueLevels.map((level) => (
            <option key={level} value={level}>
              レベル {level} ({renderStars(level)})
            </option>
          ))}
        </select>
      </div>

      {/* 問題一覧 */}
      {filteredQuestions.length === 0 ? (
        <p>選択した難易度の問題がありません</p>
      ) : (
        <div>
          <p>全 {filteredQuestions.length} 問</p>
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
              {filteredQuestions.map((question) => (
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
      )}
    </main>
  );
};

export default QuestionsPage;
