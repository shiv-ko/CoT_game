// scores_repo_test.go は scores リポジトリの DB 実装が想定通りに動くかを結合テストで検証し、リファクタリング時の安全網になります。
// 関数だけじゃなく、DB接続も含んでいるので単体テストではなく結合テスト
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// setupTestDB は Docker/PostgreSQL への疎通を確認しつつ *sql.DB を返すヘルパーで、準備できていない環境ではテストを安全にスキップさせます。
// 「テストが落ちる=アプリが壊れている」と混同しないように、環境未準備の場合は落とさずスキップする運用です。
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// 環境変数 TEST_DATABASE_URL があればそれを使用し、なければローカル開発用のデフォルト値を使用します。
	// これにより、CI 環境など異なるDB設定にも柔軟に対応できます。
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=cot_game sslmode=disable"
	}

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

// ensureTestUser は指定IDのユーザーが users テーブルに存在するよう強制し、外部キー制約でテストが失敗しないようにします。
func ensureTestUser(t *testing.T, db *sql.DB, userID int) {
	t.Helper()
	username := fmt.Sprintf("test_user_%d", userID)
	passwordHash := "test_hash"

	// ユーザーが存在しない場合のみ作成
	_, err := db.Exec(`
		INSERT INTO users (id, username, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO NOTHING
	`, userID, username, passwordHash)
	if err != nil {
		t.Logf("テストユーザー作成警告: %v", err)
	}
}

// cleanupTestData は scores と users からテスト用レコードを削除し、テストケース同士の干渉を防ぎます。
func cleanupTestData(t *testing.T, db *sql.DB, userID int) {
	t.Helper()
	_, _ = db.Exec("DELETE FROM scores WHERE user_id = $1", userID)
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", userID)
}

// TestScoresRepo_Create は Create メソッドがスコアを INSERT し、ID/CreatedAt をセットできるかを検証します。
// Create関数をテストしているから、単体テストに見えるが、DBまでテストしているので結合テストになる。
func TestScoresRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	testUserID := 999999              // テスト用のユーザーID
	ensureTestUser(t, db, testUserID) // テストユーザーを作成
	defer cleanupTestData(t, db, testUserID)

	tests := []struct {
		name    string
		record  *Score
		wantErr bool
	}{
		{
			name: "正常系：完全なレコード",
			record: &Score{
				UserID:       &testUserID,
				QuestionID:   1,
				Prompt:       "テストプロンプト",
				AIResponse:   "テスト回答",
				Score:        100,
				ModelVendor:  "gemini",
				ModelName:    stringPtr("gemini-1.5-flash"),
				AnswerNumber: float64Ptr(42.0),
				LatencyMs:    1500,
				EvaluationDetail: map[string]interface{}{
					"mode":         "exact_match",
					"mode_reason":  "完全一致",
					"numeric_diff": 0.0,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系：最小限のレコード",
			record: &Score{
				UserID:      &testUserID,
				QuestionID:  1,
				Prompt:      "最小プロンプト",
				AIResponse:  "最小回答",
				Score:       50,
				ModelVendor: "gemini",
				LatencyMs:   500,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 正常時はIDとタイムスタンプが設定されることを確認
				// ここまでチェックしておくと、INSERT が実際に成功しているか／RETURNING で値を受け取れているかが分かります。
				if tt.record.ID == 0 {
					t.Error("Create() did not set ID")
				}
				if tt.record.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
			}
		})
	}
}

