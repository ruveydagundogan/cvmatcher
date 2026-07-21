"use client";

import { useState, useCallback } from "react";

import { useWebLLM } from "@/hooks/useWebLLM";
import { useHistory } from "@/hooks/useHistory";
import { PromptInput } from "@/components/PromptInput";
import { AskButton } from "@/components/AskButton";
import { StatusCard } from "@/components/StatusCard";
import { ResponseCard } from "@/components/ResponseCard";
import { MetricsCard } from "@/components/MetricsCard";
import { DecisionScoreCard } from "@/components/DecisionScoreCard";

export default function DashboardPage() {
  const [prompt, setPrompt] = useState("");
  const { addItem } = useHistory();

  const {
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
    handleAskAI: llmHandleAskAI,
  } = useWebLLM();

  const handleAskAI = useCallback(async () => {
    if (!prompt.trim()) return;

    const currentPrompt = prompt;
    setPrompt("");

    const result = await llmHandleAskAI(currentPrompt);

    if (result && result.score > 0) {
      addItem({
        prompt: result.prompt,
        response: result.response,
        score: result.score,
        model: "Gemma-2B",
        word_count: result.wordCount,
        char_count: result.characterCount,
        inference_time: result.inferenceTime,
      });
    }
  }, [prompt, llmHandleAskAI, addItem]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 h-full">
      {/* Left - Chat */}
      <div className="flex flex-col">
        <div className="flex-1 bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 flex flex-col">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center animate-pulse-glow">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
              </svg>
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-white">AI Assistant</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">Gemma-2B • Browser LLM</p>
            </div>
          </div>

          <div className="flex-1 min-h-0">
            <ResponseCard responseText={responseText} isLoading={isInferenceRunning} />
          </div>

          <div className="mt-4">
            <StatusCard
              loadStatus={loadStatus}
              initProgress={initProgress}
              progressText={progressText}
              loadError={loadError}
            />
          </div>

          <div className="mt-4 space-y-3">
            <PromptInput value={prompt} onChange={setPrompt} />
            <AskButton
              disabled={loadStatus !== "loaded" || isInferenceRunning || !prompt.trim()}
              isLoading={isInferenceRunning}
              onClick={handleAskAI}
            />
          </div>
        </div>
      </div>

      {/* Right - Metrics & Score */}
      <div className="flex flex-col gap-6">
        <DecisionScoreCard score={score} />
        <MetricsCard
          inferenceTime={inferenceTime}
          wordCount={wordCount}
          characterCount={characterCount}
        />
      </div>
    </div>
  );
}
