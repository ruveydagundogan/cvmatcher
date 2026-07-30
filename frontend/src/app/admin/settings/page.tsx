"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

export default function SettingsPage() {
  const [maxTokens, setMaxTokens] = useState(2048);
  const [temperature, setTemperature] = useState(0.7);
  const [topP, setTopP] = useState(0.9);
  const [contextLength, setContextLength] = useState(4096);
  const [modelName, setModelName] = useState("qwen2.5:1.5b-instruct");
  const [loading, setLoading] = useState(true);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    api.get("/api/v1/admin/settings")
      .then((s) => {
        if (s) {
          setMaxTokens(s.max_tokens ?? 2048);
          setTemperature(s.temperature ?? 0.7);
          setTopP(s.top_p ?? 0.9);
          setContextLength(s.context_length ?? 4096);
          setModelName(s.model_name ?? "qwen2.5:1.5b-instruct");
        }
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const save = async () => {
    try {
      await api.put("/api/v1/admin/settings", {
        max_tokens: maxTokens,
        temperature,
        top_p: topP,
        context_length: contextLength,
        model_name: modelName,
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (e) {
      console.error(e);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">LLM Settings</h1>

      <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 max-w-xl space-y-5">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Model</label>
          <input type="text" value={modelName} onChange={(e) => setModelName(e.target.value)}
            className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Max Tokens: {maxTokens}
          </label>
          <input type="range" min="256" max="8192" step="256" value={maxTokens}
            onChange={(e) => setMaxTokens(Number(e.target.value))}
            className="w-full accent-purple-500" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Temperature: {temperature.toFixed(2)}
          </label>
          <input type="range" min="0" max="2" step="0.05" value={temperature}
            onChange={(e) => setTemperature(Number(e.target.value))}
            className="w-full accent-purple-500" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Top-P: {topP.toFixed(2)}
          </label>
          <input type="range" min="0" max="1" step="0.05" value={topP}
            onChange={(e) => setTopP(Number(e.target.value))}
            className="w-full accent-purple-500" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Context Length: {contextLength}
          </label>
          <input type="range" min="1024" max="32768" step="1024" value={contextLength}
            onChange={(e) => setContextLength(Number(e.target.value))}
            className="w-full accent-purple-500" />
        </div>

        <button onClick={save}
          className="px-6 py-2.5 rounded-xl bg-purple-600 text-white font-medium hover:bg-purple-700">
          {saved ? "Saved!" : "Save Settings"}
        </button>
      </div>
    </div>
  );
}
