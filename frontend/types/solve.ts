/**
 * Solve API関連の型定義
 * POST /api/v1/solve のリクエストとレスポンスの型を定義します。
 */

/**
 * Solve APIリクエストの型
 * ユーザーが問題を解くために送信するプロンプトとモデル情報
 */
export interface SolveRequest {
  /** 問題のID */
  question_id: number;
  /** ユーザーが入力したプロンプト（AIへの指示） */
  prompt: string;
  /** 使用するAIモデル（デフォルト: "gemini-1.5-flash"） */
  model?: string;
}

/**
 * Solve APIレスポンスの型
 * AI応答の結果と評価スコアを含む
 */
export interface SolveResponse {
  /** 問題のID */
  question_id: number;
  /** 送信されたプロンプト */
  prompt: string;
  /** AIモデルのベンダー名（例: "gemini"） */
  model_vendor: string;
  /** AIモデル名（例: "gemini-1.5-flash"） */
  model_name: string;
  /** AIが生成した回答テキスト */
  ai_output: string;
  /** AIの回答から抽出された数値（数値問題の場合） */
  answer_number: number | null;
  /** 評価スコア（0-100点） */
  score: number;
  /** 評価の詳細情報（評価モード、一致状況など） */
  evaluation: {
    /** 評価モード（例: "exact", "fuzzy"） */
    mode?: string;
    /** その他の評価詳細 */
    [key: string]: unknown;
  };
  /** AI応答までの経過時間（ミリ秒） */
  elapsed_ms: number;
  /** データベースへの保存が成功したかどうか */
  saved: boolean;
}
