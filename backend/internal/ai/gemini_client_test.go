package ai

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestGeminiClient_GenerateSuccess は、1 回のリクエストで正常なレスポンスが返り、
// 結果やメトリクスが期待通りであることを確認するハッピーパスのテスト。
func TestGeminiClient_GenerateSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 疑似サーバーを作り、1 回の呼び出しで 200 とテキストを返す。
		// 本番では外部 API への HTTP POST を行うので、ここでも POST 以外は拒否する。
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"42"}]}}]}`))
	}))
	defer ts.Close()

	var observed Metric
	cfg := Config{
		APIKey:     "test-key",
		BaseURL:    ts.URL,
		Model:      "unit-test",
		Timeout:    time.Second,
		MaxRetries: 0,
		Observer: func(m Metric) {
			// 成功時の Metric が想定通り渡されるか検証するため保持。
			observed = m
		},
	}

	client, err := NewGeminiClient(cfg, ts.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 実際にクライアントを通じてリクエストを投げる。
	resp, err := client.Generate(context.Background(), "hello")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.RawText != "42" {
		t.Fatalf("unexpected raw text: %s", resp.RawText)
	}
	// オブザーバーに渡されたメトリクスが想定通りかを確認する。
	if observed.Status != "success" || observed.Attempts != 1 {
		t.Fatalf("unexpected metric: %+v", observed)
	}
}

// TestGeminiClient_RetryOnServerError は、サーバーエラーが発生した際に
// リトライが実行されることを確認するテスト。
func TestGeminiClient_RetryOnServerError(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 1 回目は 500 (内部エラー) を返し、2 回目に成功レスポンスを返す。
		// これにより、リトライが正しく動作するかを検証できる。
		count := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		if count == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"message":"internal"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"ok"}]}}]}`))
	}))
	defer ts.Close()

	cfg := Config{
		APIKey:     "test-key",
		BaseURL:    ts.URL,
		Model:      "unit-test",
		Timeout:    time.Second,
		MaxRetries: 1,
	}

	client, err := NewGeminiClient(cfg, ts.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// リトライが機能していれば最終的には成功レスポンスが返ってくるはず。
	resp, err := client.Generate(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if resp.RawText != "ok" {
		t.Fatalf("unexpected raw text: %s", resp.RawText)
	}
	// 呼び出し回数が 2 回になっているかでリトライの有無を確認する。
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

// TestGeminiClient_GenerateAnswer は、補助的な GenerateAnswer メソッドが
// 内部で Generate を呼ぶだけで結果を返すことを確認するテスト。
func TestGeminiClient_GenerateAnswer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 疑似サーバーは常に固定のテキストを返す。
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"answer"}]}}]}`))
	}))
	defer ts.Close()

	cfg := Config{
		APIKey:  "test-key",
		BaseURL: ts.URL,
		Model:   "unit-test",
		Timeout: time.Second,
	}

	client, err := NewGeminiClient(cfg, ts.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 期待通りに文字列が返ることを検証する。
	text, err := client.GenerateAnswer(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "answer" {
		t.Fatalf("unexpected text: %s", text)
	}
}

// TestGeminiClient_Unauthorized は、認証エラーが発生した際に
// 適切なエラー種別が返されることを確認するテスト。
func TestGeminiClient_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 認証に失敗したケースを再現するため、常に 401 を返す。
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"unauthorized"}}`))
	}))
	defer ts.Close()

	cfg := Config{
		APIKey:     "bad-key",
		BaseURL:    ts.URL,
		Model:      "unit-test",
		Timeout:    time.Second,
		MaxRetries: 1,
	}

	client, err := NewGeminiClient(cfg, ts.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Generate(context.Background(), "prompt")
	if err == nil {
		t.Fatalf("expected error")
	}
	// エラー種別が Unauthorized に分類されているか確認する。
	if !IsKind(err, ErrorKindUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

// TestGeminiClient_Timeout は、タイムアウト時間内に応答が得られない場合に
// context.DeadlineExceeded が返ることを確認するテスト。
func TestGeminiClient_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 指定タイムアウトを超えるまで応答を遅延させ、タイムアウトが発生する状況を作る。
		time.Sleep(150 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"late"}]}}]}`))
	}))
	defer ts.Close()

	cfg := Config{
		APIKey:     "test-key",
		BaseURL:    ts.URL,
		Model:      "unit-test",
		Timeout:    50 * time.Millisecond,
		MaxRetries: 0,
	}

	httpClient := ts.Client()
	httpClient.Timeout = cfg.Timeout

	client, err := NewGeminiClient(cfg, httpClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = client.Generate(ctx, "prompt")
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	// タイムアウト時は context.DeadlineExceeded が返る想定。
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}

// TestNewGeminiClientFromEnv_QuotedValues は、環境変数に余計な引用符が含まれていても
// 適切にトリムされて読み込まれることを確認するテスト。
func TestNewGeminiClientFromEnv_QuotedValues(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", `"quoted-key"`)
	t.Setenv("AI_MODEL_NAME", `"gemini-test"`)
	t.Setenv("GEMINI_API_BASE", `"https://example.com/v1"`)

	client, err := NewGeminiClientFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 値が引用符なしに読み込まれているかを順に確認する。
	if client.config.APIKey != "quoted-key" {
		t.Fatalf("expected trimmed API key, got %q", client.config.APIKey)
	}
	if client.config.Model != "gemini-test" {
		t.Fatalf("expected trimmed model, got %q", client.config.Model)
	}
	if client.config.BaseURL != "https://example.com/v1" {
		t.Fatalf("expected trimmed base URL, got %q", client.config.BaseURL)
	}
}
