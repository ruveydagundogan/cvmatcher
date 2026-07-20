"use client";

interface ResponseCardProps {
  responseText: string;
  isLoading: boolean;
}

export function ResponseCard({ responseText, isLoading }: ResponseCardProps) {
  return (
    <div className="mt-12">
      <h2 className="text-xl font-semibold text-black dark:text-white mb-4">
        Response
      </h2>
      <div className="p-6 border border-gray-300 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900 min-h-24">
        <p className="text-gray-800 dark:text-gray-100 whitespace-pre-wrap">
          {responseText || (isLoading ? "Waiting for model response..." : "Response will appear here...")}
        </p>
      </div>
    </div>
  );
}
