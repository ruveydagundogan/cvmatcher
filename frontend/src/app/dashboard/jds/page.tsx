"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

interface JD {
  id: string;
  title: string;
  required_skills: string[];
  preferred_skills: string[];
  experience_level: string;
  employment_type: string;
  location: string;
  created_at: string;
}

export default function JDListPage() {
  const [jds, setJDs] = useState<JD[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showCreate, setShowCreate] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [creating, setCreating] = useState(false);

  const fetchJDs = () => {
    setLoading(true);
    api.get("/api/v1/jds?page=1&limit=50")
      .then((data) => { setJDs(data.items || []); setTotal(data.total || 0); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchJDs(); }, []);

  const handleCreate = async () => {
    if (!title.trim() || !content.trim()) return;
    setCreating(true);
    try {
      await api.post("/api/v1/jds", { title, content });
      setTitle("");
      setContent("");
      setShowCreate(false);
      fetchJDs();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this job description?")) return;
    try {
      await api.delete(`/api/v1/jds/${id}`);
      fetchJDs();
    } catch (e: any) {
      setError(e.message);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Job Descriptions</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">{total} job description{jds.length !== 1 ? "s" : ""}</p>
        </div>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-500 to-pink-600 text-white font-medium hover:opacity-90 transition-all duration-200"
        >
          {showCreate ? "Cancel" : "New JD"}
        </button>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      {showCreate && (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Create Job Description</h2>
          <div className="space-y-4">
            <input
              type="text"
              placeholder="Job Title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
            <textarea
              placeholder="Paste job description content here..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={10}
              className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500 resize-y font-mono text-sm"
            />
            <button
              onClick={handleCreate}
              disabled={creating || !title.trim() || !content.trim()}
              className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-500 to-pink-600 text-white font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-50"
            >
              {creating ? "Creating..." : "Create JD"}
            </button>
          </div>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : jds.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 dark:border-white/10 p-20 text-center bg-white/5">
          <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center mb-6">
            <svg className="w-10 h-10 text-pink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">No job descriptions yet</h2>
          <p className="text-gray-500 dark:text-gray-400 max-w-md">Create a job description to match against your CVs.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {jds.map((jd) => (
            <div key={jd.id} className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-5 hover:border-purple-500/30 transition-all duration-200">
              <div className="flex items-start justify-between mb-3">
                <Link href={`/dashboard/jds/${jd.id}`} className="font-semibold text-gray-900 dark:text-white hover:text-purple-400">
                  {jd.title}
                </Link>
              </div>
              <div className="flex flex-wrap gap-2 mb-3">
                {jd.experience_level && (
                  <span className="text-xs px-2 py-1 rounded-full bg-blue-500/10 text-blue-400">{jd.experience_level}</span>
                )}
                {jd.employment_type && (
                  <span className="text-xs px-2 py-1 rounded-full bg-green-500/10 text-green-400">{jd.employment_type}</span>
                )}
                {jd.location && (
                  <span className="text-xs px-2 py-1 rounded-full bg-orange-500/10 text-orange-400">{jd.location}</span>
                )}
              </div>
              {jd.required_skills && jd.required_skills.length > 0 && (
                <div className="flex flex-wrap gap-1 mb-3">
                  {jd.required_skills.slice(0, 5).map((s, i) => (
                    <span key={i} className="text-xs px-2 py-1 rounded-full bg-purple-500/10 text-purple-400">{s}</span>
                  ))}
                  {jd.required_skills.length > 5 && (
                    <span className="text-xs px-2 py-1 rounded-full bg-gray-500/10 text-gray-400">+{jd.required_skills.length - 5}</span>
                  )}
                </div>
              )}
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500">{new Date(jd.created_at).toLocaleDateString()}</span>
                <button onClick={() => handleDelete(jd.id)} className="text-xs text-red-400 hover:text-red-300">Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
