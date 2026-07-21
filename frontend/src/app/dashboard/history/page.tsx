"use client";

import { useHistory } from "@/hooks/useHistory";

export default function HistoryPage() {
  const { history, clearHistory } = useHistory();

  const getScoreColor = (score: number) => {
    if (score >= 80) {
      return "from-green-500 to-emerald-600";
    }
    if (score >= 60) {
      return "from-yellow-500 to-orange-600";
    }
    return "from-red-500 to-rose-600";
  };

  const getScoreBg = (score: number) => {
    if (score >= 80) {
      return "bg-green-500/10 border-green-500/20";
    }
    if (score >= 60) {
      return "bg-yellow-500/10 border-yellow-500/20";
    }
    return "bg-red-500/10 border-red-500/20";
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">History</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            {history.length} scored {history.length === 1 ? "request" : "requests"}
          </p>
        </div>

        {history.length > 0 && (
          <button
            onClick={clearHistory}
            className="px-5 py-2.5 rounded-xl bg-red-500/10 text-red-500 hover:bg-red-500/20 border border-red-500/20 font-medium transition-all duration-200"
          >
            Clear All
          </button>
        )}
      </div>

      {history.length === 0 && (
        <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 dark:border-white/10 p-20 text-center bg-white/5 dark:bg-white/5">
          <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center mb-6">
            <svg className="w-10 h-10 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">No history yet</h2>
          <p className="text-gray-500 dark:text-gray-400 max-w-md">
            Start a conversation in the Chat to see your scored requests here.
          </p>
        </div>
      )}

      <div className="space-y-4">
        {history.map((item, index) => (
          <div
            key={item.id || index}
            className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 hover:border-purple-500/30 transition-all duration-200"
          >
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                  <span className="text-white text-sm font-bold">#{index + 1}</span>
                </div>
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {item.model || "Gemma-2B"}
                  </p>
                  <p className="text-xs text-gray-400 dark:text-gray-500">
                    {item.created_at ? new Date(item.created_at).toLocaleString() : ""}
                  </p>
                </div>
              </div>

              <div className={`px-4 py-2 rounded-xl border bg-gradient-to-r ${getScoreColor(item.score)} ${getScoreBg(item.score)}`}>
                <span className="text-white font-bold">{item.score}</span>
                <span className="text-white/60 text-sm">/100</span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-xs font-medium text-purple-400 mb-2 uppercase tracking-wider">Prompt</p>
                <div className="rounded-xl bg-blue-500/5 border border-blue-500/10 p-4 text-sm text-gray-700 dark:text-gray-300 max-h-32 overflow-y-auto">
                  {item.prompt}
                </div>
              </div>

              <div>
                <p className="text-xs font-medium text-purple-400 mb-2 uppercase tracking-wider">Response</p>
                <div className="rounded-xl bg-purple-500/5 border border-purple-500/10 p-4 text-sm text-gray-700 dark:text-gray-300 max-h-32 overflow-y-auto whitespace-pre-wrap">
                  {item.response}
                </div>
              </div>
            </div>

            <div className="flex gap-4 mt-4 pt-4 border-t border-gray-100 dark:border-white/5">
              <span className="text-xs text-gray-500">
                <span className="font-medium text-gray-700 dark:text-gray-300">{item.word_count}</span> words
              </span>
              <span className="text-xs text-gray-500">
                <span className="font-medium text-gray-700 dark:text-gray-300">{item.char_count}</span> chars
              </span>
              {item.inference_time > 0 && (
                <span className="text-xs text-gray-500">
                  <span className="font-medium text-gray-700 dark:text-gray-300">{item.inference_time.toFixed(2)}s</span> inference
                </span>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
