import { ErrorResponse } from '../types/auth';
import { Question } from '../types/question';
import { SolveRequest, SolveResponse } from '../types/solve';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

/**
 * APIリクエストを送信する汎用関数
 * @param path APIのエンドポイントパス
 * @param options リクエストオプション (method, bodyなど)
 * @returns レスポンスのJSONオブジェクト
 * @throws エラーレスポンスまたはネットワークエラー
 */
async function fetchApi<T>(path: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_BASE_URL}${path}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json();
      throw new Error(errorData.message || 'API request failed');
    }

    return response.json() as Promise<T>;
  } catch (error) {
    console.error('API Error:', error);
    if (error instanceof Error) {
      throw error;
    } else {
      throw new Error('An unknown error occurred');
    }
  }
}

/**
 * 問題一覧を取得する
 * @returns 問題の配列
 */
export async function fetchQuestions(): Promise<Question[]> {
  return fetchApi<Question[]>('/api/v1/questions');
}

/**
 * 問題を解くためのプロンプトをAIに送信し、評価結果を取得する
 * @param request - Solve APIリクエスト（問題ID、プロンプト、モデル）
 * @returns Solve APIレスポンス（AI応答、スコア、評価詳細）
 * @throws APIエラーまたはネットワークエラー
 */
export async function submitSolve(request: SolveRequest): Promise<SolveResponse> {
  return fetchApi<SolveResponse>('/api/v1/solve', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}

export default fetchApi;
