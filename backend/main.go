package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/shiv/CoT_game/backend/handlers"
	"github.com/shiv/CoT_game/backend/routes"
)

// init は main() の前に実行され、.env ファイルから環境変数を読み込みます。
// これにより、DATABASE_URL のような必要な設定がすべて利用可能になります。
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env ファイルが見つかりません。システムの環境変数に依存します。")
	}
}

// createDbPool は新しいデータベース接続プールを作成して返します。
// データベース接続に失敗した場合は panic を起こします。これは、アプリケーションがDBなしでは
// 機能しないためです。
func createDbPool(ctx context.Context) *pgxpool.Pool {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("エラー: DATABASE_URL 環境変数が設定されていません。")
	}

	dbpool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("接続プールの作成に失敗しました: %v\n", err)
	}

	// データベースへの接続が有効か確認します。
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v\n", err)
	}

	log.Println("データベースに正常に接続しました。")
	return dbpool
}

func main() {
	// アプリケーションのコンテキストを作成します。
	ctx := context.Background()

	// データベース接続プールを作成します。
	dbpool := createDbPool(ctx)
	defer dbpool.Close()

	// デフォルトのミドルウェアを使用してGinルーターを初期化します。
	router := gin.Default()

	// バージョニングのためにメインのAPIグループを作成します (例: /api/v1)。
	// これにより、APIの保守性が向上し、将来のバージョンへの対応が可能になります。
	apiV1 := router.Group("/api/v1")

	// データベースプールを使用してハンドラを初期化します。
	questionHandler := handlers.NewQuestionHandler(dbpool)

	// questions API のルートを登録します。
	routes.RegisterQuestionRoutes(apiV1, questionHandler)

	// シンプルなヘルスチェック用のエンドポイントです。
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// サーバーを起動し、ポート8080でリッスンします。
	log.Println("サーバーを http://localhost:8080 で起動します")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
