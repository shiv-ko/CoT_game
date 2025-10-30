// Package models はデータベースや外部入出力に対応するドメインオブジェクトを提供します。
package models

import "time"

// Question は内部処理用の完全な問題情報を表す構造体です。
// データベースの questions テーブルに対応します。
type Question struct {
	ID               int       `json:"id"`
	Level            int       `json:"level"`
	ProblemStatement string    `json:"problem_statement"`
	CorrectAnswer    string    `json:"correct_answer"`
	Tags             []string  `json:"tags"`
	CreatedAt        time.Time `json:"created_at"`
}

// QuestionResponse はクライアントに返す問題情報を表す構造体です。
// セキュリティとゲーム性の観点から、問題文と正解は含まれません。
// ユーザーは問題文を見ずにプロンプトを工夫してAIに正解を導き出させます。
type QuestionResponse struct {
	ID        int       `json:"id"`
	Level     int       `json:"level"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
