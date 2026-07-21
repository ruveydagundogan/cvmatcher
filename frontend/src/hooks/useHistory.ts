"use client";

import { useState, useEffect, useCallback } from "react";

export interface HistoryItem {
  id: string;
  prompt: string;
  response: string;
  score: number;
  model: string;
  word_count: number;
  char_count: number;
  inference_time: number;
  created_at: string;
}

function getStorageKey(): string {
  const userId = localStorage.getItem("userId") || localStorage.getItem("token") || "anonymous";
  return `llm_history_${userId}`;
}

export function useHistory() {
  const [history, setHistory] = useState<HistoryItem[]>([]);

  useEffect(() => {
    const stored = localStorage.getItem(getStorageKey());
    if (stored) {
      try {
        setHistory(JSON.parse(stored));
      } catch {
        setHistory([]);
      }
    }
  }, []);

  const addItem = useCallback((item: Omit<HistoryItem, "id" | "created_at">) => {
    const newItem: HistoryItem = {
      ...item,
      id: crypto.randomUUID(),
      created_at: new Date().toISOString(),
    };

    setHistory((prev) => {
      const updated = [newItem, ...prev];
      localStorage.setItem(getStorageKey(), JSON.stringify(updated));
      return updated;
    });
  }, []);

  const clearHistory = useCallback(() => {
    localStorage.removeItem(getStorageKey());
    setHistory([]);
  }, []);

  return { history, addItem, clearHistory };
}
