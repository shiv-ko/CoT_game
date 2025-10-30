'use client';

import { useParams, useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';

import { fetchQuestions, submitSolve } from '../../../../services/api';
import { Question } from '../../../../types/question';
import { SolveRequest, SolveResponse } from '../../../../types/solve';

import SolveResult from './components/SolveResult';
import styles from './QuestionDetail.module.css';

/**
 * 問題詳細ページコンポーネント
 * 指定されたIDの問題の詳細を表示し、ユーザーがプロンプトを入力して解答を提出できるようにします。
 */
export default function QuestionDetailPage() {
  const params = useParams();
  const router = useRouter();
  const questionId = Number(params.id);

  // 状態管理
  const [question, setQuestion] = useState<Question | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [prompt, setPrompt] = useState<string>('');
  const [modelName] = useState<string>('gemini-2.0-flash-lite'); // 固定値として使用
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [result, setResult] = useState<SolveResponse | null>(null);
  const [promptError, setPromptError] = useState<string | null>(null);

  /**
   * 問題詳細を取得する
   */
  useEffect(() => {
    const loadQuestion = async () => {
      try {
        setLoading(true);
        setError(null);
        const questions = await fetchQuestions();
        const found = questions.find((q) => q.id === questionId);

        if (!found) {
          setError('指定された問題が見つかりませんでした');
          return;
        }

        setQuestion(found);
      } catch (err) {
        setError(err instanceof Error ? err.message : '問題の取得に失敗しました');
      } finally {
        setLoading(false);
      }
    };

    if (questionId) {
      void loadQuestion();
    }
  }, [questionId]);

  /**
   * プロンプト入力のバリデーション
   * @param value - 入力されたプロンプト
   * @returns バリデーションエラーメッセージ(エラーがない場合はnull)
   */
  const validatePrompt = (value: string): string | null => {
    if (value.trim().length === 0) {
      return 'プロンプトを入力してください';
    }
    if (value.length > 2000) {
      return 'プロンプトは2000文字以内で入力してください';
    }
    return null;
  };

  /**
   * プロンプト変更時のハンドラ
   * @param e - 入力イベント
   */
  const handlePromptChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    setPrompt(value);
    setPromptError(validatePrompt(value));
  };

  /**
   * フォーム送信ハンドラ
   * @param e - フォーム送信イベント
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // バリデーション
    const validationError = validatePrompt(prompt);
    if (validationError) {
      setPromptError(validationError);
      return;
    }

    if (!question) {
      return;
    }

    try {
      setSubmitting(true);
      setPromptError(null);

      const request: SolveRequest = {
        question_id: questionId,
        prompt: prompt,
        model: modelName,
      };

      const response = await submitSolve(request);
      setResult(response);
    } catch (err) {
      setPromptError(err instanceof Error ? err.message : '解答の送信に失敗しました');
    } finally {
      setSubmitting(false);
    }
  };

  /**
   * 再挑戦ボタンのハンドラ
   * フォームをリセットして再度挑戦できるようにします
   */
  const handleRetry = () => {
    setResult(null);
    setPrompt('');
    setPromptError(null);
  };

  /**
   * 一覧に戻るボタンのハンドラ
   */
  const handleBackToList = () => {
    router.push('/questions');
  };

  // ローディング中
  if (loading) {
    return (
      <div className={styles.loadingState}>
        <p>読み込み中...</p>
      </div>
    );
  }

  // エラー発生時
  if (error || !question) {
    return (
      <div className={styles.errorState}>
        <p>{error || '問題が見つかりません'}</p>
        <button onClick={handleBackToList} className={styles.backButton}>
          一覧に戻る
        </button>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* ヘッダー */}
      <div className={styles.header}>
        <button onClick={handleBackToList} className={styles.backButton}>
          ← 一覧に戻る
        </button>
        <h1 className={styles.title}>問題 #{question.id}</h1>
        <div className={styles.metadata}>
          <span>難易度: {question.level}</span>
        </div>
        <p className={styles.description}>
          この問題に対して、AIに指示を出すプロンプトを入力してください。
          <br />
          あなたのプロンプトの質によって、AIの回答の正確さが変わります。
        </p>
      </div>

      {/* 結果が表示されていない場合のみフォームを表示 */}
      {!result && (
        <form onSubmit={handleSubmit} className={styles.form}>
          {/* プロンプト入力エリア */}
          <div className={styles.formGroup}>
            <label htmlFor="prompt" className={styles.label}>
              AIへのプロンプトを入力してください
            </label>
            <textarea
              id="prompt"
              value={prompt}
              onChange={handlePromptChange}
              placeholder="例: この問題を解いてください"
              rows={10}
              className={`${styles.textarea} ${promptError ? styles.textareaError : ''}`}
              disabled={submitting}
            />
            {promptError && <p className={styles.errorMessage}>{promptError}</p>}
            <p className={styles.charCount}>{prompt.length} / 2000 文字</p>
          </div>

          {/* 送信ボタン */}
          <button
            type="submit"
            disabled={submitting || !!promptError}
            className={`${styles.submitButton} ${submitting || promptError ? styles.submitButtonDisabled : ''}`}
          >
            {submitting ? '送信中...' : '解答を送信'}
          </button>
        </form>
      )}

      {/* 結果表示 */}
      {result && <SolveResult result={result} onRetry={handleRetry} />}
    </div>
  );
}
