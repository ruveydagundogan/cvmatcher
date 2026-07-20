"use client";

import { useEffect, useRef, useState } from "react";
import type { InitProgressReport, MLCEngine } from "@mlc-ai/web-llm";

export default function Home() {
  const [prompt, setPrompt] = useState("");
  const engineRef = useRef<MLCEngine | null>(null);
  const [loadStatus, setLoadStatus] = useState<"idle" | "loading" | "loaded" | "error">("idle");
  const [initProgress, setInitProgress] = useState<number>(0);
  const [progressText, setProgressText] = useState<string>("");
  const [loadError, setLoadError] = useState<string | null>(null);
  const [isInferenceRunning, setIsInferenceRunning] = useState(false);
  const [responseText, setResponseText] = useState<string>("");

  const handleAskAI = async () => {
    if (!engineRef.current) {
      return;
    }

    setIsInferenceRunning(true);
    setResponseText("");

    try {
      const completion = await engineRef.current.chat.completions.create({
        messages: [{ role: "user", content: prompt }],
        max_tokens: 256,
      });

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

  return (
    <div className="min-h-screen bg-white dark:bg-black">
      <main className="max-w-2xl mx-auto px-6 py-12">
        {/* Header */}
        <div className="mb-12">
          <h1 className="text-4xl font-bold text-black dark:text-white mb-2">
            LLM Decision Score
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Ask the AI a question and see how it decides.
          </p>
        </div>

        {/* Prompt Input */}
        <div className="mb-8">
          <label className="block text-sm font-semibold text-black dark:text-white mb-3">
            Your Question
          </label>
          <textarea
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Type your question here..."
            className="w-full h-24 px-4 py-3 border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 text-black dark:text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Ask AI Button */}
        <button
          type="button"
          disabled={loadStatus !== "loaded" || isInferenceRunning}
          onClick={handleAskAI}
          className={`w-full py-3 rounded-lg font-semibold text-white ${loadStatus === "loaded" && !isInferenceRunning ? "bg-blue-600 hover:bg-blue-700" : "bg-gray-400 dark:bg-gray-700 cursor-not-allowed opacity-60"}`}
        >
          {isInferenceRunning ? "Thinking..." : "Ask AI"}
        </button>

        {/* Model Loading Status */}
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

        {/* Response Section */}
        <div className="mt-12">
          <h2 className="text-xl font-semibold text-black dark:text-white mb-4">
            Response
          </h2>
          <div className="p-6 border border-gray-300 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900 min-h-24">
            <p className="text-gray-800 dark:text-gray-100 whitespace-pre-wrap">
              {responseText || (isInferenceRunning ? "Waiting for model response..." : "Response will appear here...")}
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
