// Package main provides a test utility for verifying Gemini API connectivity and response times.
package main

import (
	"context"
	"log"
	"time"

	"github.com/shiv/CoT_game/backend/internal/ai"
)

func main() {
	log.Println("Gemini API接続テストを開始...")

	client, err := ai.NewGeminiClientFromEnv()
	if err != nil {
		log.Fatalf("クライアント作成エラー: %v", err)
	}

	// テストする問題一覧
	problems := []struct {
		name   string
		prompt string
	}{
		{"1+1", `以下の問題を解いてください。

【問題】
1+1を計算してください。

【指示】
答えを教えて

必ず最終的な答えだけを出力してください。`},
		{"strawberry", `以下の問題を解いてください。

【問題】
strawberryの中にrは何個ある？

【指示】
答えを教えて

必ず最終的な答えだけを出力してください。`},
	}

	for _, p := range problems {
		log.Printf("\n=== テスト: %s ===", p.name)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		start := time.Now()
		response, err := client.GenerateAnswer(ctx, p.prompt)
		elapsed := time.Since(start)

		if err != nil {
			log.Printf("❌ エラー発生 (経過時間: %v): %v", elapsed, err)
			cancel()
			continue
		}

		log.Printf("✅ 成功 (経過時間: %v)", elapsed)
		log.Printf("回答: %s", response)
		cancel()
	}
}
