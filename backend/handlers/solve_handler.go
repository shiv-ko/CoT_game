// Package handlers は HTTP リクエストを処理しレスポンスを生成するコントローラ層を提供します。
// solve_handler.go は /api/v1/solve のリクエストを受け取り、AI 呼び出し→評価→DB 保存→レスポンス生成という一連の業務フローを司る中心的なハンドラです。
package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiv/CoT_game/backend/internal/ai"
	"github.com/shiv/CoT_game/backend/internal/eval"
	"github.com/shiv/CoT_game/backend/internal/repository"
)

// SolveHandler は solve エンドポイントの依存関係を保持します。
type SolveHandler struct {
	AIClient  ai.Client                   // AI へ質問を投げるための依存。モック可能にするため interface で受け取ります。
	ScoreRepo repository.ScoresRepository // スコア保存・取得を担うリポジトリ。DB 直書きよりテストしやすい構造です。
	DB        *sql.DB                     // 正解を問い合わせるための生 SQL 接続。将来的に専用リポジトリを切り出す余地があります。
}

// NewSolveHandler は新しい SolveHandler を作成します。
func NewSolveHandler(aiClient ai.Client, scoreRepo repository.ScoresRepository, db *sql.DB) *SolveHandler {
	return &SolveHandler{
		AIClient:  aiClient,
		ScoreRepo: scoreRepo,
		DB:        db,
	}
}

// SolveRequest は /api/v1/solve の入力リクエストを表します。
type SolveRequest struct {
	QuestionID int    `json:"question_id" binding:"required"`
	Prompt     string `json:"prompt" binding:"required"`
	Model      string `json:"model"`
}

// SolveResponse は /api/v1/solve のレスポンスを表します。
type SolveResponse struct {
	QuestionID   int                    `json:"question_id"`
	Prompt       string                 `json:"prompt"`
	ModelVendor  string                 `json:"model_vendor"`
	ModelName    string                 `json:"model_name"`
	AIOutput     string                 `json:"ai_output"`     // AI が出力したテキスト全文。クライアントで表示します。
	AnswerNumber *float64               `json:"answer_number"` // 数値回答が抽出できた場合のみ値が入ります（例: 算数の答え）。
	Score        int                    `json:"score"`         // 評価ロジックで決まった点数。100 点満点を想定。
	Evaluation   map[string]interface{} `json:"evaluation"`    // 評価モードなどの補足情報。UI の詳細表示に役立ちます。
	ElapsedMs    int64                  `json:"elapsed_ms"`    // AI 応答までにかかった時間（ミリ秒）。
	Saved        bool                   `json:"saved"`         // DB 保存が成功したかどうか。false でもスコア自体は返します。
}

