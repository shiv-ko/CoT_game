// Package ai は AI モデルと対話するためのクライアントインターフェースを定義します。
package ai

import (
	"context"
	"errors"
	"time"
)

// Client は外部AIサービスとの対話を統一するためのインターフェースです。
// 上位レイヤーは具体的な実装（Gemini など）を意識せず、抽象化された契約だけに依存できます。
type Client interface {
	// Generate はプロンプトを送信し、AIが返すテキストとレイテンシ情報をまとめて取得します。
	Generate(ctx context.Context, prompt string) (Response, error)
	// GenerateAnswer は文字列だけが必要なときの薄いラッパーです。
	GenerateAnswer(ctx context.Context, prompt string) (string, error)
}

// Response はAIクライアント呼び出しの結果を保持します。
// 将来的にメタデータ（token usage など）を拡張しやすいよう構造体で管理します。
type Response struct {
	RawText string        // 生成された文章そのもの
	Latency time.Duration // 処理に要した時間（リトライ込みの最終試行）
}

// Config はAIクライアントの共通設定値を扱います。
type Config struct {
	APIKey     string        // 認証トークン。必須。
	BaseURL    string        // エンドポイントのベース URL。テスト時に差し替え可能。
	Model      string        // 利用するモデル名。
	Timeout    time.Duration // context に設定するデフォルトタイムアウト。
	MaxRetries int           // 再試行回数（追加試行数）。
	Observer   func(Metric)  // 呼び出し結果を収集するオプションフック。
}

// Metric は呼び出しごとの軽量な観測情報です。
type Metric struct {
	Model      string        // 呼び出し対象モデル。
	Attempts   int           // 実際に試行した回数（1+リトライ回数）。
	Latency    time.Duration // 最終試行のレイテンシ。
	PromptSize int           // プロンプトの長さ（文字数ベースの簡易値）。
	Status     string        // "success" もしくは "failure"。
	Err        error         // 失敗時のエラー。成功時は nil。
}

// Validate は必須項目をチェックし、不備があればエラーを返します。
func (c Config) Validate() error {
	switch {
	case c.APIKey == "":
		return errors.New("ai: GEMINI_API_KEY is not set")
	case c.BaseURL == "":
		return errors.New("ai: base URL is empty")
	case c.Model == "":
		return errors.New("ai: model is empty")
	case c.Timeout <= 0:
		return errors.New("ai: timeout must be positive")
	}
	return nil
}

// ErrorKind はAPIエラーの分類です。
type ErrorKind string

const (
	// ErrorKindUnauthorized は認証・認可エラーを表します。
	ErrorKindUnauthorized ErrorKind = "unauthorized"
	// ErrorKindClientError は 4xx のうち認証以外のエラーを表します。
	ErrorKindClientError ErrorKind = "client_error"
	// ErrorKindServerError は 5xx など再試行可能なエラーを表します。
	ErrorKindServerError ErrorKind = "server_error"
)

// Error はAIクライアントから返されるドメインエラーです。
type Error struct {
	Kind    ErrorKind
	Code    int
	Message string
	Temp    bool
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Temporary は一時的エラーかどうかを返します（net.Error 互換）。
func (e *Error) Temporary() bool {
	return e != nil && e.Temp
}

// Is は errors.Is 互換の比較を提供します。
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok || e == nil {
		return false
	}
	return e.Kind == t.Kind && (t.Code == 0 || t.Code == e.Code)
}

// IsKind は特定のKindか判定するユーティリティです。
func IsKind(err error, kind ErrorKind) bool {
	var ae *Error
	if errors.As(err, &ae) {
		return ae.Kind == kind
	}
	return false
}

// Sentinel errors for comparisons
var (
	ErrUnauthorized = &Error{Kind: ErrorKindUnauthorized}
	ErrClientError  = &Error{Kind: ErrorKindClientError}
	ErrServerError  = &Error{Kind: ErrorKindServerError, Temp: true}
)