// TestScoresRepo_FindLeaderboard は期間・件数指定でランキング取得ができ、期間指定ミス時にエラーを返すことを確認します。
func TestScoresRepo_FindLeaderboard(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	testUserID1 := 999997
	testUserID2 := 999998
	ensureTestUser(t, db, testUserID1) // テストユーザーを作成
	ensureTestUser(t, db, testUserID2) // テストユーザーを作成
	defer cleanupTestData(t, db, testUserID1)
	defer cleanupTestData(t, db, testUserID2)

	// テストデータを挿入
	// ランキングは集計結果を返すだけなので、最低でも 2 ユーザー分のスコアがあると挙動が見えやすいです。
	_ = repo.Create(ctx, &Score{
		UserID:      &testUserID1,
		QuestionID:  1,
		Prompt:      "test1",
		AIResponse:  "ans1",
		Score:       95,
		ModelVendor: "gemini",
		LatencyMs:   100,
	})
	_ = repo.Create(ctx, &Score{
		UserID:      &testUserID2,
		QuestionID:  1,
		Prompt:      "test2",
		AIResponse:  "ans2",
		Score:       85,
		ModelVendor: "gemini",
		LatencyMs:   200,
	})

	tests := []struct {
		name    string
		period  string
		limit   int
		wantErr bool
	}{
		{
			name:    "全期間ランキング",
			period:  "all",
			limit:   10,
			wantErr: false,
		},
		{
			name:    "1日間ランキング",
			period:  "day",
			limit:   10,
			wantErr: false,
		},
		{
			name:    "不正な期間",
			period:  "invalid",
			limit:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := repo.FindLeaderboard(ctx, tt.period, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLeaderboard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.period == "all" {
				// 正常系で結果が取得できることを確認（テストデータが入っている場合）
				if len(rows) == 0 {
					t.Log("FindLeaderboard() returned empty results (may be expected if no test data)")
				} else {
					t.Logf("FindLeaderboard() returned %d rows", len(rows))
					// 最初の行が最高スコアであることを確認
					for i := 1; i < len(rows); i++ {
						if rows[i-1].BestScore < rows[i].BestScore {
							t.Errorf("Leaderboard not sorted: row %d score=%d < row %d score=%d",
								i-1, rows[i-1].BestScore, i, rows[i].BestScore)
						}
					}
				}
			}
		})
	}
}

// TestScoresRepo_FindUserScores はユーザー別スコア履歴が新しい順で返ることと、存在しないユーザーでもエラーにしないことをテストします。
func TestScoresRepo_FindUserScores(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	testUserID := 999996
	ensureTestUser(t, db, testUserID) // テストユーザーを作成
	defer cleanupTestData(t, db, testUserID)

	// 既存データをクリーンアップしてから新規挿入
	cleanupTestData(t, db, testUserID)
	ensureTestUser(t, db, testUserID) // クリーンアップ後に再作成

	// テストデータを2件挿入
	// created_at の降順を確認するため、少なくとも 2 レコード必要です。
	err1 := repo.Create(ctx, &Score{
		UserID:      &testUserID,
		QuestionID:  1, // 既存の問題IDを使用
		Prompt:      "prompt1",
		AIResponse:  "ans1",
		Score:       100,
		ModelVendor: "gemini",
		LatencyMs:   100,
	})
	if err1 != nil {
		t.Fatalf("Failed to create first score: %v", err1)
	}

	time.Sleep(10 * time.Millisecond) // created_at の順序を保証

	err2 := repo.Create(ctx, &Score{
		UserID:      &testUserID,
		QuestionID:  1, // 同じ問題IDを使用
		Prompt:      "prompt2",
		AIResponse:  "ans2",
		Score:       90,
		ModelVendor: "gemini",
		LatencyMs:   200,
	})
	if err2 != nil {
		t.Fatalf("Failed to create second score: %v", err2)
	}

	tests := []struct {
		name    string
		userID  int
		limit   int
		wantErr bool
		minRows int
	}{
		{
			name:    "ユーザースコア取得",
			userID:  testUserID,
			limit:   10,
			wantErr: false,
			minRows: 2,
		},
		{
			name:    "存在しないユーザー",
			userID:  888888,
			limit:   10,
			wantErr: false,
			minRows: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores, err := repo.FindUserScores(ctx, tt.userID, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(scores) < tt.minRows {
					t.Errorf("FindUserScores() returned %d rows, want at least %d", len(scores), tt.minRows)
				}

				// 降順でソートされていることを確認
				for i := 1; i < len(scores); i++ {
					if scores[i-1].CreatedAt.Before(scores[i].CreatedAt) {
						t.Error("FindUserScores() not sorted by created_at DESC")
					}
				}
			}
		})
	}
}

// stringPtr はリテラル文字列のポインタを返し、テストデータで *string フィールドへ簡潔に値を入れるための小道具です。
func stringPtr(s string) *string {
	// Go では文字列リテラルのポインタを直接取れないため、簡単なヘルパーを用意しておくとテストが読みやすくなります。
	return &s
}

// float64Ptr は数値リテラルのポインタ版ヘルパーで、AnswerNumber など *float64 フィールドの値定義を簡潔にします。
func float64Ptr(f float64) *float64 {
	// ↑ と同じ理由で float64 版も用意しています。
	return &f
}
