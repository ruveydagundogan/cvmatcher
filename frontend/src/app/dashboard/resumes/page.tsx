"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

interface CV {
  id: string;
  title: string;
  status: string;
  parsed_skills: string[];
  parsed_summary: string;
  created_at: string;
}

export default function CVListPage() {
  const [cvs, setCVs] = useState<CV[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showCreate, setShowCreate] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [creating, setCreating] = useState(false);

  const fetchCVs = () => {
    setLoading(true);
    api.get("/api/v1/cvs?page=1&limit=50")
      .then((data) => { setCVs(data.items || []); setTotal(data.total || 0); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchCVs(); }, []);

  const handleCreate = async () => {
    if (!title.trim() || !content.trim()) return;
    setCreating(true);
    try {
      await api.post("/api/v1/cvs", { title, content });
      setTitle("");
      setContent("");
      setShowCreate(false);
      fetchCVs();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this CV?")) return;
    try {
      await api.delete(`/api/v1/cvs/${id}`);
      fetchCVs();
    } catch (e: any) {
      setError(e.message);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">CVs</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">{total} CV{cvs.length !== 1 ? "s" : ""}</p>
        </div>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-all duration-200"
        >
          {showCreate ? "Cancel" : "New CV"}
        </button>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      {showCreate && (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Create CV</h2>
          <div className="space-y-4">
            <input
              type="text"
              placeholder="CV Title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
            <textarea
              placeholder="Paste CV content here..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={10}
              className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500 resize-y font-mono text-sm"
            />
            <button
              onClick={handleCreate}
              disabled={creating || !title.trim() || !content.trim()}
              className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-50"
            >
              {creating ? "Creating..." : "Create CV"}
            </button>
          </div>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : cvs.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 dark:border-white/10 p-20 text-center bg-white/5">
          <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center mb-6">
            <svg className="w-10 h-10 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">No CVs yet</h2>
          <p className="text-gray-500 dark:text-gray-400 max-w-md">Create a CV to start matching with job descriptions.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {cvs.map((cv) => (
            <div key={cv.id} className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-5 hover:border-purple-500/30 transition-all duration-200">
              <div className="flex items-start justify-between mb-3">
                <Link href={`/dashboard/resumes/${cv.id}`} className="font-semibold text-gray-900 dark:text-white hover:text-purple-400">
                  {cv.title}
                </Link>
                <span className={`text-xs px-3 py-1 rounded-full font-medium ${
                  cv.status === "completed" ? "bg-green-500/10 text-green-500" :
                  cv.status === "pending" ? "bg-yellow-500/10 text-yellow-500" :
                  "bg-red-500/10 text-red-500"
                }`}>
                  {cv.status}
                </span>
              </div>
              {cv.parsed_summary && (
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-3 line-clamp-2">{cv.parsed_summary}</p>
              )}
              {cv.parsed_skills && cv.parsed_skills.length > 0 && (
                <div className="flex flex-wrap gap-1 mb-3">
                  {cv.parsed_skills.slice(0, 5).map((s, i) => (
                    <span key={i} className="text-xs px-2 py-1 rounded-full bg-purple-500/10 text-purple-400">{s}</span>
                  ))}
                  {cv.parsed_skills.length > 5 && (
                    <span className="text-xs px-2 py-1 rounded-full bg-gray-500/10 text-gray-400">+{cv.parsed_skills.length - 5}</span>
                  )}
                </div>
              )}
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500">{new Date(cv.created_at).toLocaleDateString()}</span>
                <button onClick={() => handleDelete(cv.id)} className="text-xs text-red-400 hover:text-red-300">Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