// PostSolve は POST /api/v1/solve のハンドラです。
// リクエストのバリデーションをして、問題なかったら保存のリポジトリに投げる。
// ユーザーのプロンプトをAIに送信し、評価してDBに保存します。
// PostSolve godoc
// @Summary      Solve a question
// @Description  Submit a prompt to solve a question using AI
// @Tags         solve
// @Accept       json
// @Produce      json
// @Param        request body handlers.SolveRequest true "Solve Request"
// @Success      200  {object}  handlers.SolveResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /solve [post]
func (h *SolveHandler) PostSolve(c *gin.Context) {
	var req SolveRequest

	// リクエストボディのバリデーション
	// gin の ShouldBindJSON は、構造体タグで required を指定しておくと必須チェックまでまとめて行ってくれます。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "リクエスト形式が不正です",
			"detail":  err.Error(),
		})
		return
	}

	// プロンプトの長さチェック
	if len(req.Prompt) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_prompt",
			"message": "プロンプトが空です",
		})
		return
	}
	if len(req.Prompt) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "prompt_too_long",
			"message": "プロンプトが長すぎます（最大2000文字）",
		})
		return
	}

	// デフォルトモデルを設定
	if req.Model == "" {
		// モデル指定が無いケースでも使いやすいよう、サービス側でデフォルト値を決めています。
		req.Model = "gemini-2.0-flash-lite"
	}

	ctx := c.Request.Context()

	// 問題の存在確認と正解の取得
	// 問題が存在しない場合は 404 を返してフロントに伝えます。
	correctAnswer, err := h.getCorrectAnswer(ctx, req.QuestionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "question_not_found",
				"message": "指定された問題が見つかりません",
			})
		} else {
			log.Printf("問題取得エラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "database_error",
				"message": "問題の取得に失敗しました",
			})
		}
		return
	}

	// 問題文の取得
	// AIに問題文を含めたプロンプトを送信するために必要です。
	problemStatement, err := h.getProblemStatement(ctx, req.QuestionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "question_not_found",
				"message": "指定された問題が見つかりません",
			})
		} else {
			log.Printf("問題文取得エラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "database_error",
				"message": "問題文の取得に失敗しました",
			})
		}
		return
	}

	// システムプロンプトとユーザープロンプトを結合
	// 問題文はユーザーには見せませんが、AIには送信する必要があります。
	combinedPrompt := buildCombinedPrompt(problemStatement, req.Prompt)

	// デバッグ: 送信されるプロンプトを確認
	log.Printf("=== プロンプト確認 ===")
	log.Printf("問題ID: %d", req.QuestionID)
	log.Printf("ユーザープロンプト: %s", req.Prompt)
	log.Printf("問題文: %s", problemStatement)
	log.Printf("結合後のプロンプト:\n%s", combinedPrompt)
	log.Printf("===================")

	// AI呼び出し開始時刻
	// time.Since と組み合わせることでレイテンシを簡単に測定できます。
	startTime := time.Now()

	// AIクライアントで結合されたプロンプトを送信
	// AIClient は interface のため、実運用では外部 API を呼び、テストではモックを差し込めます。
	aiResp, err := h.AIClient.Generate(ctx, combinedPrompt)
	if err != nil {
		log.Printf("AI呼び出しエラー: %v", err)
		statusCode := http.StatusBadGateway
		if ai.IsKind(err, ai.ErrorKindClientError) {
			statusCode = http.StatusBadRequest
		} else if ai.IsKind(err, ai.ErrorKindUnauthorized) {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, gin.H{
			"error":   "ai_error",
			"message": "AI応答の取得に失敗しました",
			"detail":  err.Error(),
		})
		return
	}

	// レイテンシ計算
	elapsedMs := time.Since(startTime).Milliseconds()

	// AIの完全な回答（DBに保存用）
	fullAIResponse := aiResp.RawText

	// クライアントに返す回答（「最終回答: 」以降のみ）
	// 問題文が推測されないよう、説明部分を除外します
	clientResponse := extractFinalAnswer(aiResp.RawText)

	// 評価ロジック実行
	// eval パッケージに責務を分離することで、ハンドラは「AI の結果をどう扱うか」に集中できます。
	score, answerNumber, mode, detail := eval.Evaluate(fullAIResponse, correctAnswer)

	// 評価メタデータを構築
	// detail 全体は JSONB に保存しますが、レスポンスに最低限の情報を添えておくと UI 側で扱いやすくなります。
	evaluationMeta := map[string]interface{}{
		"mode":   mode,
		"detail": detail,
	}

	// スコアレコードをDBに保存
	// 完全な回答をDBに保存（デバッグ・分析用）
	scoreRecord := &repository.Score{
		UserID:           nil, // ゲストユーザー（認証未実装のため）
		QuestionID:       req.QuestionID,
		Prompt:           req.Prompt,
		AIResponse:       fullAIResponse,
		Score:            score,
		ModelVendor:      "gemini",
		ModelName:        &req.Model,
		AnswerNumber:     answerNumber,
		LatencyMs:        int(elapsedMs),
		EvaluationDetail: detail,
	}

	// 保存は可能な限り試みますが、失敗しても回答自体はクライアントに返せるようにします。
	// ここで、リポジトリ層を使って保存処理を行います。
	saved := true
	if err := h.ScoreRepo.Create(ctx, scoreRecord); err != nil {
		log.Printf("スコア保存エラー: %v", err)
		saved = false
		// 保存失敗しても結果は返す（クライアントには成功を伝える）
	}

	// レスポンス生成
	// フロントエンドには「最終回答: 」以降のみを返して、問題文の推測を防ぎます
	resp := SolveResponse{
		QuestionID:   req.QuestionID,
		Prompt:       req.Prompt,
		ModelVendor:  "gemini",
		ModelName:    req.Model,
		AIOutput:     clientResponse, // 最終回答のみ
		AnswerNumber: answerNumber,
		Score:        score,
		Evaluation:   evaluationMeta,
		ElapsedMs:    elapsedMs,
		Saved:        saved,
	}

	c.JSON(http.StatusOK, resp)
}

// getProblemStatement は question_id から問題文を取得します。
// システムプロンプトの構築に使用されます。
func (h *SolveHandler) getProblemStatement(ctx context.Context, questionID int) (string, error) {
	query := "SELECT problem_statement FROM questions WHERE id = $1"
	var problemStatement string
	// PostgreSQLのバインド変数を使用してSQLインジェクションを防ぎます。
	err := h.DB.QueryRowContext(ctx, query, questionID).Scan(&problemStatement)
	return problemStatement, err
}

// getCorrectAnswer はquestion_idから正解を取得します。
func (h *SolveHandler) getCorrectAnswer(ctx context.Context, questionID int) (string, error) {
	query := "SELECT correct_answer FROM questions WHERE id = $1"
	var correctAnswer string
	// 実環境では questionID をバインドして SQL インジェクションを防ぎます。QueryRowContext → Scan の流れは DB 操作の基本形です。
	err := h.DB.QueryRowContext(ctx, query, questionID).Scan(&correctAnswer)
	return correctAnswer, err
}

// extractFinalAnswer はAIの完全な回答から「最終回答: 」以降の部分のみを抽出します。
// 問題文が推測されないよう、説明部分は除外してクライアントに返します。
func extractFinalAnswer(fullResponse string) string {
	// 「最終回答: 」または「最終回答:」を探す
	markers := []string{"最終回答: ", "最終回答:", "最終回答："}

	for _, marker := range markers {
		if idx := strings.Index(fullResponse, marker); idx != -1 {
			// マーカー以降の部分を取得
			afterMarker := fullResponse[idx+len(marker):]
			// 前後の空白を削除して返す
			return strings.TrimSpace(afterMarker)
		}
	}

	// マーカーが見つからない場合は全文を返す（フォールバック）
	return fullResponse
}

// buildCombinedPrompt はシステムプロンプトとユーザープロンプトを結合します。
// 問題文はユーザーには見せませんが、AIが問題を解くために必要です。
// 返されるプロンプトには、問題文、ユーザーの指示が含まれます。
// ユーザーのプロンプトを最大限尊重しつつ、最終回答を明確にするよう促します。
func buildCombinedPrompt(problemStatement, userPrompt string) string {
	return fmt.Sprintf(`以下の問題について、ユーザーの指示に従って回答してください。

【問題】
%s

【ユーザーの指示】
%s

【重要】回答の最後に、必ず「最終回答: 」に続けて答えを明記してください。`, problemStatement, userPrompt)
}
