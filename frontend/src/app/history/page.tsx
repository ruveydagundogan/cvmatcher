"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface HistoryItem {
  prompt: string;
  response: string;
  score: number;
}

export default function HistoryPage() {
  const [history, setHistory] = useState<HistoryItem[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/llm/history")
      .then((res) => res.json())
      .then((data) => {
        setHistory(data.history ?? data);
      });
  }, []);

  const clearHistory = async () => {
    await fetch("http://localhost:8080/llm/history", {
      method: "DELETE",
    });

    setHistory([]);
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) {
      return "bg-green-100 text-green-700";
    }

    if (score >= 60) {
      return "bg-yellow-100 text-yellow-700";
    }

    return "bg-red-100 text-red-700";
  };

  return (
    <main className="max-w-6xl mx-auto px-8 py-10">

      {/* Header */}
      <div className="flex items-center justify-between mb-10">

        <h1 className="text-5xl font-bold">
          History
        </h1>

        <div className="flex items-center gap-5">

          <Link
            href="/dashboard"
            className="text-blue-600 hover:underline text-xl"
          >
            Dashboard
          </Link>

          <button
            onClick={clearHistory}
            className="bg-red-600 hover:bg-red-700 text-white px-6 py-3 rounded-xl font-semibold transition"
          >
            Clear History
          </button>

        </div>

      </div>

      {/* Empty State */}
      {history.length === 0 && (
        <div className="mt-20 flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-300 p-20 text-center">

          <div className="text-6xl mb-6">
            📝
          </div>

          <h2 className="text-3xl font-semibold mb-3">
            No history yet
          </h2>

          <p className="text-gray-500 max-w-lg">
            Ask the AI your first question from the dashboard.
            Every response and its decision score will appear here.
          </p>

        </div>
      )}

      {/* History List */}
      <div className="space-y-8">

        {history.map((item, index) => (

          <div
            key={index}
            className="border rounded-2xl p-8 shadow-sm"
          >

            <div className="flex items-center justify-between mb-8">

              <h2 className="text-3xl font-bold">
                Request #{index + 1}
              </h2>

              <div
                className={`px-5 py-2 rounded-full font-bold ${getScoreColor(item.score)}`}
              >
                {item.score}/100
              </div>

            </div>

            <div className="mb-8">

              <h3 className="font-semibold text-2xl mb-3">
                Prompt
              </h3>

              <div className="rounded-xl bg-blue-50 p-5">
                {item.prompt}
              </div>

            </div>

            <div>

              <h3 className="font-semibold text-2xl mb-3">
                Response
              </h3>

              <div className="rounded-xl bg-gray-100 p-5 whitespace-pre-wrap max-h-72 overflow-y-auto">
                {item.response}
              </div>

            </div>

          </div>

        ))}

      </div>

    </main>
  );
}