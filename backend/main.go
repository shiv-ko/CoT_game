/*
Go製のAPIサーバー」を起動するエントリポイント
役割は大きく4つ：
1. .envから設定を読む
2. Gemini（AI API）クライアントを初期化
3. PostgreSQLへの接続プールを作成
4. Gin（Webフレームワーク）でHTTPルーターを立ててAPIを公開
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/shiv/CoT_game/backend/handlers"
	"github.com/shiv/CoT_game/backend/internal/ai"
	"github.com/shiv/CoT_game/backend/routes"
)

// init は main() の前に実行され、.env ファイルから環境変数を読み込みます。
// これにより、DATABASE_URL のような必要な設定がすべて利用可能になります。
func init() {
	// godotenv パッケージを使用して、.envを簡単に読み込む。
	// err != nil なら .env ファイルがないので、環境変数に依存する。
	if err := godotenv.Load(); err != nil {
		log.Println(".env ファイルが見つかりません。システムの環境変数に依存します。")
	}
}

// createDbPool は新しいデータベース接続プールを作成して返します。
// データベース接続に失敗した場合は panic を起こします。これは、アプリケーションがDBなしでは機能しないため。
func createDbPool(ctx context.Context) (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("エラー: DATABASE_URL 環境変数が設定されていません。")
	}

	// pgxpool を使って接続プールを作成します。
	// pgxpool は高性能でスレッドセーフなPostgreSQL接続プールです。
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("接続プールの作成に失敗しました: %w", err)
	}

	// Pingを使ってデータベースへの接続が有効か確認します。
	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("データベースへの接続に失敗しました: %w", err)
	}

	log.Println("データベースに正常に接続しました。")
	return dbpool, nil
}

// initの実行後にmainが実行される。
func main() {
	// アプリケーションのコンテキストを作成します。
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Printf("起動に失敗しました: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Gemini クライアントを初期化し、APIキー未設定時は起動を停止します。
	if _, err := ai.NewGeminiClientFromEnv(); err != nil {
		return fmt.Errorf("gemini クライアントの初期化に失敗しました: %w", err)
	}

	// データベース接続プールを作成します。
	dbpool, err := createDbPool(ctx)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	// デフォルトのミドルウェアを使用してGinルーターを初期化します。
	router := gin.Default()

	// CORSミドルウェアの設定
	// フロントエンド (http://localhost:3000) からのリクエストを許可します。
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		return fmt.Errorf("サーバーの起動に失敗しました: %w", err)
	}
	return nil
}
