"use client";

import { useEffect, useState } from "react";

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
      .then((data) => setHistory(data));
  }, []);

  return (
    <main className="max-w-4xl mx-auto p-8">

      <h1 className="text-3xl font-bold mb-8">
        History
      </h1>

      <button
        className="bg-red-600 text-white px-4 py-2 rounded mb-6"
        onClick={async () => {
          await fetch("http://localhost:8080/llm/history", {
            method: "DELETE",
          });

          setHistory([]);
        }}
      >
        Clear History
      </button>

      {history.length === 0 && (
        <p>No history yet.</p>
      )}

      {history.map((item, index) => (
        <div
          key={index}
          className="border rounded-xl p-5 mb-5"
        >
          <h2 className="font-bold">
            Prompt
          </h2>

          <p className="mb-4">
            {item.prompt}
          </p>

          <h2 className="font-bold">
            Response
          </h2>

          <p className="mb-4 whitespace-pre-wrap">
            {item.response}
          </p>

          <h2 className="font-bold text-blue-600">
            Score: {item.score}/100
          </h2>
        </div>
      ))}
    </main>
  );
}