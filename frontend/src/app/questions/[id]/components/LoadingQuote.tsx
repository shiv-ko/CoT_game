'use client';

import React from 'react';

import styles from './LoadingQuote.module.css';

export const LoadingQuote: React.FC = () => {
  return (
    <div className={styles.loadingOverlay}>
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <div className={styles.quoteContainer}>
          <p className={styles.quote}>忍耐こそが、力よりも大きな成果を生む。</p>
          <p className={styles.author}>by エドマンド・バーク（Edmund Burke）</p>
        </div>
        <p className={styles.message}>解答を評価しています...</p>
      </div>
    </div>
  );
};
