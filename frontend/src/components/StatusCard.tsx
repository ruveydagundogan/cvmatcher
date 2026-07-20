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
  return (
    <div className="mt-4 text-sm">
      {loadStatus === "idle" && <p className="text-gray-500">Model not started.</p>}
      {loadStatus === "loading" && (
        <p className="text-blue-600">
          Loading model... {Math.round(initProgress * 100)}% {progressText}
        </p>
      )}
      {loadStatus === "loaded" && <p className="text-green-600">Model loaded (ready).</p>}
      {loadStatus === "error" && (
        <p className="text-red-600">Model load error: {loadError}</p>
      )}
    </div>
  );
}
