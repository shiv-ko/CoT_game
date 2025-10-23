// Package eval は AI の回答を正解と照合し、差分に応じたスコアと詳細情報を算出するロジックをまとめたパッケージ。
package eval

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	// toleranceThreshold は「ここまでの誤差なら高得点を与える」目安の値（0.01）。
	// 従来の 0.01 差で 95 点という仕様を、連続スコア計算に引き継ぐための境界として扱う。
	toleranceThreshold = 0.01
)

var (
	// scoreDecayRate は指数関数 e^(-k*diff) の k を調整した値。
	// diff が 0.01 (toleranceThreshold) のときに約 95 点になるよう求めている。
	scoreDecayRate = -math.Log(0.95) / toleranceThreshold
	// numberPattern は回答文から最初に現れる数値らしい文字列を抜き出すためのパターン。
	// プラス・マイナス符号や小数点と整数部の組を想定したシンプルな正規表現にとどめている。
	numberPattern = regexp.MustCompile(`[-+]?\d+(?:\.\d+)?`)
)

// Evaluate は AI 回答を正解と比較し、数値の一致度に基づくスコア・抽出値・メタ情報を返す。
// 完全一致なら 100 点、それ以外は数値誤差に応じて連続的に減点し、必要な場合はモード名と理由も detail に記録する。
func Evaluate(answerText string, correct string) (score int, extracted *float64, mode string, detail map[string]any) {
	// 回答と正解の前後スペースを除去し、純粋な値として比較しやすくする。
	trimmedAnswer := strings.TrimSpace(answerText)
	trimmedCorrect := strings.TrimSpace(correct)

	// detail は評価過程の情報を溜め込むメタデータ。
	// デバッグや可視化で使えるよう、生の入力・前処理結果・スコア計算式を格納している。
	detail = map[string]any{
		"answer_raw":       answerText,
		"correct_raw":      correct,
		"answer_trimmed":   trimmedAnswer,
		"correct_trimmed":  trimmedCorrect,
		"score_strategy":   "score=100*exp(-k*diff)",
		"score_decay_rate": scoreDecayRate,
	}

	if trimmedAnswer == trimmedCorrect {
		// 完全一致の場合は文字列比較で終わらせ、満点と理由を保存する。
		// 数値として解釈できるなら diff=0 をメタ情報に積み増す。
		mode = "exact_match"
		score = 100
		detail["mode_reason"] = "回答文字列が正解と完全一致"
		detail["normalized_score"] = score
		if val, err := strconv.ParseFloat(trimmedAnswer, 64); err == nil {
			detail["extracted_text"] = trimmedAnswer
			detail["extracted_numeric"] = val
			detail["numeric_diff"] = 0.0
			valueCopy := val
			extracted = &valueCopy
		}
		return
	}

	// 正解文字列が数値として読めるか先に調べ、後続の誤差計算に備える。
	correctVal, correctErr := strconv.ParseFloat(trimmedCorrect, 64)
	if correctErr == nil {
		detail["correct_numeric"] = correctVal
	}

	// 回答内から数値らしき部分を抽出。見つからなければスコア 0 で終了。
	matched := numberPattern.FindString(trimmedAnswer)
	if matched == "" {
		mode = "no_numeric"
		score = 0
		detail["mode_reason"] = "回答に数値が含まれていない"
		return
	}

	// 正規表現で得た文字列を float に変換できるか確認する。
	parsed, err := strconv.ParseFloat(matched, 64)
	if err != nil {
		mode = "no_numeric"
		score = 0
		detail["mode_reason"] = "数値抽出に失敗"
		detail["extracted_text"] = matched
		detail["parse_error"] = err.Error()
		return
	}

	// 数値抽出に成功した場合は、テキスト版と数値版を detail に記録する。
	detail["extracted_text"] = matched
	detail["extracted_numeric"] = parsed
	valueCopy := parsed
	extracted = &valueCopy

	if correctErr == nil {
		// 正解も数値なら calculatePreciseDiff で高精度に差分を算出し、連続スコアを決定する。
		diff, precise := calculatePreciseDiff(trimmedCorrect, matched)
		detail["diff_precision"] = "float"
		if precise {
			detail["diff_precision"] = "rational"
		} else {
			diff = math.Abs(parsed - correctVal)
		}

		detail["numeric_diff"] = diff
		score = computeScore(diff)
		detail["normalized_score"] = score
		if diff == 0 {
			mode = "numeric_exact"
			detail["mode_reason"] = "数値比較で誤差が0"
		} else {
			mode = "numeric_score"
			detail["mode_reason"] = "数値誤差に基づき連続スコアを算出"
		}
		return
	}

	mode = "extracted_only"
	score = 0
	detail["mode_reason"] = "数値抽出は成功したが正解が数値として解釈できない"

	return
}

// calculatePreciseDiff は文字列表現を用いた高精度計算で誤差を算出する。
// 変換に失敗した場合は (0, false) を返し、呼び出し側でフォールバックしてもらう。
func calculatePreciseDiff(correctStr string, extractedStr string) (diff float64, ok bool) {
	// SetString は文字列表現をそのまま有理数に変換できる場合のみ true を返す。
	correctRat, ok := new(big.Rat).SetString(correctStr)
	if !ok {
		return 0, false
	}
	// 抽出側についても同様に有理数化し、成功した場合だけ高精度計算を行う。
	extractedRat, ok := new(big.Rat).SetString(extractedStr)
	if !ok {
		return 0, false
	}
	// 差分 diffRat を計算し、負であれば絶対値に変換する。
	diffRat := new(big.Rat).Sub(extractedRat, correctRat)
	if diffRat.Sign() < 0 {
		diffRat.Neg(diffRat)
	}
	// 呼び出し側で扱いやすいよう float64 へ変換しつつ、高精度判定の成功を示す。
	diff, _ = diffRat.Float64()
	return diff, true
}

// computeScore は「誤差 diff が大きいほど急速に減点される」指数グラフで連続スコアを計算する。
// diff が 0 なら 100 点、0.01 で約 95 点、diff が大きくなるほど 0 点へ近づく。
func computeScore(diff float64) int {
	if diff <= 0 {
		return 100
	}
	raw := 100 * math.Exp(-scoreDecayRate*diff)
	if raw < 0 {
		raw = 0
	}
	if raw > 100 {
		raw = 100
	}
	return int(math.Round(raw))
}
