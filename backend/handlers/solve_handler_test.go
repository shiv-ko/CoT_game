// solve_handler_test.go は SolveHandler の挙動をモックやテスト DB で再現し、HTTP レイヤーの入力・出力・エラー処理が仕様通りかを保証します。
// これもDBまで含むので結合テストになる。
package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/shiv/CoT_game/backend/internal/ai"
	"github.com/shiv/CoT_game/backend/internal/repository"
)

// MockAIClient はテスト用のAIクライアントモックです。
type MockAIClient struct {
	Response ai.Response
	Err      error
}

func (m *MockAIClient) Generate(_ context.Context, _ string) (ai.Response, error) {
	if m.Err != nil {
		return ai.Response{}, m.Err
	}
	return m.Response, nil
}

func (m *MockAIClient) GenerateAnswer(_ context.Context, _ string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Response.RawText, nil
}

// setupTestDB はテスト用のDB接続を作成します。
// テスト用コンテナが立っていないときに無理に失敗させず、Skip でテストスイート全体を止めない方針です。
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=cot_game sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("テストDB接続失敗（スキップ）: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		t.Skipf("テストDB疎通失敗（スキップ）: %v", err)
		return nil
	}

	return db
}

// cleanupTestScores はテストデータをクリーンアップします。
func cleanupTestScores(t *testing.T, db *sql.DB, questionID int) {
	t.Helper()
	_, _ = db.Exec("DELETE FROM scores WHERE question_id = $1", questionID)
}

func TestSolveHandler_PostSolve_Success(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	// テスト用の問題が存在することを確認（ID=1の問題を想定）
	// 本番 DB と同じスキーマを使っているので、フィクスチャを用意できていない場合はスキップします。
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE id = 1)").Scan(&exists)
	if err != nil || !exists {
		t.Skip("テスト用の問題（ID=1）が存在しないためスキップ")
	}

	defer cleanupTestScores(t, db, 1)

	// モックAIクライアント：正解「2」を返す
	// latency は Evaluate の外で使うので、敢えて値を埋めておきます。
	mockAI := &MockAIClient{
		Response: ai.Response{
			RawText: "答えは2です。",
			Latency: 100 * time.Millisecond,
		},
	}

	scoreRepo := repository.NewScoresRepository(db)
	handler := NewSolveHandler(mockAI, scoreRepo, db)

	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/solve", handler.PostSolve)

	// リクエストボディ
	reqBody := SolveRequest{
		QuestionID: 1,
		Prompt:     "1+1を計算してください",
		Model:      "gemini-1.5-flash",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// HTTPリクエスト
	req := httptest.NewRequest(http.MethodPost, "/api/v1/solve", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// httptest.NewRecorder によって疑似 HTTP サーバーを用意し、Router.ServeHTTP に通すとミドルウェアも含めた実際の挙動を確認できます。
	router.ServeHTTP(w, req)

	// ステータスコード確認
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
		t.Logf("response body: %s", w.Body.String())
		return
	}

	// レスポンスボディのパース
	// httptest.Recorder の Body は bytes.Buffer なので、そのまま JSON デコードできます。
	var resp SolveResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// 検証
	if resp.QuestionID != 1 {
		t.Errorf("expected question_id=1, got %d", resp.QuestionID)
	}
	if resp.Score <= 0 {
		t.Errorf("expected positive score, got %d", resp.Score)
	}
	if resp.AIOutput != "答えは2です。" {
		t.Errorf("expected ai_output='答えは2です。', got %s", resp.AIOutput)
	}
	if !resp.Saved {
		t.Error("expected saved=true")
	}

	// DBに保存されたか確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM scores WHERE question_id = 1").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count scores: %v", err)
	}
	if count == 0 {
		t.Error("expected score to be saved in DB")
	}
}

func TestSolveHandler_PostSolve_ValidationError(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	mockAI := &MockAIClient{
		Response: ai.Response{RawText: "test"},
	}
	scoreRepo := repository.NewScoresRepository(db)
	handler := NewSolveHandler(mockAI, scoreRepo, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/solve", handler.PostSolve)

	tests := []struct {
		name       string
		reqBody    interface{}
		wantStatus int
	}{
		{
			name: "プロンプトが空",
			reqBody: SolveRequest{
				QuestionID: 1,
				Prompt:     "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "プロンプトが長すぎる",
			reqBody: SolveRequest{
				QuestionID: 1,
				Prompt:     string(make([]byte, 2001)),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "question_idが不正",
			reqBody: map[string]interface{}{
				"prompt": "test",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// map[string]interface{} で書いたパターンも JSON にできるので、柔軟にリクエストを組み立てられます。
			bodyBytes, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/solve", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestSolveHandler_PostSolve_QuestionNotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	mockAI := &MockAIClient{
		Response: ai.Response{RawText: "test"},
	}
	scoreRepo := repository.NewScoresRepository(db)
	handler := NewSolveHandler(mockAI, scoreRepo, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/solve", handler.PostSolve)

	reqBody := SolveRequest{
		QuestionID: 999999, // 存在しない問題ID
		Prompt:     "test prompt",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/solve", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 存在しない問題 ID を指定した場合、ハンドラが 404 を返すかを検証します。
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestSolveHandler_PostSolve_AIError(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	// 問題1が存在することを確認
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE id = 1)").Scan(&exists)
	if err != nil || !exists {
		t.Skip("テスト用の問題（ID=1）が存在しないためスキップ")
	}

	// AIエラーをシミュレート
	// Temp=true のエラーは一時的な障害を意味し、ハンドラが 502 を返すか確認します。
	mockAI := &MockAIClient{
		Err: &ai.Error{
			Kind:    ai.ErrorKindServerError,
			Code:    500,
			Message: "AI server error",
			Temp:    true,
		},
	}
	scoreRepo := repository.NewScoresRepository(db)
	handler := NewSolveHandler(mockAI, scoreRepo, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/solve", handler.PostSolve)

	reqBody := SolveRequest{
		QuestionID: 1,
		Prompt:     "test prompt",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/solve", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected status 502, got %d", w.Code)
	}
}
