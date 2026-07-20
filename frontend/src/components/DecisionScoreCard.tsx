"use client";

interface DecisionScoreCardProps {
  score: number | null;
}

export function DecisionScoreCard({ score }: DecisionScoreCardProps) {
  if (score === null) {
    return null;
  }

  return (
    <div className="mt-4">
      <h4 className="text-sm font-semibold text-black dark:text-white mb-2">
        Decision Score
      </h4>
      <div className="p-3 border border-gray-300 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900">
        <p className="text-lg font-bold text-blue-600 dark:text-blue-400">
          {score} / 100
        </p>
      </div>
    </div>
  );
}
