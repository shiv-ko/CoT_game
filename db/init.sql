CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    level INT NOT NULL,
    problem_statement TEXT NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE scores (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    question_id INT REFERENCES questions(id),
    prompt TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    score INT NOT NULL,
    model_vendor TEXT NOT NULL DEFAULT 'gemini',
    model_name TEXT NULL,
    answer_number NUMERIC NULL,
    latency_ms INT NOT NULL DEFAULT 0,
    evaluation_detail JSONB NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert some initial data for testing
INSERT INTO users (username, password_hash) VALUES ('testuser', 'testhash');
INSERT INTO questions (level, problem_statement, correct_answer) VALUES (1, '1+1を計算してください。', '2');
INSERT INTO questions (level, problem_statement, correct_answer) VALUES (3, 'strawberryの中にrは何個ある？', '3');
INSERT INTO questions (level, problem_statement, correct_answer) VALUES (4, '「すもももももももものうち」の中に「も」は何個ある？', '8');
INSERT INTO questions (level, problem_statement, correct_answer) VALUES (3, '太郎君は今年12歳です。お父さんは太郎君の3倍の年齢です。お父さんは何歳ですか？', '36');
INSERT INTO questions (level, problem_statement, correct_answer) VALUES (2, '1, 2, 4, 8, 16, ... 次に来る数は？', '32');


