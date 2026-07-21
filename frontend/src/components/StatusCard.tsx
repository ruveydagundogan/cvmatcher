"use client";

interface StatusCardProps {
  loadStatus: "idle" | "loading" | "loaded" | "error";
  initProgress: number;
  progressText: string;
  loadError: string | null;
}

export function StatusCard({
  loadStatus,
  initProgress,
  progressText,
  loadError,
}: StatusCardProps) {
  if (loadStatus === "idle") {
    return null;
  }

  return (
    <div className="rounded-xl bg-white/5 border border-white/10 p-3">
      {loadStatus === "loading" && (
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-purple-400 font-medium">Loading model...</span>
            <span className="text-gray-400">{Math.round(initProgress * 100)}%</span>
          </div>
          <div className="w-full h-1.5 bg-white/5 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-blue-500 to-purple-600 rounded-full transition-all duration-300"
              style={{ width: `${initProgress * 100}%` }}
            />
          </div>
          {progressText && (
            <p className="text-xs text-gray-500 truncate">{progressText}</p>
          )}
        </div>
      )}

      {loadStatus === "loaded" && (
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-sm text-green-400 font-medium">Model ready</span>
        </div>
      )}

      {loadStatus === "error" && (
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-red-500" />
          <span className="text-sm text-red-400 font-medium">Error: {loadError}</span>
        </div>
      )}
    </div>
  );
}
