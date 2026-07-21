"use client";

interface MetricsCardProps {
  inferenceTime: number | null;
  wordCount: number;
  characterCount: number;
}

export function MetricsCard({
  inferenceTime,
  wordCount,
  characterCount,
}: MetricsCardProps) {
  if (inferenceTime === null) {
    return null;
  }

  return (
    <div className="mt-8">
      <h3 className="text-lg font-semibold text-black dark:text-white mb-3">
        Model Metrics
      </h3>

      <div className="grid grid-cols-2 gap-4">
        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Model</p>
          <p className="text-xl font-semibold">Gemma-2B</p>
        </div>

        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Inference Time</p>
          <p className="text-xl font-semibold">
            {inferenceTime.toFixed(2)} sec
          </p>
        </div>

        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Status</p>
          <p className="text-green-600 font-semibold">
            Loaded
          </p>
        </div>

        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Backend</p>
          <p className="text-blue-600 font-semibold">
            Connected
          </p>
        </div>

        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Word Count</p>
          <p className="text-xl font-semibold">
            {wordCount}
          </p>
        </div>

        <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-4">
          <p className="text-sm text-gray-500">Character Count</p>
          <p className="text-xl font-semibold">
            {characterCount}
          </p>
        </div>
      </div>
    </div>
  );
}