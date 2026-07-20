"use client";

interface MetricsCardProps {
  inferenceTime: number | null;
}

export function MetricsCard({ inferenceTime }: MetricsCardProps) {
  if (inferenceTime === null) {
    return null;
  }

  return (
    <div className="mt-8">
      <h3 className="text-lg font-semibold text-black dark:text-white mb-3">
        Model Metrics
      </h3>
      <div className="p-4 border border-gray-300 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900">
        <p className="text-gray-800 dark:text-gray-100">
          Inference Time: {inferenceTime.toFixed(2)} seconds
        </p>
      </div>
    </div>
  );
}
