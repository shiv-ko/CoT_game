import fetchApi from './api';
import { User, LoginResponse } from '../types/auth';

/**
 * 新しいユーザーを登録する
 * @param userData サインアップ情報 (username, email, password)
 * @returns 登録されたユーザー情報とトークン
 */
export const signup = async (
  userData: Omit<User, 'id'> & { password: string },
): Promise<LoginResponse> => {
  return fetchApi<LoginResponse>('/api/signup', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

/**
 * ユーザーをログインさせる
 * @param credentials ログイン情報 (email, password)
 * @returns ログインしたユーザー情報とトークン
 */
export const login = async (
  credentials: Pick<User, 'email'> & { password: string },
): Promise<LoginResponse> => {
  return fetchApi<LoginResponse>('/api/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
};
