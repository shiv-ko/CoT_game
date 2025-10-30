'use client';

import React, { useEffect, useMemo, useState } from 'react';

import { fetchQuestions } from '../../../services/api';
import { Question } from '../../../types/question';

import EmptyState from './components/EmptyState';
import ErrorState from './components/ErrorState';
import LevelFilter from './components/LevelFilter';
import LoadingState from './components/LoadingState';
import QuestionsTable from './components/QuestionsTable';
import styles from './QuestionsPage.module.css';

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
   * レベルでフィルタリング(レベル順にソート)
   */
  useEffect(() => {
    let filtered: Question[];
    if (selectedLevel === 'all') {
      filtered = questions;
    } else {
      filtered = questions.filter((q) => q.level === selectedLevel);
    }
    // レベル順にソート(簡単な順)
    setFilteredQuestions([...filtered].sort((a, b) => a.level - b.level));
  }, [selectedLevel, questions]);

  /**
   * リトライボタンのハンドラ
   */
  const handleRetry = () => {
    loadQuestions();
  };

  // ユニークなレベルを取得(メモ化)
  const uniqueLevels = useMemo(
    () => Array.from(new Set(questions.map((q) => q.level))).sort((a, b) => a - b),
    [questions],
  );

  // ローディング中
  if (loading) {
    return <LoadingState />;
  }

  // エラー時
  if (error) {
    return <ErrorState error={error} onRetry={handleRetry} />;
  }

  // 問題が空の場合
  if (questions.length === 0) {
    return <EmptyState />;
  }

  return (
    <main className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>問題一覧</h1>
      </div>

      <div className={styles.statsBar}>
        <div className={styles.statItem}>
          <span className={styles.statValue}>{questions.length}</span>
          <span className={styles.statLabel}>総問題数</span>
        </div>
        <div className={styles.statItem}>
          <span className={styles.statValue}>{uniqueLevels.length}</span>
          <span className={styles.statLabel}>難易度レベル</span>
        </div>
        <div className={styles.statItem}>
          <span className={styles.statValue}>{filteredQuestions.length}</span>
          <span className={styles.statLabel}>表示中</span>
        </div>
      </div>

      <div className={styles.filterSection}>
        <LevelFilter
          selectedLevel={selectedLevel}
          uniqueLevels={uniqueLevels}
          onLevelChange={setSelectedLevel}
        />
      </div>

      <div className={styles.content}>
        <QuestionsTable questions={filteredQuestions} />
      </div>
    </main>
  );
};

export default QuestionsPage;
