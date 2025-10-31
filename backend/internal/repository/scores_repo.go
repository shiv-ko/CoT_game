// Package repository はデータベースアクセスとドメインロジックの間を仲介するリポジトリ層を提供します。
// scores_repo.go は scores テーブルに絞った CRUD と集計処理をまとめ、サービス層から一貫した API として利用できるようにする役割を担います。
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Score は scores テーブルに対応する構造体です。
// Web アプリから受け取った回答結果を丸ごと保持し、
// 後でランキング集計や履歴表示に利用するための入れ物と考えると理解しやすいです。
type Score struct {
	ID               int                    `json:"id"`
	UserID           *int                   `json:"user_id"` // 認証済みユーザーのID。ゲスト投稿も扱いたいので null を許容しています。
	QuestionID       int                    `json:"question_id"`
	Prompt           string                 `json:"prompt"`
	AIResponse       string                 `json:"ai_response"`
	Score            int                    `json:"score"`
	ModelVendor      string                 `json:"model_vendor"`
	ModelName        *string                `json:"model_name"`
	AnswerNumber     *float64               `json:"answer_number"`
	LatencyMs        int                    `json:"latency_ms"`        // AI 応答までの時間（ミリ秒）。体験の快適さを可視化するために保存します。
	EvaluationDetail map[string]interface{} `json:"evaluation_detail"` // 採点結果の詳細メモ。採点ロジックが増えても柔軟に持てるように JSONB で保存。
	CreatedAt        time.Time              `json:"created_at"`        // DB 側で決まる投稿時刻。履歴ソートや期間集計に必須です。
}

// LeaderboardRow はランキング結果の1行を表します。
type LeaderboardRow struct {
	UserID    *int      `json:"user_id"`
	Username  string    `json:"username"`
	BestScore int       `json:"best_score"`
	Attempts  int       `json:"attempts"`
	LastAt    time.Time `json:"last_at"`
}

// ScoresRepository は scores テーブルに対する操作を定義するインターフェースです。
// インターフェース→関数の定義だけをここに置き、実装は下に続く構造体で行います。
type ScoresRepository interface {
	// Create は新しいスコアレコードをデータベースに保存します。
	// HTTP ハンドラなど上位レイヤーから呼ばれ、1 回のプレイ結果を 1 行として登録します。
	Create(ctx context.Context, record *Score) error

	// FindLeaderboard は指定された期間と上限数でランキングを取得します。
	// period: "day", "week", "all" のいずれか
	// 期間によって SQL の WHERE 条件を差し替え、上位 n 件だけ返します。
	FindLeaderboard(ctx context.Context, period string, limit int) ([]LeaderboardRow, error)

	// FindUserScores は指定されたユーザーのスコア履歴を取得します。
	// マイページなどで最新の解答履歴を表示する用途を想定しています。
	FindUserScores(ctx context.Context, userID int, limit int) ([]Score, error)
}

// scoresRepo は ScoresRepository の実装です。
// これはデータベースと接続するための道具
type scoresRepo struct {
	db *sql.DB
}

// NewScoresRepository は ScoresRepository の新しいインスタンスを作成します。
// インターフェース型を返すことで、テスト時には別実装に差し替えることも可能になります。
// scoresRepoが作業員だとしたら、NewScoresRepositoryはその作業員を派遣するための工場のようなものです。
// ↑NewScoresRepositoryという工場に、dbという作業員を渡す
// それにより、ScoresRepositoryインターフェースを満たす作業員が得られます。
// ↑つまり、これをmainなどで呼ぶことで、インターフェースの中の関数が使えるようになる。
func NewScoresRepository(db *sql.DB) ScoresRepository {
	return &scoresRepo{db: db}
}

