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

  const metrics = [
    {
      label: "Model",
      value: "Gemma-2B",
      icon: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      ),
      color: "from-blue-500 to-cyan-500",
    },
    {
      label: "Inference",
      value: `${inferenceTime.toFixed(2)}s`,
      icon: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
      color: "from-purple-500 to-pink-500",
    },
    {
      label: "Words",
      value: wordCount.toString(),
      icon: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
        </svg>
      ),
      color: "from-green-500 to-emerald-500",
    },
    {
      label: "Characters",
      value: characterCount.toString(),
      icon: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      ),
      color: "from-orange-500 to-red-500",
    },
  ];

  return (
    <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
          <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white">Metrics</h3>
          <p className="text-xs text-gray-500">Inference details</p>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        {metrics.map((metric) => (
          <div
            key={metric.label}
            className="bg-white/5 dark:bg-white/5 rounded-xl border border-gray-100 dark:border-white/5 p-4 hover:border-purple-500/30 transition-all duration-200"
          >
            <div className={`w-8 h-8 rounded-lg bg-gradient-to-br ${metric.color} flex items-center justify-center mb-3`}>
              <span className="text-white">{metric.icon}</span>
            </div>
            <p className="text-xs text-gray-500 mb-1">{metric.label}</p>
            <p className="text-lg font-bold text-gray-900 dark:text-white">{metric.value}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
