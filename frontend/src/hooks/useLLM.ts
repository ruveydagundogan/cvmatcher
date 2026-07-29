"use client";

import { useState, useCallback } from "react";
import { API_BASE_URL } from "@/lib/config";

export interface ScoreResult {
  prompt: string;
  response: string;
  score: number;
  inferenceTime: number;
  wordCount: number;
  characterCount: number;
}

export function useLLM() {
  const [isInferenceRunning, setIsInferenceRunning] = useState(false);
  const [responseText, setResponseText] = useState("");
  const [inferenceTime, setInferenceTime] = useState<number | null>(null);
  const [score, setScore] = useState<number | null>(null);
  const [wordCount, setWordCount] = useState(0);
  const [characterCount, setCharacterCount] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const handleAskAI = useCallback(async (prompt: string): Promise<ScoreResult | null> => {
    setIsInferenceRunning(true);
    setResponseText("");
    setInferenceTime(null);
    setScore(null);
    setWordCount(0);
    setCharacterCount(0);
    setError(null);

    try {
      const token = localStorage.getItem("token");
      const headers: Record<string, string> = { "Content-Type": "application/json" };
      if (token) headers["Authorization"] = `Bearer ${token}`;

      const startTime = performance.now();
      const chatRes = await fetch(`${API_BASE_URL}/api/v1/llm/chat`, {
        method: "POST",
        headers,
        body: JSON.stringify({ prompt, max_tokens: 256 }),
      });
      const endTime = performance.now();
      const elapsedTime = (endTime - startTime) / 1000;
      setInferenceTime(elapsedTime);

      if (!chatRes.ok) {
        const errBody = await chatRes.json().catch(() => ({}));
        const errMsg = errBody?.error || `Backend returned ${chatRes.status}`;
        setError(errMsg);
        setResponseText(errMsg);
        return null;
      }

      const chatData = await chatRes.json();
      const responseContent = chatData?.data?.response || chatData?.response || "";
      setResponseText(responseContent);

      const charCount = responseContent.length;
      setCharacterCount(charCount);

      const wordCnt = responseContent.trim().split(/\s+/).filter(Boolean).length;
      setWordCount(wordCnt);

      let scoreValue: number | null = null;
      try {
        const scoreRes = await fetch(`${API_BASE_URL}/api/v1/score`, {
          method: "POST",
          headers,
          body: JSON.stringify({ prompt, response: responseContent }),
        });

        if (scoreRes.ok) {
          const scoreData = await scoreRes.json();
          scoreValue = scoreData?.data?.score ?? scoreData?.score ?? null;
          setScore(scoreValue);
        }
      } catch {
        console.error("Score endpoint failed");
      }

      return {
        prompt,
        response: responseContent,
        score: scoreValue || 0,
        inferenceTime: elapsedTime,
        wordCount: wordCnt,
        characterCount: charCount,
      };
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      setError(msg);
      setResponseText(msg);
      return null;
    } finally {
      setIsInferenceRunning(false);
    }
  }, []);

  return {
    loadStatus: error ? "error" as const : "loaded" as const,
    loadError: error,
    isInferenceRunning,
    responseText,
    inferenceTime,
    score,
    wordCount,
    characterCount,
    handleAskAI,
  };
}
