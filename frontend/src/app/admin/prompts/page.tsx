"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

interface Prompt {
  id: string;
  name: string;
  content: string;
  active: boolean;
}

export default function PromptsPage() {
  const [prompts, setPrompts] = useState<Prompt[]>([]);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState("");
  const [content, setContent] = useState("");

  const load = async () => {
    try {
      setPrompts(await api.get("/api/v1/admin/prompts"));
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const create = async () => {
    if (!name.trim() || !content.trim()) return;
    try {
      await api.post("/api/v1/admin/prompts", { name, content });
      setName(""); setContent("");
      load();
    } catch (e) {
      console.error(e);
    }
  };

  const activate = async (id: string) => {
    try {
      await api.post(`/api/v1/admin/prompts/${id}/activate`, {});
      load();
    } catch (e) {
      console.error(e);
    }
  };

  const remove = async (id: string) => {
    try {
      await api.delete(`/api/v1/admin/prompts/${id}`);
      load();
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">System Prompts</h1>

      <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">New System Prompt</h2>
        <input type="text" value={name} onChange={(e) => setName(e.target.value)}
          placeholder="Prompt name" className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white mb-3" />
        <textarea value={content} onChange={(e) => setContent(e.target.value)} rows={6}
          placeholder="System prompt content..." className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white mb-3 font-mono text-sm" />
        <button onClick={create}
          className="px-6 py-2.5 rounded-xl bg-purple-600 text-white font-medium hover:bg-purple-700">
          Create Prompt
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : prompts.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No system prompts yet</div>
      ) : (
        <div className="space-y-3">
          {prompts.map((p) => (
            <div key={p.id} className="bg-white dark:bg-slate-900/50 rounded-xl border border-gray-200 dark:border-white/10 p-4">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <h3 className="font-semibold text-gray-900 dark:text-white">{p.name}</h3>
                  {p.active && (
                    <span className="text-xs px-2 py-0.5 rounded-full bg-green-500/10 text-green-400">Active</span>
                  )}
                </div>
                <div className="flex gap-2">
                  {!p.active && (
                    <button onClick={() => activate(p.id)}
                      className="text-xs px-3 py-1 rounded-lg bg-blue-500/10 text-blue-400 hover:bg-blue-500/20">
                      Activate
                    </button>
                  )}
                  <button onClick={() => remove(p.id)}
                    className="text-xs text-red-400 hover:text-red-300">Delete</button>
                </div>
              </div>
              <pre className="text-xs text-gray-600 dark:text-gray-400 whitespace-pre-wrap line-clamp-3 font-mono">{p.content}</pre>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
