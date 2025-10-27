import React from 'react';

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
    <div style={{ marginBottom: '1.5rem' }}>
      <label htmlFor="level-filter" style={{ marginRight: '0.5rem' }}>
        難易度フィルタ:
      </label>
      <select
        id="level-filter"
        value={selectedLevel}
        onChange={(e) => onLevelChange(e.target.value === 'all' ? 'all' : Number(e.target.value))}
        style={{ padding: '0.5rem' }}
      >
        <option value="all">すべて</option>
        {uniqueLevels.map((level) => (
          <option key={level} value={level}>
            レベル {level} ({renderStars(level)})
          </option>
        ))}
      </select>
    </div>
  );
};

export default LevelFilter;
