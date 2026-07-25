"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import KnowledgeSearch from "@/components/knowledge/KnowledgeSearch";

interface KnowledgeEntry {
  id: string;
  title: string;
  content: string;
  tags: string[];
  category: string;
  source: string;
  created_at: string;
}

export default function KnowledgePage() {
  const [entries, setEntries] = useState<KnowledgeEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [category, setCategory] = useState("");
  const [tagInput, setTagInput] = useState("");
  const [tags, setTags] = useState<string[]>([]);
  const [showForm, setShowForm] = useState(false);

  const load = async () => {
    try {
      const data = await api.get("/api/v1/knowledge?limit=50");
      setEntries(data.items || []);
      setTotal(data.total || 0);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const addTag = () => {
    const t = tagInput.trim();
    if (t && !tags.includes(t)) setTags([...tags, t]);
    setTagInput("");
  };

  const create = async () => {
    if (!title.trim() || !content.trim()) return;
    try {
      await api.post("/api/v1/knowledge", { title, content, category, tags });
      setTitle(""); setContent(""); setCategory(""); setTags([]);
      setShowForm(false);
      load();
    } catch (e) {
      console.error(e);
    }
  };

  const removeEntry = async (id: string) => {
    try {
      await api.delete(`/api/v1/knowledge/${id}`);
      load();
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Knowledge Base</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">DeepKwiki — search and manage your knowledge</p>
        </div>
        <button onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 rounded-xl bg-purple-600 text-white text-sm font-medium hover:bg-purple-700">
          {showForm ? "Cancel" : "Add Entry"}
        </button>
      </div>

      {showForm && (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6 space-y-4">
          <input type="text" value={title} onChange={(e) => setTitle(e.target.value)}
            placeholder="Title" className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <textarea value={content} onChange={(e) => setContent(e.target.value)} rows={5}
            placeholder="Content" className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <input type="text" value={category} onChange={(e) => setCategory(e.target.value)}
            placeholder="Category" className="w-full px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <div className="flex gap-2">
            <input type="text" value={tagInput} onChange={(e) => setTagInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && addTag()}
              placeholder="Add tag..." className="flex-1 px-4 py-2 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
            <button onClick={addTag} className="px-3 py-2 rounded-xl bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-sm">Add</button>
          </div>
          {tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5">
              {tags.map((t, i) => (
                <span key={i} className="text-xs px-2 py-1 rounded-full bg-purple-500/10 text-purple-400">
                  {t} <button onClick={() => setTags(tags.filter((_, j) => j !== i))} className="ml-1">&times;</button>
                </span>
              ))}
            </div>
          )}
          <button onClick={create}
            className="px-6 py-2.5 rounded-xl bg-purple-600 text-white font-medium hover:bg-purple-700">
            Save
          </button>
        </div>
      )}

      <div className="mb-6">
        <KnowledgeSearch />
      </div>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : entries.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No knowledge entries yet</div>
      ) : (
        <div className="grid gap-3">
          {entries.map((e) => (
            <div key={e.id} className="bg-white dark:bg-slate-900/50 rounded-xl border border-gray-200 dark:border-white/10 p-4">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-white">{e.title}</h3>
                  <p className="text-sm text-gray-500 mt-1 line-clamp-2">{e.content}</p>
                </div>
                <button onClick={() => removeEntry(e.id)}
                  className="text-xs text-red-400 hover:text-red-300 ml-4">Delete</button>
              </div>
              {(e.category || (e.tags && e.tags.length > 0)) && (
                <div className="flex items-center gap-2 mt-2">
                  {e.category && <span className="text-xs px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-400">{e.category}</span>}
                  {e.tags?.map((t, i) => (
                    <span key={i} className="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-500">{t}</span>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
