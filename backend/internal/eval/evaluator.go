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
	// r0 は相対誤差がこの値のときにスコアが50点になる基準値
	// スケール適応型の場合は使用されない
	r0 = 0.05

	// p はロジスティック関数の曲率パラメータ（デフォルト: 2.0）
	// 値が大きいほど減衰が急になる
	p = 2.0

	// tolAbsHint は正解が0に近い場合のスケール補助値
	// 相対誤差を計算する際の分母が0にならないようにする
	tolAbsHint = 1e-2

	// 整数スケール問題の閾値
	// 正解の絶対値がこの値以下の場合、整数ベースのスコアリングを使用
	integerScaleThreshold = 1000.0

	// 整数スケールでの基準誤差数（この誤差で50点になる）
	integerBaseError = 2.0
)

// numberPattern は回答文から数値らしい文字列を抜き出すためのパターン。
// プラス・マイナス符号や小数点と整数部の組を想定したシンプルな正規表現にとどめている。
// 複数の数値がある場合は、最後に出現するものを最終回答として採用する。
var numberPattern = regexp.MustCompile(`[-+]?\d+(?:\.\d+)?`)

// Evaluate は AI 回答を正解と比較し、数値の一致度に基づくスコア・抽出値・メタ情報を返す。
// 完全一致なら 100 点、それ以外は数値誤差に応じて連続的に減点し、必要な場合はモード名と理由も detail に記録する。
func Evaluate(answerText string, correct string) (score int, extracted *float64, mode string, detail map[string]any) {
	// 回答と正解の前後スペースを除去し、純粋な値として比較しやすくする。
	trimmedAnswer := strings.TrimSpace(answerText)
	trimmedCorrect := strings.TrimSpace(correct)

	// detail は評価過程の情報を溜め込むメタデータ。
	// デバッグや可視化で使えるよう、生の入力・前処理結果・スコア計算式を格納している。
	detail = map[string]any{
		"answer_raw":      answerText,
		"correct_raw":     correct,
		"answer_trimmed":  trimmedAnswer,
		"correct_trimmed": trimmedCorrect,
		"score_strategy":  "v3: スケール適応型（整数問題は絶対誤差、大きな数は相対誤差）",
		"r0":              r0,
		"p":               p,
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
			detail["absolute_diff"] = 0.0
			detail["relative_error"] = 0.0
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

	// 回答内から数値らしき部分を抽出。複数ある場合は最後のものを採用。
	// AIの説明的な回答では、最終的な答えが文末に来ることが多いため。
	allMatches := numberPattern.FindAllString(trimmedAnswer, -1)
	if len(allMatches) == 0 {
		mode = "no_numeric"
		score = 0
		detail["mode_reason"] = "回答に数値が含まれていない"
		return
	}
	matched := allMatches[len(allMatches)-1] // 最後の数値を採用
	if len(allMatches) > 1 {
		detail["all_numbers_found"] = allMatches
		detail["extraction_note"] = "複数の数値が見つかったため、最後のものを最終回答として採用"
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

		detail["absolute_diff"] = diff

		if diff == 0 {
			mode = "numeric_exact"
			score = 100
			detail["mode_reason"] = "数値比較で誤差が0"
			detail["normalized_score"] = score
			return
		}

		// 正解の大きさに応じてスコアリング方式を選択
		absCorrect := math.Abs(correctVal)
		if absCorrect <= integerScaleThreshold {
			// 整数スケール問題：絶対誤差ベースでスコアリング
			score = computeIntegerScaleScore(diff, absCorrect)
			mode = "numeric_score_integer"
			detail["mode_reason"] = "整数スケール問題として絶対誤差ベースで評価（v3）"
			detail["scale_type"] = "integer"
			detail["base_error"] = integerBaseError
		} else {
			// 大きな数の問題：相対誤差ベースでスコアリング
			relativeError := calculateRelativeError(diff, correctVal)
			detail["relative_error"] = relativeError
			score = computeScore(relativeError)
			mode = "numeric_score_relative"
			detail["mode_reason"] = "大規模数値問題として相対誤差ベースで評価（v3）"
			detail["scale_type"] = "relative"
		}
		detail["normalized_score"] = score
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

// calculateRelativeError は絶対誤差と正解値から相対誤差を計算する。
// 正解が0に近い場合は tolAbsHint を使用して分母が0にならないようにする。
func calculateRelativeError(absoluteDiff float64, correctVal float64) float64 {
	denominator := math.Max(math.Abs(correctVal), tolAbsHint)
	return absoluteDiff / denominator
}

// computeIntegerScaleScore は整数スケールの問題に対して絶対誤差ベースでスコアを計算する。
// 正解が小さい整数（例: 3）の場合、相対誤差ではなく絶対誤差で評価する方が適切。
//
// スコアリング方式:
// - 誤差0: 100点
// - 誤差1: 75点（1つずれ）
// - 誤差2: 50点（2つずれ、基準値）
// - 誤差3: 30点
// - 誤差4: 18点
// - 誤差5以上: 10点以下に急速に減少
//
// 数式: score = 100 / (1 + (diff/base)^p)
// ここで base = integerBaseError (デフォルト: 2.0)
func computeIntegerScaleScore(absoluteDiff float64, correctVal float64) int {
	if absoluteDiff <= 0 {
		return 100
	}

	// 正解が極小（1未満）の場合は、より厳しい評価
	base := integerBaseError
	if correctVal < 1.0 {
		base = 0.5 // 0.5の誤差で50点
	}

	// ロジスティック関数: score = 100 / (1 + (diff/base)^p)
	ratio := absoluteDiff / base
	raw := 100.0 / (1.0 + math.Pow(ratio, p))

	// スコアを0-100の範囲に制限
	if raw < 0 {
		raw = 0
	}
	if raw > 100 {
		raw = 100
	}

	return int(math.Round(raw))
}

// computeScore は相対誤差に基づいてロジスティック関数でスコアを計算する。
// rel が 0 なら 100 点、r0(5%)で 50 点、それ以上は急速に減少する。
// スコア式: score = 100 / (1 + (rel/r0)^p)
func computeScore(relativeError float64) int {
	if relativeError <= 0 {
		return 100
	}

	// ロジスティック関数: score = 100 / (1 + (rel/r0)^p)
	ratio := relativeError / r0
	raw := 100.0 / (1.0 + math.Pow(ratio, p))

	// スコアを0-100の範囲に制限
	if raw < 0 {
		raw = 0
	}
	if raw > 100 {
		raw = 100
	}

	return int(math.Round(raw))
}
