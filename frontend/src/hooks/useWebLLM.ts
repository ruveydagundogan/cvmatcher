"use client";

import { useEffect, useRef, useState } from "react";
import type { InitProgressReport, MLCEngine } from "@mlc-ai/web-llm";

export interface WebLLMState {
  loadStatus: "idle" | "loading" | "loaded" | "error";
  initProgress: number;
  progressText: string;
  loadError: string | null;
  isInferenceRunning: boolean;
  responseText: string;
  inferenceTime: number | null;
}

export interface WebLLMHandlers {
  handleAskAI: (prompt: string) => Promise<void>;
}

export interface UseWebLLMReturn extends WebLLMState, WebLLMHandlers {}

export function useWebLLM(): UseWebLLMReturn {
  const engineRef = useRef<MLCEngine | null>(null);
  const [loadStatus, setLoadStatus] = useState<"idle" | "loading" | "loaded" | "error">("idle");
  const [initProgress, setInitProgress] = useState<number>(0);
  const [progressText, setProgressText] = useState<string>("");
  const [loadError, setLoadError] = useState<string | null>(null);
  const [isInferenceRunning, setIsInferenceRunning] = useState(false);
  const [responseText, setResponseText] = useState<string>("");
  const [inferenceTime, setInferenceTime] = useState<number | null>(null);

  const handleAskAI = async (prompt: string) => {
    if (!engineRef.current) {
      return;
    }

    setIsInferenceRunning(true);
    setResponseText("");
    setInferenceTime(null);

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

      setResponseText(content || "No response returned.");
    } catch (error) {
      setResponseText(error instanceof Error ? error.message : String(error));
    } finally {
      setIsInferenceRunning(false);
    }
  };

  useEffect(() => {
    let mounted = true;

    async function initWebLLM() {
      setLoadStatus("loading");
      setInitProgress(0);
      setProgressText("Starting model initialization...");

      try {
        const { CreateMLCEngine } = await import("@mlc-ai/web-llm");

        // Model ID verified from the installed @mlc-ai/web-llm package source.
        const selectedModel = "gemma-2b-it-q4f16_1-MLC";

        const initProgressCallback = (report: InitProgressReport) => {
          if (!mounted) return;
          setInitProgress(report.progress);
          setProgressText(report.text);
        };

        const createdEngine = await CreateMLCEngine(selectedModel, {
          initProgressCallback,
        });

        if (mounted) {
          engineRef.current = createdEngine;
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
    handleAskAI,
  };
}
