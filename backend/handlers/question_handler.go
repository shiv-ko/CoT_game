package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shiv/CoT_game/backend/models"
)

// QuestionHandler はデータベース接続プールを保持します。
// このアプローチはデータベース接続の依存性注入（DI）を利用しており、
// ハンドラのテスト容易性や再利用性を高めます。
type QuestionHandler struct {
	DB *pgxpool.Pool
}

// NewQuestionHandler は question に関連するルートの新しいハンドラを作成します。
// オブジェクトの初期化には、一般的なファクトリーパターンを採用しています。
func NewQuestionHandler(db *pgxpool.Pool) *QuestionHandler {
	return &QuestionHandler{DB: db}
}

// GetQuestions はデータベースから全ての問題を取得します。
// データベースのクエリやデータスキャン中に発生しうるエラーをハンドリングし、
// 適切なHTTPステータスコードとエラーメッセージを返します。
func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	// 全ての問題を選択するクエリ。一貫性を保つためにlevelで並び替えます。
	query := "SELECT id, level, problem_statement, correct_answer, created_at FROM questions ORDER BY level"

	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		// サーバー側のデバッグ目的で詳細なエラーをログに出力します。
		log.Printf("質問のクエリ実行中にエラーが発生しました: %v", err)
		// クライアントには汎用的なエラーメッセージを返します。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースからの質問の取得に失敗しました。"})
		return
	}
	defer rows.Close()

	// 質問データを保持するためのスライスを作成します。
	var questions []models.Question
	for rows.Next() {
		var q models.Question
		// 行データをQuestion構造体にスキャンします。
		if err := rows.Scan(&q.ID, &q.Level, &q.ProblemStatement, &q.CorrectAnswer, &q.CreatedAt); err != nil {
			log.Printf("質問行のスキャン中にエラーが発生しました: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "質問データの処理に失敗しました。"})
			return
		}
		questions = append(questions, q)
	}

	// 行のイテレーション中に発生したエラーを確認します。
	if err := rows.Err(); err != nil {
		log.Printf("質問行のイテレーション中にエラーが発生しました: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "質問の取得中にエラーが発生しました。"})
		return
	}

	// データベースに質問が見つからない場合、フロントエンドの互換性のために
	// `null` の代わりに空のリスト `[]` を返します。
	if questions == nil {
		questions = make([]models.Question, 0)
	}

	// 200 OKステータスと共に質問のリストを返します。
	c.JSON(http.StatusOK, questions)
}
