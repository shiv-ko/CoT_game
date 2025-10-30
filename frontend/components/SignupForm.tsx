'use client';

import React, { useState } from 'react';

import { signup } from '../services/authService';

import styles from './SignupForm.module.css';

/**
 * サインアップフォームコンポーネント
 */
const SignupForm: React.FC = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  /**
   * フォームの送信処理
   * @param e イベントオブジェクト
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      const response = await signup({ username, email, password });
      setSuccess(`Welcome, ${response.user.username}! Your account has been created.`);
      // TODO: トークンを保存し、ログイン状態に遷移する処理
      console.info('Token:', response.token);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      }
    }
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <h2 className={styles.formTitle}>アカウント作成</h2>
      {error && <div className={`${styles.message} ${styles.errorMessage}`}>{error}</div>}
      {success && <div className={`${styles.message} ${styles.successMessage}`}>{success}</div>}
      <div className={styles.formGroup}>
        <label htmlFor="username" className={styles.label}>
          ユーザー名
        </label>
        <input
          id="username"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          className={styles.input}
          placeholder="your-username"
        />
      </div>
      <div className={styles.formGroup}>
        <label htmlFor="email" className={styles.label}>
          メールアドレス
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className={styles.input}
          placeholder="your@email.com"
        />
      </div>
      <div className={styles.formGroup}>
        <label htmlFor="password" className={styles.label}>
          パスワード
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className={styles.input}
          placeholder="••••••••"
        />
      </div>
      <button type="submit" className={styles.submitButton}>
        サインアップ
      </button>
      <div className={styles.loginLink}>
        既にアカウントをお持ちですか? <a href="/login">ログイン</a>
      </div>
    </form>
  );
};

export default SignupForm;
