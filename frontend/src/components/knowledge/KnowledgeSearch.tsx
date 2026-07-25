"use client";

import { useState } from "react";

interface KnowledgeEntry {
  id: string;
  title: string;
  content: string;
  tags: string[];
  category: string;
  source: string;
  created_at: string;
}

interface SearchResult {
  entry: KnowledgeEntry;
  score: number;
}

export default function KnowledgeSearch() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [searching, setSearching] = useState(false);

  const handleSearch = async () => {
    if (!query.trim()) return;
    setSearching(true);
    try {
      const mod = await import("@/lib/api");
      const data = await mod.api.get(`/api/v1/knowledge/search?q=${encodeURIComponent(query)}&limit=10`);
      setResults(data || []);
    } catch (e) {
      console.error("search error", e);
    } finally {
      setSearching(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex gap-3">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          placeholder="Search knowledge base..."
          className="flex-1 px-4 py-2.5 rounded-xl bg-white dark:bg-slate-900/50 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500/50"
        />
        <button
          onClick={handleSearch}
          disabled={searching}
          className="px-6 py-2.5 rounded-xl bg-purple-600 text-white font-medium hover:bg-purple-700 disabled:opacity-50"
        >
          {searching ? "..." : "Search"}
        </button>
      </div>

      {results.length > 0 && (
        <div className="space-y-3">
          {results.map((r) => (
            <div key={r.entry.id} className="bg-white dark:bg-slate-900/50 rounded-xl border border-gray-200 dark:border-white/10 p-4">
              <div className="flex items-start justify-between mb-2">
                <h3 className="font-semibold text-gray-900 dark:text-white">{r.entry.title}</h3>
                {r.entry.category && (
                  <span className="text-xs px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-400">{r.entry.category}</span>
                )}
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">{r.entry.content}</p>
              {r.entry.tags && r.entry.tags.length > 0 && (
                <div className="flex flex-wrap gap-1.5 mt-2">
                  {r.entry.tags.map((tag, i) => (
                    <span key={i} className="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-500">{tag}</span>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {!searching && query && results.length === 0 && (
        <p className="text-gray-500 text-center py-8">No results found</p>
      )}
    </div>
  );
}
