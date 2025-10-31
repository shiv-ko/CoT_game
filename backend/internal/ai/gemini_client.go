package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// defaultBaseURL は公式 REST エンドポイント。環境変数で差し替え可能。
	defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	// defaultModel は app 全体でデフォルト利用するモデル名。
	defaultModel = "gemini-2.0-flash-lite"
	// defaultTimeout は API 呼び出し 1 件あたりの目安タイムアウト。
	// Gemini APIのレスポンスが15秒程度かかることがあるため、余裕を持って30秒に設定。
	defaultTimeout = 30 * time.Second
)

// GeminiClient は Gemini REST API を利用してテキスト生成を行うクライアントです。
type GeminiClient struct {
	httpClient *http.Client
	config     Config
}

// NewGeminiClient は設定をバリデーションし、使用可能な Gemini クライアントを返します。
// Config の MaxRetries は「追加で許可する再試行回数」を意味し、0 の場合はリトライ無しです。
func NewGeminiClient(cfg Config, client *http.Client) (*GeminiClient, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.Model == "" {
		cfg.Model = defaultModel
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if client == nil {
		// 明示的に http.Client を渡さない場合はここで生成する。
		client = &http.Client{Timeout: cfg.Timeout}
	} else if client.Timeout == 0 {
		// 呼び出し元がタイムアウト未設定のクライアントを渡したときは、デフォルトを補完する。
		client.Timeout = cfg.Timeout
	}

	return &GeminiClient{
		httpClient: client,
		config:     cfg,
	}, nil
}

// NewGeminiClientFromEnv は環境変数から設定を生成します。
//
// 利用例:
//
//	client, err := ai.NewGeminiClientFromEnv()
//	if err != nil {
//		return err
//	}
//	text, err := client.GenerateAnswer(ctx, "1+1=?")
//	// text には Gemini からの回答が格納される
func NewGeminiClientFromEnv() (*GeminiClient, error) {
	baseURL := readEnv("GEMINI_API_BASE")
	model := readEnv("AI_MODEL_NAME")
	cfg := Config{
		APIKey:     readEnv("GEMINI_API_KEY"),
		BaseURL:    defaultBaseURL,
		Model:      defaultModel,
		Timeout:    defaultTimeout,
		MaxRetries: 1,
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if model != "" {
		cfg.Model = model
	}
	return NewGeminiClient(cfg, nil)
}

// Generate は Gemini API にプロンプトを送信し、レスポンスを返します。
func (c *GeminiClient) Generate(ctx context.Context, prompt string) (Response, error) {
	if strings.TrimSpace(prompt) == "" {
		// 空文字列を送るとベンダー側が 400 を返すため、早めに防御する。
		return Response{}, errors.New("ai: prompt is empty")
	}
	payload, err := json.Marshal(newGenerateContentRequest(prompt))
	if err != nil {
		// JSON へのシリアライズは理論上失敗しないが、念のためエラーを伝播する。
		return Response{}, fmt.Errorf("ai: failed to marshal request: %w", err)
	}

	ctx, cancel := c.ensureTimeout(ctx)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// 現在の試行を開始した時刻。レイテンシとメトリクス計算に使用する。
		attemptStart := time.Now()
		resp, err := c.invoke(ctx, payload)
		if err == nil {
			resp.Latency = time.Since(attemptStart)
			c.emitMetric(successMetric(c.config.Model, attempt+1, resp.Latency, len(prompt)))
			return resp, nil
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			c.emitMetric(failureMetric(c.config.Model, attempt+1, time.Since(attemptStart), len(prompt), ctxErr))
			return Response{}, ctxErr
		}

		var apiErr *Error
		if errors.As(err, &apiErr) {
			if !apiErr.Temp || attempt == c.config.MaxRetries {
				c.emitMetric(failureMetric(c.config.Model, attempt+1, time.Since(attemptStart), len(prompt), err))
				return Response{}, err
			}
		} else {
			// ネットワークエラーなど API 固有でない失敗もここに到達する。
			// 追加の試行余地が無い場合は直ちに返す。
			if attempt == c.config.MaxRetries {
				c.emitMetric(failureMetric(c.config.Model, attempt+1, time.Since(attemptStart), len(prompt), err))
				return Response{}, err
			}
		}

		lastErr = err
		// 単純な線形バックオフ。1 回目 250ms, 2 回目 500ms ...。
		backoff := time.Duration(attempt+1) * 250 * time.Millisecond
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			c.emitMetric(failureMetric(c.config.Model, attempt+1, time.Since(attemptStart), len(prompt), ctx.Err()))
			return Response{}, ctx.Err()
		}
	}

	if lastErr != nil {
		c.emitMetric(failureMetric(c.config.Model, c.config.MaxRetries+1, 0, len(prompt), lastErr))
		return Response{}, lastErr
	}
	c.emitMetric(failureMetric(c.config.Model, c.config.MaxRetries+1, 0, len(prompt), errors.New("ai: request failed without specific error")))
	return Response{}, errors.New("ai: request failed without specific error")
}

