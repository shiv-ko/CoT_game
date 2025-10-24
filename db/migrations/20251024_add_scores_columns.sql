-- Migration: Add extended columns to scores table for evaluation metadata
-- Created: 2025-10-24
-- Purpose: Store AI model details, extracted numeric answer, latency, and evaluation detail JSON

-- Add new columns to scores table
ALTER TABLE scores
ADD COLUMN model_vendor TEXT NOT NULL DEFAULT 'gemini',
ADD COLUMN model_name TEXT NULL,
ADD COLUMN answer_number NUMERIC NULL,
ADD COLUMN latency_ms INT NOT NULL DEFAULT 0,
ADD COLUMN evaluation_detail JSONB NULL;

-- Comments for documentation
COMMENT ON COLUMN scores.model_vendor IS 'AI model vendor (e.g., gemini, openai). Default: gemini';
COMMENT ON COLUMN scores.model_name IS 'Specific model name/version (e.g., gemini-1.5-flash)';
COMMENT ON COLUMN scores.answer_number IS 'Extracted numeric value from AI response';
COMMENT ON COLUMN scores.latency_ms IS 'Response latency in milliseconds. Default: 0';
COMMENT ON COLUMN scores.evaluation_detail IS 'JSON object containing evaluation metadata (diff, mode, precision, etc.)';

-- Rollback instructions:
-- To revert this migration, run the following:
/*
ALTER TABLE scores
DROP COLUMN model_vendor,
DROP COLUMN model_name,
DROP COLUMN answer_number,
DROP COLUMN latency_ms,
DROP COLUMN evaluation_detail;
*/
