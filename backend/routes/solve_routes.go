// Package routes は HTTP ルートのグルーピングとマッピングを管理します。
// solve_routes.go は solve ハンドラを REST ルーティングへ結び付け、API モジュールが一貫した URL 設計を保てるよう整理します。
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shiv/CoT_game/backend/handlers"
)

// RegisterSolveRoutes は Solve エンドポイントを登録します。
func RegisterSolveRoutes(api *gin.RouterGroup, h *handlers.SolveHandler) {
	// POST /api/v1/solve
	// RouterGroup は /api/v1 のような共通 prefix をまとめるための仕組みです。
	// ここでは「/solve に POST されたら SolveHandler.PostSolve を呼ぶ」という関連付けを 1 行で表現しています。
	api.POST("/solve", h.PostSolve)
}
