// Package models はデータベースや外部入出力に対応するドメインオブジェクトを提供します。
package models

import "time"

// Question はデータベースの questions テーブルに対応する構造体です。
// CorrectAnswer の `json:"-"` タグは、APIレスポンスに正解が含まれるのを防ぎます。
type Question struct {
	ID               int       `json:"id"`
	Level            int       `json:"level"`
	ProblemStatement string    `json:"problem_statement"`
	CorrectAnswer    string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}
