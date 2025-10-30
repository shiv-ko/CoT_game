/**
 * 問題の型定義（クライアント側）
 * セキュリティとゲーム性の観点から、problem_statementとcorrect_answerは含まれません。
 * ユーザーは問題文を見ずにプロンプトを工夫してAIに正解を導き出させます。
 */
export interface Question {
  id: number;
  level: number;
  created_at: string;
}
