/**
 * レベルを星で表示
 * @param level 問題のレベル (1-5)
 * @returns 星の文字列 (例: ★★★☆☆)
 */
export const renderStars = (level: number): string => {
  const maxStars = 5;
  const filledStars = '★'.repeat(level);
  const emptyStars = '☆'.repeat(maxStars - level);
  return filledStars + emptyStars;
};
