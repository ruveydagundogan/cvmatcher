"use client";

import { useState } from "react";
import { useWebLLM } from "@/hooks/useWebLLM";
import { PromptInput } from "@/components/PromptInput";
import { AskButton } from "@/components/AskButton";
import { StatusCard } from "@/components/StatusCard";
import { ResponseCard } from "@/components/ResponseCard";
import { MetricsCard } from "@/components/MetricsCard";
import { DecisionScoreCard } from "@/components/DecisionScoreCard";

export default function Home() {
  const [prompt, setPrompt] = useState("");
  const {
    loadStatus,
    initProgress,
    progressText,
    loadError,
    isInferenceRunning,
    responseText,
    inferenceTime,
    score,
    handleAskAI,
  } = useWebLLM();

  const onAskAI = async () => {
    await handleAskAI(prompt);
    setPrompt("");
  };

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

        <PromptInput value={prompt} onChange={setPrompt} />

        {/* Ask AI Button */}
        <AskButton
          disabled={loadStatus !== "loaded" || isInferenceRunning}
          isLoading={isInferenceRunning}
          onClick={onAskAI}
        />

        <StatusCard
          loadStatus={loadStatus}
          initProgress={initProgress}
          progressText={progressText}
          loadError={loadError}
        />

        <ResponseCard responseText={responseText} isLoading={isInferenceRunning} />

        <MetricsCard inferenceTime={inferenceTime} />

        <DecisionScoreCard score={score} />
      </main>
    </div>
  );
}
