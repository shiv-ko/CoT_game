-- Migration: Add extended columns to scores table for evaluation metadata
-- Created: 2025-10-24
-- Purpose: Store AI model details, extracted numeric answer, latency, and evaluation detail JSON

-- 新しい列を追加
-- ALTERはテーブル構造を変更するための書き出しコマンド
ALTER TABLE scores\
-- AIベンダーの名前・デフォルトは 'gemini'
ADD COLUMN model_vendor TEXT NOT NULL DEFAULT 'gemini',
-- AIモデル名/バージョン
ADD COLUMN model_name TEXT NULL,
-- AI応答から抽出された数値
ADD COLUMN answer_number NUMERIC NULL,
-- レイテンシ (ミリ秒)、デフォルトは 0
ADD COLUMN latency_ms INT NOT NULL DEFAULT 0,
-- メタデータを含む詳細情報をjson形式で格納できる列
-- 細かい情報を詰め込めるようにNULL許容
ADD COLUMN evaluation_detail JSONB NULL;

-- 各列に説明文をつけている。
-- これにより開発者が列の目的を理解しやすくなる。
COMMENT ON COLUMN scores.model_vendor IS 'AI model vendor (e.g., gemini, openai). Default: gemini';
COMMENT ON COLUMN scores.model_name IS 'Specific model name/version (e.g., gemini-1.5-flash)';
COMMENT ON COLUMN scores.answer_number IS 'Extracted numeric value from AI response';
COMMENT ON COLUMN scores.latency_ms IS 'Response latency in milliseconds. Default: 0';
COMMENT ON COLUMN scores.evaluation_detail IS 'JSON object containing evaluation metadata (diff, mode, precision, etc.)';


-- ロールバック用のコマンド:
-- もとに戻す場合は以下のコマンドを実行する:
/*
ALTER TABLE scores
DROP COLUMN model_vendor,
DROP COLUMN model_name,
DROP COLUMN answer_number,
DROP COLUMN latency_ms,
DROP COLUMN evaluation_detail;
*/