// Create では新しいスコアレコードをデータベースに保存します。
// この関数定義はscoresRepoという構造体にCreateメソッドを実装しています。
// ここでCreateの具体的な作業内容を実装する。
func (r *scoresRepo) Create(ctx context.Context, record *Score) error {
	// evaluation_detail を JSONB に変換
	// Go の map はそのままでは PostgreSQL の JSONB 列に入れられないため、
	// いったん JSON 文字列（[]byte）にシリアライズしてから保存します。
	var detailJSON interface{}
	if record.EvaluationDetail != nil {
		jsonBytes, err := json.Marshal(record.EvaluationDetail)
		if err != nil {
			return fmt.Errorf("failed to marshal evaluation_detail: %w", err)
		}
		detailJSON = jsonBytes
	} else {
		// nil の場合は SQL の NULL として扱う
		detailJSON = nil
	}

	// SQL はヒアドキュメントで書くと列の並びが視覚的に追いやすくなります。
	// このVALUES部分には、後で実際の値を渡す。
	query := `
		INSERT INTO scores (
			user_id, question_id, prompt, ai_response, score,
			model_vendor, model_name, answer_number, latency_ms, evaluation_detail,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	// QueryRowContext は 1 行だけ返るクエリを実行し、Scan で構造体に詰めます。
	// returning で ID と created_at を受け取ることで、呼び出し側が直後に利用できます。
	err := r.db.QueryRowContext(
		ctx,
		query,
		record.UserID,
		record.QuestionID,
		record.Prompt,
		record.AIResponse,
		record.Score,
		record.ModelVendor,
		record.ModelName,
		record.AnswerNumber,
		record.LatencyMs,
		detailJSON,
		time.Now(), // Go 側で現在時刻をセットしておくと、呼び出しが終わった時点で値が分かります。
	).Scan(&record.ID, &record.CreatedAt)
	if err != nil {
		// %w を使うと元のエラー情報を保ったまま上位に伝えられ、原因調査が容易になります。
		return fmt.Errorf("failed to insert score: %w", err)
	}

	return nil
}

// FindLeaderboard は指定された期間と上限数でランキングを取得します。
func (r *scoresRepo) FindLeaderboard(ctx context.Context, period string, limit int) ([]LeaderboardRow, error) {
	// 期間条件を決定
	// SQL の断片を直接差し込むため、switch で許可する文字列のみ選ぶと安全です。
	var whereClause string
	switch period {
	case "day":
		whereClause = "AND s.created_at >= NOW() - INTERVAL '1 day'"
	case "week":
		whereClause = "AND s.created_at >= NOW() - INTERVAL '7 days'"
	case "all":
		whereClause = ""
	default:
		return nil, fmt.Errorf("invalid period: %s (must be day, week, or all)", period)
	}

	// 期間条件は whereClause に差し込み、ランキングのスコアと最新回答日時で並べ替えます。
	query := fmt.Sprintf(`
		SELECT
			s.user_id,
			COALESCE(u.username, 'guest') as username,
			MAX(s.score) as best_score,
			COUNT(*) as attempts,
			MAX(s.created_at) as last_at
		FROM scores s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE 1=1 %s
		GROUP BY s.user_id, u.username
		ORDER BY best_score DESC, last_at DESC
		LIMIT $1
	`, whereClause)

	// QueryContext で複数行取得。limit はバインド変数で渡すことで SQL インジェクションを防ぎます。
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// rows.Close のエラーは通常無視しても問題ないが、linter 対策のためログ出力を想定
			_ = closeErr
		}
	}()

	var results []LeaderboardRow
	for rows.Next() {
		var row LeaderboardRow
		err := rows.Scan(
			&row.UserID,
			&row.Username,
			&row.BestScore,
			&row.Attempts,
			&row.LastAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard row: %w", err)
		}
		// 1 行読み取るたびにスライスへ追加。順位は ORDER BY の結果そのままです。
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("leaderboard rows iteration error: %w", err)
	}

	// 呼び出し側は []LeaderboardRow をそのまま JSON にして返すもよし、テンプレートに渡すもよしです。
	return results, nil
}

// FindUserScores は指定されたユーザーのスコア履歴を取得します。
func (r *scoresRepo) FindUserScores(ctx context.Context, userID int, limit int) ([]Score, error) {
	// ORDER BY created_at DESC で最新回答から順に並べ替えます。
	query := `
		SELECT
			id, user_id, question_id, prompt, ai_response, score,
			model_vendor, model_name, answer_number, latency_ms,
			evaluation_detail, created_at
		FROM scores
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	// userID と limit は必ずプレースホルダを使って安全に渡します。
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query user scores: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// rows.Close のエラーは通常無視しても問題ないが、linter 対策のためログ出力を想定
			_ = closeErr
		}
	}()

	var results []Score
	for rows.Next() {
		var s Score
		var detailJSON []byte

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.QuestionID,
			&s.Prompt,
			&s.AIResponse,
			&s.Score,
			&s.ModelVendor,
			&s.ModelName,
			&s.AnswerNumber,
			&s.LatencyMs,
			&detailJSON,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan score row: %w", err)
		}

		// JSONB を map に変換
		if len(detailJSON) > 0 {
			// DB から返る JSONB は []byte なので、構造体の map に戻します。
			// エラーがあると詳細表示ができないため、そのまま上位に返します。
			if err := json.Unmarshal(detailJSON, &s.EvaluationDetail); err != nil {
				return nil, fmt.Errorf("failed to unmarshal evaluation_detail: %w", err)
			}
		}

		// 1 レコードずつ結果スライスに詰めていきます。limit が小さければメモリ消費も抑えられます。
		results = append(results, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("user scores rows iteration error: %w", err)
	}

	return results, nil
}
