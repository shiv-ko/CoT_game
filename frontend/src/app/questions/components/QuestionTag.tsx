import React, { useState } from 'react';

import styles from './QuestionTag.module.css';

import { getTagById, type Tag } from '@/types/tag';

interface QuestionTagProps {
  tagId: string;
}

/**
 * 問題のタグを表示するコンポーネント
 */
export function QuestionTag({ tagId }: QuestionTagProps) {
  const tag = getTagById(tagId);

  if (!tag) {
    return null;
  }

  return (
    <span className={styles.tag} style={{ backgroundColor: tag.color }} title={tag.description}>
      <span className={styles.icon}>{tag.icon}</span>
      <span className={styles.label}>{tag.label}</span>
    </span>
  );
}

interface QuestionTagListProps {
  tags?: string[];
}

/**
 * 複数のタグを表示するコンポーネント
 */
export function QuestionTagList({ tags }: QuestionTagListProps) {
  if (!tags || tags.length === 0) {
    return null;
  }

  return (
    <div className={styles.tagList}>
      {tags.map((tagId) => (
        <QuestionTag key={tagId} tagId={tagId} />
      ))}
    </div>
  );
}

interface PromptTipsProps {
  tags?: string[];
}

/**
 * タグに基づいたプロンプトヒントを表示するコンポーネント
 * トグルで開閉可能
 */
export function PromptTips({ tags }: PromptTipsProps) {
  const [isOpen, setIsOpen] = useState(false);

  if (!tags || tags.length === 0) {
    return null;
  }

  const tips = tags
    .map((tagId) => getTagById(tagId))
    .filter((tag): tag is Tag => tag !== undefined)
    .map((tag) => tag.prompt_tips);

  if (tips.length === 0) {
    return null;
  }

  return (
    <div className={styles.promptTips}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={styles.tipsToggle}
        aria-expanded={isOpen}
      >
        <span className={styles.toggleIcon}>{isOpen ? '▼' : '▶'}</span>
        <span className={styles.tipsTitle}>💡 プロンプトのヒント</span>
      </button>
      {isOpen && (
        <ul className={styles.tipsList}>
          {tips.map((tip, index) => (
            <li key={index} className={styles.tipItem}>
              {tip}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
