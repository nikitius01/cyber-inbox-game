const labels = {
  easy: 'Easy',
  medium: 'Medium',
  hard: 'Hard',
  nightmare: 'Nightmare',
  AI: 'AI',
};

export default function DifficultyBadge({ value }) {
  return <span className={`difficulty difficulty-${value}`}>{labels[value] || value}</span>;
}

