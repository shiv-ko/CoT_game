import React from 'react';

import styles from './LevelFilter.module.css';
import { renderStars } from './utils';

interface LevelFilterProps {
  selectedLevel: number | 'all';
  uniqueLevels: number[];
  onLevelChange: (level: number | 'all') => void;
}

/**
 * 難易度フィルタコンポーネント
 */
const LevelFilter: React.FC<LevelFilterProps> = ({
  selectedLevel,
  uniqueLevels,
  onLevelChange,
}) => {
  return (
    <div className={styles.filterContainer}>
      <label htmlFor="level-filter" className={styles.label}>
        難易度フィルタ:
      </label>
      <select
        id="level-filter"
        value={selectedLevel}
        onChange={(e) => onLevelChange(e.target.value === 'all' ? 'all' : Number(e.target.value))}
        className={styles.select}
      >
        <option value="all">すべて</option>
        {uniqueLevels.map((level) => (
          <option key={level} value={level}>
            レベル {level} ({renderStars(level)})
          </option>
        ))}
      </select>
      <div className={styles.info}>
        <span className={styles.badge}>
          {selectedLevel === 'all' ? '全問題' : `レベル ${selectedLevel}`}
        </span>
      </div>
    </div>
  );
};

export default LevelFilter;
