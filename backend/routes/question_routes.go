// Package routes は HTTP ルートのグルーピングとマッピングを管理します。
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shiv/CoT_game/backend/handlers"
)

// RegisterQuestionRoutes は Question リソースのルート（エンドポイント）を設定します。
// 全ての問題関連エンドポイントを共通のパスプレフィックス（例: /api/v1）でグループ化し、
// 一貫性のあるバージョニング可能なAPI構造を保証します。
func RegisterQuestionRoutes(api *gin.RouterGroup, h *handlers.QuestionHandler) {
	// /questions の新しいグループを作成します。
	// 例: /api/v1/questions
	questionRoutes := api.Group("/questions")
	{
		// /api/v1/questions（末尾に / でもOK）に GET を登録。
		// 最終的なパスは /api/v1/questions となります。
		// このエンドポイントは、全ての問題のリストを取得します。
		questionRoutes.GET("", h.GetQuestions)
	}
}
