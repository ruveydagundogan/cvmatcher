"use client";

interface AskButtonProps {
  disabled: boolean;
  isLoading: boolean;
  onClick: () => void;
}

export function AskButton({ disabled, isLoading, onClick }: AskButtonProps) {
  return (
    <button
      type="button"
      disabled={disabled}
      onClick={onClick}
      className={`w-full py-3 rounded-lg font-semibold text-white ${
        disabled
          ? "bg-gray-400 dark:bg-gray-700 cursor-not-allowed opacity-60"
          : "bg-blue-600 hover:bg-blue-700"
      }`}
    >
      {isLoading ? "Thinking..." : "Ask AI"}
    </button>
  );
}
