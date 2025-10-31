/**
 * タグの定義情報
 */
export interface Tag {
  id: string; // タグの一意識別子
  label: string; // 表示用ラベル
  icon: string; // アイコン絵文字
  description: string; // タグの説明
  prompt_tips: string; // AI向けプロンプトヒント
  color: string; // 表示用カラーコード
}

/**
 * タグIDからタグメタデータへのマッピング
 */
export const TAG_DEFINITIONS: Record<string, Tag> = {
  calculation: {
    id: 'calculation',
    label: '計算問題',
    icon: '🔢',
    description: '数値計算や算術演算が必要な問題',
    prompt_tips: '数値を正確に計算させるため、段階的な計算プロセスを促すプロンプトが有効です。',
    color: '#3B82F6', // blue-500
  },
  character_counting: {
    id: 'character_counting',
    label: '文字数カウント',
    icon: '📊',
    description: '文字列の長さや特定文字の出現回数を数える問題',
    prompt_tips: '文字を一つずつ数えるように指示し、確認を促すプロンプトが効果的です。',
    color: '#10B981', // green-500
  },
  text_analysis: {
    id: 'text_analysis',
    label: 'テキスト解析',
    icon: '📝',
    description: '文章の構造や内容を分析する問題',
    prompt_tips: 'テキストを丁寧に読み解くよう促し、分析の観点を明示するプロンプトが有効です。',
    color: '#8B5CF6', // violet-500
  },
  text_problem: {
    id: 'text_problem',
    label: '文章題',
    icon: '📖',
    description: '文章から情報を読み取り、問題を解く必要がある問題',
    prompt_tips: '問題文を正確に理解させ、条件を整理してから解かせるプロンプトが効果的です。',
    color: '#F59E0B', // amber-500
  },
  pattern_recognition: {
    id: 'pattern_recognition',
    label: 'パターン認識',
    icon: '🔍',
    description: '規則性やパターンを見つけ出す問題',
    prompt_tips: '観察と仮説検証を促し、複数の例から規則を導き出すプロンプトが有効です。',
    color: '#EC4899', // pink-500
  },
  logic_puzzle: {
    id: 'logic_puzzle',
    label: '論理パズル',
    icon: '🧩',
    description: '論理的思考や推論が必要なパズル問題',
    prompt_tips: '論理的に段階を踏んで考えるよう促し、矛盾を確認するプロンプトが効果的です。',
    color: '#EF4444', // red-500
  },
  general_knowledge: {
    id: 'general_knowledge',
    label: '一般知識',
    icon: '🌍',
    description: '一般的な知識や常識を問う問題',
    prompt_tips: '知識を活用するよう促し、根拠を示すよう求めるプロンプトが有効です。',
    color: '#06B6D4', // cyan-500
  },
  estimation: {
    id: 'estimation',
    label: '推定・概算',
    icon: '📐',
    description: 'おおよその値を推定する問題',
    prompt_tips: '仮定を明確にし、段階的に推定するプロンプトが効果的です。',
    color: '#84CC16', // lime-500
  },
};

/**
 * タグIDからタグメタデータを取得
 */
export function getTagById(id: string): Tag | undefined {
  return TAG_DEFINITIONS[id];
}

/**
 * 全てのタグメタデータを取得
 */
export function getAllTags(): Tag[] {
  return Object.values(TAG_DEFINITIONS);
}
