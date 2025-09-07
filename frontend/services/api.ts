import { ErrorResponse } from '../types/auth';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';

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

export default fetchApi;
