import React from 'react';

import { SolveResponse } from '../../../../../types/solve';

import styles from './SolveResult.module.css';

/**
 * SolveResult コンポーネントのプロパティ
 */
interface SolveResultProps {
  /** Solve APIのレスポンスデータ */
  result: SolveResponse;
  /** 再挑戦ボタンがクリックされた時のコールバック関数 */
  onRetry: () => void;
}

/**
 * スコアに応じたCSSクラス名を取得する関数
 * @param score - 評価スコア(0-100)
 * @returns CSSモジュールのクラス名
 */
function getScoreClassName(score: number): string {
  if (score >= 90) {
    return styles.scoreExcellent;
  }
  if (score >= 70) {
    return styles.scoreGood;
  }
  if (score >= 50) {
    return styles.scoreFair;
  }
  return styles.scorePoor;
}

/**
 * スコアに応じたメッセージを取得する関数
 * @param score - 評価スコア(0-100)
 * @returns スコアに対する評価メッセージ
 */
function getScoreMessage(score: number): string {
  if (score >= 90) {
    return '素晴らしい!';
  }
  if (score >= 70) {
    return 'もう少しです!';
  }
  if (score >= 50) {
    return '頑張りましょう';
  }
  return '再挑戦してみましょう';
}

/**
 * Solve結果表示コンポーネント
 * AIの応答、評価スコア、評価詳細を視覚的に表示します。
 */
const SolveResult: React.FC<SolveResultProps> = ({ result, onRetry }) => {
  const scoreClassName = getScoreClassName(result.score);
  const scoreMessage = getScoreMessage(result.score);

  return (
    <div className={styles.container}>
      {/* ヘッダー: スコア表示 */}
      <div className={styles.scoreSection}>
        <div className={styles.scoreLabel}>あなたのスコア</div>
        <div className={`${styles.scoreValue} ${scoreClassName}`}>
          {result.score}
          <span className={styles.scoreUnit}>点</span>
        </div>
        <div className={styles.scoreMessage}>{scoreMessage}</div>
      </div>

      {/* AI応答の表示 */}
      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>AIの回答</h3>
        <div className={styles.aiOutput}>{result.ai_output}</div>
      </div>

      {/* 詳細情報 */}
      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>詳細情報</h3>
        <div className={styles.detailsGrid}>
          <div className={styles.detailItem}>
            <strong>モデル:</strong> {result.model_vendor} ({result.model_name})
          </div>
          <div className={styles.detailItem}>
            <strong>応答時間:</strong> {result.elapsed_ms}ms
          </div>
          {result.answer_number !== null && (
            <div className={styles.detailItem}>
              <strong>抽出された数値:</strong> {result.answer_number}
            </div>
          )}
          <div className={styles.detailItem}>
            <strong>評価モード:</strong> {result.evaluation.mode || 'N/A'}
          </div>
          <div className={styles.detailItem}>
            <strong>保存状態:</strong> {result.saved ? '✓ 保存済み' : '✗ 未保存'}
          </div>
        </div>
      </div>

      {/* 再挑戦ボタン */}
      <div className={styles.retrySection}>
        <button onClick={onRetry} className={styles.retryButton}>
          もう一度挑戦する
        </button>
      </div>
    </div>
  );
};

export default SolveResult;
