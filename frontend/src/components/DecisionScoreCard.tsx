"use client";

interface DecisionScoreCardProps {
  score: number | null;
}

export function DecisionScoreCard({ score }: DecisionScoreCardProps) {
  if (score === null) {
    return null;
  }

  const getScoreColor = () => {
    if (score >= 80) {
      return "text-green-600";
    }

    if (score >= 60) {
      return "text-yellow-600";
    }

    return "text-red-600";
  };

  const getBackgroundColor = () => {
    if (score >= 80) {
      return "bg-green-50 border-green-200";
    }

    if (score >= 60) {
      return "bg-yellow-50 border-yellow-200";
    }

    return "bg-red-50 border-red-200";
  };

  return (
    <div className="mt-8">
      <h3 className="text-lg font-semibold text-black dark:text-white mb-3">
        Decision Score
      </h3>

      <div
        className={`rounded-lg border p-6 ${getBackgroundColor()}`}
      >
        <p className={`text-4xl font-bold ${getScoreColor()}`}>
          {score} / 100
        </p>
      </div>
    </div>
  );
}