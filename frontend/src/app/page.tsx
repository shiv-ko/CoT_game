import Link from 'next/link';

import SignupForm from '../../components/SignupForm';

import styles from './Home.module.css';

export default function Home() {
  return (
    <main className={styles.container}>
      <div className={styles.hero}>
        <h1 className={styles.title}>Prompt Battle</h1>
        <p className={styles.subtitle}>AIを操る、究極のプロンプトバトル</p>
        <p className={styles.description}>
          問題を見ずに、最高のプロンプトを作成してAIに正解を導かせよう。
          <br />
          あなたのプロンプトエンジニアリング能力が試される!
        </p>
      </div>

      <div className={styles.formCard}>
        <SignupForm />
      </div>

      <div className={styles.features}>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>🎯</div>
          <h3 className={styles.featureTitle}>スマートな挑戦</h3>
          <p className={styles.featureDescription}>
            問題文を見ずに、AIへの指示だけで正解を導き出そう
          </p>
        </div>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>🏆</div>
          <h3 className={styles.featureTitle}>スコアで競争</h3>
          <p className={styles.featureDescription}>
            正確なプロンプトで高得点を獲得し、ランキング上位を目指そう
          </p>
        </div>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>🚀</div>
          <h3 className={styles.featureTitle}>スキル向上</h3>
          <p className={styles.featureDescription}>
            実践を通じてプロンプトエンジニアリング能力を磨こう
          </p>
        </div>
      </div>

      <Link href="/questions" className={styles.getStartedButton}>
        問題に挑戦する →
      </Link>
    </main>
  );
}
