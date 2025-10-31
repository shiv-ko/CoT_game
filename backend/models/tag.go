// Package models はデータベースや外部入出力に対応するドメインオブジェクトを提供します。
package models

// Tag はタグの定義情報を表す構造体です。
// 各タグは問題の分類と、AIへのプロンプトヒントを提供します。
type Tag struct {
	ID          string `json:"id"`          // タグの一意識別子（例: "calculation"）
	Label       string `json:"label"`       // 表示用ラベル（例: "計算問題"）
	Icon        string `json:"icon"`        // アイコン絵文字（例: "🔢"）
	Description string `json:"description"` // タグの説明
	PromptTips  string `json:"prompt_tips"` // AI向けプロンプトヒント
	Color       string `json:"color"`       // 表示用カラーコード（例: "#3B82F6"）
}

// TagDefinitions はタグIDからタグメタデータへのマッピングです。
// フロントエンドとバックエンドで一貫したタグ情報を提供します。
var TagDefinitions = map[string]Tag{
	"calculation": {
		ID:          "calculation",
		Label:       "計算問題",
		Icon:        "🔢",
		Description: "数値計算や算術演算が必要な問題",
		PromptTips:  "数値を正確に計算させるため、段階的な計算プロセスを促すプロンプトが有効です。",
		Color:       "#3B82F6", // blue-500
	},
	"character_counting": {
		ID:          "character_counting",
		Label:       "文字数カウント",
		Icon:        "📊",
		Description: "文字列の長さや特定文字の出現回数を数える問題",
		PromptTips:  "文字を一つずつ数えるように指示し、確認を促すプロンプトが効果的です。",
		Color:       "#10B981", // green-500
	},
	"text_analysis": {
		ID:          "text_analysis",
		Label:       "テキスト解析",
		Icon:        "📝",
		Description: "文章の構造や内容を分析する問題",
		PromptTips:  "テキストを丁寧に読み解くよう促し、分析の観点を明示するプロンプトが有効です。",
		Color:       "#8B5CF6", // violet-500
	},
	"text_problem": {
		ID:          "text_problem",
		Label:       "文章題",
		Icon:        "📖",
		Description: "文章から情報を読み取り、問題を解く必要がある問題",
		PromptTips:  "問題文を正確に理解させ、条件を整理してから解かせるプロンプトが効果的です。",
		Color:       "#F59E0B", // amber-500
	},
	"pattern_recognition": {
		ID:          "pattern_recognition",
		Label:       "パターン認識",
		Icon:        "🔍",
		Description: "規則性やパターンを見つけ出す問題",
		PromptTips:  "観察と仮説検証を促し、複数の例から規則を導き出すプロンプトが有効です。",
		Color:       "#EC4899", // pink-500
	},
	"logic_puzzle": {
		ID:          "logic_puzzle",
		Label:       "論理パズル",
		Icon:        "🧩",
		Description: "論理的思考や推論が必要なパズル問題",
		PromptTips:  "論理的に段階を踏んで考えるよう促し、矛盾を確認するプロンプトが効果的です。",
		Color:       "#EF4444", // red-500
	},
	"general_knowledge": {
		ID:          "general_knowledge",
		Label:       "一般知識",
		Icon:        "🌍",
		Description: "一般的な知識や常識を問う問題",
		PromptTips:  "知識を活用するよう促し、根拠を示すよう求めるプロンプトが有効です。",
		Color:       "#06B6D4", // cyan-500
	},
	"estimation": {
		ID:          "estimation",
		Label:       "推定・概算",
		Icon:        "📐",
		Description: "おおよその値を推定する問題",
		PromptTips:  "仮定を明確にし、段階的に推定するプロンプトが効果的です。",
		Color:       "#84CC16", // lime-500
	},
}

// GetTagByID はタグIDからタグメタデータを取得します。
// タグが見つからない場合は空のTag構造体を返します。
func GetTagByID(id string) (Tag, bool) {
	tag, exists := TagDefinitions[id]
	return tag, exists
}

// GetAllTags は全てのタグメタデータをスライスで返します。
func GetAllTags() []Tag {
	tags := make([]Tag, 0, len(TagDefinitions))
	for _, tag := range TagDefinitions {
		tags = append(tags, tag)
	}
	return tags
}
