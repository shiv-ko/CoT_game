// evaluator_test.go は Evaluate 関数の連続スコア計算とモード判定を網羅的に確認する単体テストをまとめたファイル。
package eval

import (
	"math"
	"testing"
)

// floatPtr は即席で *float64 を作るためのヘルパー。テーブル定義内で使えるよう切り出している。
func floatPtr(v float64) *float64 {
	return &v
}

// TestEvaluate は評価ロジックの主要ケースを網羅し、計算された diff と score が整合するかを検証する。
// テーブルテスト形式で入力と期待値を並べ、初心者でも "ケースごとの条件 → 期待結果" が追いやすい構成にしている。
func TestEvaluate(t *testing.T) {
	testcases := []struct {
		name             string
		answer           string
		correct          string
		expectScore      int
		expectMode       string
		expectExtracted  *float64
		expectDiff       *float64
		expectNormalized bool
	}{
		{
			name:             "ExactMatch",
			answer:           "42",
			correct:          "42",
			expectScore:      100,
			expectMode:       "exact_match",
			expectExtracted:  floatPtr(42.0),
			expectDiff:       floatPtr(0.0),
			expectNormalized: true,
		},
		{
			name:             "NumericExactDifferentFormatting",
			answer:           "010",
			correct:          "10",
			expectScore:      100,
			expectMode:       "numeric_exact",
			expectExtracted:  floatPtr(10.0),
			expectDiff:       floatPtr(0.0),
			expectNormalized: true,
		},
		{
			name:             "ToleranceWithin",
			answer:           "The result is 9.991",
			correct:          "10",
			expectScore:      computeScore(0.009),
			expectMode:       "numeric_score",
			expectExtracted:  floatPtr(9.991),
			expectDiff:       floatPtr(0.009),
			expectNormalized: true,
		},
		{
			name:             "ToleranceOutside",
			answer:           "9.989",
			correct:          "10",
			expectScore:      computeScore(0.011),
			expectMode:       "numeric_score",
			expectExtracted:  floatPtr(9.989),
			expectDiff:       floatPtr(0.011),
			expectNormalized: true,
		},
		{
			name:        "NoNumeric",
			answer:      "No numbers here",
			correct:     "10",
			expectScore: 0,
			expectMode:  "no_numeric",
		},
		{
			name:             "MultipleNumbersUsesFirst",
			answer:           "First 11 then 10",
			correct:          "10",
			expectScore:      computeScore(1.0),
			expectMode:       "numeric_score",
			expectExtracted:  floatPtr(11.0),
			expectDiff:       floatPtr(1.0),
			expectNormalized: true,
		},
		{
			name:             "NegativeExact",
			answer:           "-5",
			correct:          "-5",
			expectScore:      100,
			expectMode:       "exact_match",
			expectExtracted:  floatPtr(-5.0),
			expectDiff:       floatPtr(0.0),
			expectNormalized: true,
		},
		{
			name:             "DecimalPrecision",
			answer:           "3.1415",
			correct:          "3.1416",
			expectScore:      computeScore(0.0001),
			expectMode:       "numeric_score",
			expectExtracted:  floatPtr(3.1415),
			expectDiff:       floatPtr(0.0001),
			expectNormalized: true,
		},
		{
			name:        "EmptyAnswer",
			answer:      "",
			correct:     "10",
			expectScore: 0,
			expectMode:  "no_numeric",
		},
		{
			name:             "ExtractedButCorrectNonNumeric",
			answer:           "Value: 7",
			correct:          "ten",
			expectScore:      0,
			expectMode:       "extracted_only",
			expectExtracted:  floatPtr(7.0),
			expectNormalized: false,
		},
	}

	for _, tc := range testcases {
		// tc をローカル変数として再代入することで、サブテスト内のクロージャがループ変数を正しく捕捉できるようにする。
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel は各サブテストを並列実行し、テスト全体の実行時間を短縮する。
			// 共有資源に触れないユニットテストなので安全に並列化できる。
			t.Parallel()

			// Evaluate を実行し、返ってきたスコアやモード、detail 情報をケースごとに照合する。
			score, extracted, mode, detail := Evaluate(tc.answer, tc.correct)

			// スコアは整数一致で比較。期待と異なる場合は Fatalf で即テスト終了。
			if score != tc.expectScore {
				t.Fatalf("score mismatch: got %d, want %d", score, tc.expectScore)
			}

			// mode も同様に文字列比較で検証し、分岐ミスを検出する。
			if mode != tc.expectMode {
				t.Fatalf("mode mismatch: got %s, want %s", mode, tc.expectMode)
			}

			// 抽出値の検証。ポインタが nil かどうかも含めて期待値と比較する。
			// nil と非 nil を明示的にチェックすることで、値が存在すべきか否かの仕様を確認できる。
			if tc.expectExtracted == nil {
				if extracted != nil {
					t.Fatalf("expected nil extracted, got %v", *extracted)
				}
			} else {
				if extracted == nil {
					t.Fatal("expected extracted value, got nil")
				}
				// 抽出された数値が期待値とほぼ等しいか確認。浮動小数点の比較なので厳密な等価ではなく微小差を許容する。
				if math.Abs(*extracted-*tc.expectExtracted) > 1e-9 {
					t.Fatalf("extracted mismatch: got %f, want %f", *extracted, *tc.expectExtracted)
				}
			}

			// detail マップに入った numeric_diff は、評価ロジックの内部状態を確認する指標。
			// ケースによっては存在しない（no_numeric など）ため、期待に応じて有無と数値をチェックする。
			diffValue, diffOK := detail["numeric_diff"].(float64)
			if tc.expectDiff == nil {
				if diffOK {
					t.Fatalf("unexpected numeric_diff detail: %v", diffValue)
				}
			} else {
				if !diffOK {
					t.Fatal("expected numeric_diff detail, but missing")
				}
				if math.Abs(diffValue-*tc.expectDiff) > 1e-9 {
					t.Fatalf("numeric_diff mismatch: got %f, want %f", diffValue, *tc.expectDiff)
				}
				// computeScore を再計算し、detail 上の diff と返却された score の両方が同じ式に従っているか確認する。
				expectedScore := computeScore(diffValue)
				if score != expectedScore {
					t.Fatalf("score should align with computeScore(diff): diff=%f got=%d want=%d", diffValue, score, expectedScore)
				}
			}

			// normalized_score は detail に保存されたスコアのコピー。
			// 正解が数値として扱えたケースでは設定されるはずなので、有無と値の一致を検証する。
			normalizedVal, hasNormalized := detail["normalized_score"].(int)
			if tc.expectNormalized {
				if !hasNormalized {
					t.Fatal("expected normalized_score detail, but missing")
				}
				if normalizedVal != score {
					t.Fatalf("normalized_score mismatch: got %d, want %d", normalizedVal, score)
				}
			} else if hasNormalized {
				t.Fatalf("unexpected normalized_score detail: %v", normalizedVal)
			}
		})
	}
}