// GenerateAnswer はテキストのみが必要な場合の簡易アクセサです。
func (c *GeminiClient) GenerateAnswer(ctx context.Context, prompt string) (string, error) {
	resp, err := c.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}
	return resp.RawText, nil
}

func (c *GeminiClient) invoke(ctx context.Context, body []byte) (Response, error) {
	endpoint := c.endpoint()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("ai: failed to create request: %w", err)
	}

	// REST API は API キーをヘッダ・クエリのいずれでも受け付ける。
	// 可読性向上のため両方に設定している。
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.config.APIKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	data, err := io.ReadAll(io.LimitReader(res.Body, 5<<20))
	if err != nil {
		return Response{}, fmt.Errorf("ai: failed to read response: %w", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return parseGenerateContentResponse(data)
	}

	// 非 2xx の場合は GCP のエラー形式を優先的に解釈し、分類したエラーを返す。
	apiErr := parseAPIError(res.StatusCode, data)
	if apiErr == nil {
		apiErr = &Error{
			Kind:    ErrorKindServerError,
			Code:    res.StatusCode,
			Message: fmt.Sprintf("ai: unexpected status %d", res.StatusCode),
			Temp:    res.StatusCode >= 500,
		}
	}

	return Response{}, apiErr
}

func (c *GeminiClient) endpoint() string {
	base := strings.TrimSuffix(c.config.BaseURL, "/")
	model := url.PathEscape(c.config.Model)
	return fmt.Sprintf("%s/models/%s:generateContent?key=%s", base, model, url.QueryEscape(c.config.APIKey))
}

func newGenerateContentRequest(prompt string) generateContentRequest {
	// Gemini の generateContent は Contents の配列を要求する。
	return generateContentRequest{
		Contents: []content{
			{
				Role: "user",
				Parts: []part{{
					Text: prompt,
				}},
			},
		},
	}
}

func parseGenerateContentResponse(data []byte) (Response, error) {
	var decoded generateContentResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		return Response{}, fmt.Errorf("ai: failed to decode response: %w", err)
	}

	if len(decoded.Candidates) == 0 || len(decoded.Candidates[0].Content.Parts) == 0 {
		return Response{}, errors.New("ai: empty response from Gemini")
	}

	// 最初の Candidate の最初の Part を利用する。Gemini 側の仕様ではここに主要回答が入る。
	text := decoded.Candidates[0].Content.Parts[0].Text
	return Response{RawText: text}, nil
}

func parseAPIError(status int, data []byte) *Error {
	var apiErr apiErrorResponse
	if err := json.Unmarshal(data, &apiErr); err != nil {
		return nil
	}

	if apiErr.Error.Message == "" {
		return nil
	}

	e := &Error{
		Code:    status,
		Message: apiErr.Error.Message,
	}

	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		e.Kind = ErrorKindUnauthorized
		e.Temp = false
	case http.StatusTooManyRequests:
		e.Kind = ErrorKindServerError
		e.Temp = true
	default:
		if status >= 500 {
			e.Kind = ErrorKindServerError
			e.Temp = true
		} else {
			e.Kind = ErrorKindClientError
			e.Temp = false
		}
	}

	return e
}

func (c *GeminiClient) ensureTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.config.Timeout)
}

func (c *GeminiClient) emitMetric(metric Metric) {
	if c.config.Observer != nil {
		c.config.Observer(metric)
	}
}

func successMetric(model string, attempts int, latency time.Duration, promptLen int) Metric {
	return Metric{
		Model:      model,
		Attempts:   attempts,
		Latency:    latency,
		PromptSize: promptLen,
		Status:     "success",
	}
}

func failureMetric(model string, attempts int, latency time.Duration, promptLen int, err error) Metric {
	return Metric{
		Model:      model,
		Attempts:   attempts,
		Latency:    latency,
		PromptSize: promptLen,
		Status:     "failure",
		Err:        err,
	}
}

func readEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	trimmed := strings.TrimSpace(val)
	if len(trimmed) >= 2 {
		if (strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"")) ||
			(strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "'")) {
			trimmed = trimmed[1 : len(trimmed)-1]
		}
	}
	// 値に含まれていた外側のクォートを除去したうえで再度トリムする。
	return strings.TrimSpace(trimmed)
}

type generateContentRequest struct {
	Contents []content `json:"contents"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text,omitempty"`
}

type generateContentResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content content `json:"content"`
}

type apiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}
