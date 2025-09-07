/**
 * ユーザー情報を表す型
 */
export interface User {
  id: number;
  username: string;
  email: string;
}

/**
 * ログインAPIのレスポンスの型
 */
export interface LoginResponse {
  token: string;
  user: User;
}

/**
 * エラーレスポンスの型
 */
export interface ErrorResponse {
  message: string;
}
