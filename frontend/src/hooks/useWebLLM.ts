"use client";

import { useEffect, useRef, useState } from "react";
import type { InitProgressReport, MLCEngine } from "@mlc-ai/web-llm";

const SELECTED_MODEL = "gemma-2b-it-q4f16_1-MLC";
const MODEL_LOADED_KEY = "webllm_engine_ready";

let sharedEngine: MLCEngine | null = null;
let enginePromise: Promise<MLCEngine> | null = null;

async function getOrCreateEngine(
  onProgress?: (report: InitProgressReport) => void
): Promise<MLCEngine> {
  if (sharedEngine) return sharedEngine;

  if (enginePromise) return enginePromise;

  enginePromise = (async () => {
    const { CreateMLCEngine } = await import("@mlc-ai/web-llm");
    const engine = await CreateMLCEngine(SELECTED_MODEL, {
      initProgressCallback: onProgress,
    });
    sharedEngine = engine;
    try {
      sessionStorage.setItem(MODEL_LOADED_KEY, "true");
    } catch {}
    return engine;
  })();

  return enginePromise;
}

export interface WebLLMState {
  loadStatus: "idle" | "loading" | "loaded" | "error";
  initProgress: number;
  progressText: string;
  loadError: string | null;
  isInferenceRunning: boolean;
  responseText: string;
  inferenceTime: number | null;
  score: number | null;
  wordCount: number;
  characterCount: number;
}

export interface ScoreResult {
  prompt: string;
  response: string;
  score: number;
  inferenceTime: number;
  wordCount: number;
  characterCount: number;
}

export interface WebLLMHandlers {
  handleAskAI: (prompt: string) => Promise<ScoreResult | null>;
}

export interface UseWebLLMReturn extends WebLLMState, WebLLMHandlers { }

export function useWebLLM(): UseWebLLMReturn {
  const engineRef = useRef<MLCEngine | null>(null);
  const [loadStatus, setLoadStatus] = useState<"idle" | "loading" | "loaded" | "error">(() => {
    if (sharedEngine) return "loaded";
    try {
      if (sessionStorage.getItem(MODEL_LOADED_KEY) === "true") return "loading";
    } catch {}
    return "idle";
  });
  const [initProgress, setInitProgress] = useState<number>(0);
  const [progressText, setProgressText] = useState<string>("");
  const [loadError, setLoadError] = useState<string | null>(null);
  const [isInferenceRunning, setIsInferenceRunning] = useState(false);
  const [responseText, setResponseText] = useState<string>("");
  const [inferenceTime, setInferenceTime] = useState<number | null>(null);
  const [score, setScore] = useState<number | null>(null);
  const [wordCount, setWordCount] = useState(0);
  const [characterCount, setCharacterCount] = useState(0);

  const handleAskAI = async (prompt: string): Promise<ScoreResult | null> => {
    if (!engineRef.current) {
      return null;
    }

    setIsInferenceRunning(true);
    setResponseText("");
    setInferenceTime(null);
    setScore(null);
    setWordCount(0);
    setCharacterCount(0);

    try {
      const startTime = performance.now();
      const completion = await engineRef.current.chat.completions.create({
        messages: [{ role: "user", content: prompt }],
        max_tokens: 256,
      });
      const endTime = performance.now();
      const elapsedTime = (endTime - startTime) / 1000;
      setInferenceTime(elapsedTime);

      const assistantMessage = completion.choices?.[0]?.message;
      const rawContent: unknown = assistantMessage?.content;
      const content = typeof rawContent === "string"
        ? rawContent
        : Array.isArray(rawContent)
          ? rawContent.map((part: any) => part.type === "text" ? part.text : "").join("")
          : "";

      const responseContent = content || "No response returned.";

      setResponseText(responseContent);

      const charCount = responseContent.length;
      setCharacterCount(charCount);

      const wordCnt = responseContent
        .trim()
        .split(/\s+/)
        .filter(Boolean).length;
      setWordCount(wordCnt);

      let scoreValue: number | null = null;

      try {
        const token = localStorage.getItem("token");
        const headers: Record<string, string> = {
          "Content-Type": "application/json",
        };
        if (token) {
          headers["Authorization"] = `Bearer ${token}`;
        }

        const scoreResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/v1/score`, {
          method: "POST",
          headers,
          body: JSON.stringify({
            prompt: prompt,
            response: responseContent,
          }),
        });

        if (scoreResponse.ok) {
          const scoreData = await scoreResponse.json();
          if (scoreData.success && scoreData.data?.score !== undefined) {
            scoreValue = scoreData.data.score;
          } else {
            scoreValue = scoreData.score ?? null;
          }
          setScore(scoreValue);
        } else {
          console.error("Failed to get score from backend:", scoreResponse.status);
        }
      } catch (scoreError) {
        console.error("Error calling /score endpoint:", scoreError);
      }

      return {
        prompt,
        response: responseContent,
        score: scoreValue || 0,
        inferenceTime: elapsedTime,
        wordCount: wordCnt,
        characterCount: charCount,
      };
    } catch (error) {
      setResponseText(error instanceof Error ? error.message : String(error));
      return null;
    } finally {
      setIsInferenceRunning(false);
    }
  };

  useEffect(() => {
    let mounted = true;

    async function initWebLLM() {
      if (sharedEngine) {
        engineRef.current = sharedEngine;
        setLoadStatus("loaded");
        return;
      }

      setLoadStatus("loading");
      setInitProgress(0);
      setProgressText("Starting model initialization...");

      try {
        const engine = await getOrCreateEngine((report) => {
          if (!mounted) return;
          setInitProgress(report.progress);
          setProgressText(report.text);
        });

        if (mounted) {
          engineRef.current = engine;
          setLoadStatus("loaded");
        }
      } catch (error) {
        if (mounted) {
          setLoadError(error instanceof Error ? error.message : String(error));
          setLoadStatus("error");
        }
      }
    }

    initWebLLM();

    return () => {
      mounted = false;
    };
  }, []);

  return {
    loadStatus,
    initProgress,
    progressText,
    loadError,
    isInferenceRunning,
    responseText,
    inferenceTime,
    score,
    wordCount,
    characterCount,
    handleAskAI,
  };
}
