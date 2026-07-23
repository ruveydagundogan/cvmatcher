"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

interface JDDetail {
  id: string;
  title: string;
  content: string;
  required_skills: string[];
  preferred_skills: string[];
  experience_level: string;
  employment_type: string;
  location: string;
  created_at: string;
  updated_at: string;
}

export default function JDDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [jd, setJD] = useState<JDDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [analyzing, setAnalyzing] = useState(false);
  const [editing, setEditing] = useState(false);
  const [editData, setEditData] = useState<Partial<JDDetail>>({});

  const fetchJD = () => {
    setLoading(true);
    api.get(`/api/v1/jds/${id}`)
      .then((data) => { setJD(data); setEditData(data); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchJD(); }, [id]);

  const handleAnalyze = async () => {
    setAnalyzing(true);
    try {
      await api.post(`/api/v1/jds/${id}/analyze`);
      fetchJD();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setAnalyzing(false);
    }
  };

  const handleSave = async () => {
    try {
      await api.put(`/api/v1/jds/${id}`, editData);
      setEditing(false);
      fetchJD();
    } catch (e: any) {
      setError(e.message);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  if (!jd) {
    return <div className="text-gray-500">Job description not found</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <Link href="/dashboard/jds" className="text-sm text-purple-400 hover:text-purple-300 mb-2 inline-block">&larr; Back to JDs</Link>
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{jd.title}</h1>
          <div className="flex items-center gap-3">
            {!editing && (
              <>
                <button onClick={() => setEditing(true)} className="px-4 py-2 rounded-xl bg-white/10 text-gray-300 hover:bg-white/20 font-medium transition-all duration-200 border border-white/10">
                  Edit
                </button>
                <button
                  onClick={handleAnalyze}
                  disabled={analyzing}
                  className="px-4 py-2 rounded-xl bg-gradient-to-r from-purple-500 to-pink-600 text-white font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-50"
                >
                  {analyzing ? "Analyzing..." : "Analyze with AI"}
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Description</h2>
            <pre className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap font-mono max-h-96 overflow-y-auto">{jd.content}</pre>
          </div>
        </div>

        <div className="space-y-6">
          {editing ? (
            <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Edit</h2>
              <div className="space-y-4">
                <div>
                  <label className="text-sm text-gray-500 mb-1 block">Title</label>
                  <input type="text" value={editData.title || ""} onChange={(e) => setEditData({...editData, title: e.target.value})}
                    className="w-full px-3 py-2 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500" />
                </div>
                <div>
                  <label className="text-sm text-gray-500 mb-1 block">Experience Level</label>
                  <input type="text" value={editData.experience_level || ""} onChange={(e) => setEditData({...editData, experience_level: e.target.value})}
                    className="w-full px-3 py-2 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500" />
                </div>
                <div>
                  <label className="text-sm text-gray-500 mb-1 block">Employment Type</label>
                  <input type="text" value={editData.employment_type || ""} onChange={(e) => setEditData({...editData, employment_type: e.target.value})}
                    className="w-full px-3 py-2 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500" />
                </div>
                <div>
                  <label className="text-sm text-gray-500 mb-1 block">Location</label>
                  <input type="text" value={editData.location || ""} onChange={(e) => setEditData({...editData, location: e.target.value})}
                    className="w-full px-3 py-2 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-purple-500" />
                </div>
                <div className="flex gap-2">
                  <button onClick={handleSave} className="px-4 py-2 rounded-xl bg-green-600 text-white font-medium hover:opacity-90">Save</button>
                  <button onClick={() => { setEditing(false); setEditData(jd); }} className="px-4 py-2 rounded-xl bg-white/10 text-gray-300 hover:bg-white/20">Cancel</button>
                </div>
              </div>
            </div>
          ) : (
            <>
              {jd.required_skills && jd.required_skills.length > 0 && (
                <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Required Skills ({jd.required_skills.length})</h2>
                  <div className="flex flex-wrap gap-2">
                    {jd.required_skills.map((s, i) => (
                      <span key={i} className="px-3 py-1.5 rounded-xl bg-red-500/10 text-red-400 text-sm border border-red-500/20">{s}</span>
                    ))}
                  </div>
                </div>
              )}
              {jd.preferred_skills && jd.preferred_skills.length > 0 && (
                <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Preferred Skills ({jd.preferred_skills.length})</h2>
                  <div className="flex flex-wrap gap-2">
                    {jd.preferred_skills.map((s, i) => (
                      <span key={i} className="px-3 py-1.5 rounded-xl bg-blue-500/10 text-blue-400 text-sm border border-blue-500/20">{s}</span>
                    ))}
                  </div>
                </div>
              )}
              <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Details</h2>
                <div className="space-y-2 text-sm">
                  {jd.experience_level && <p><span className="text-gray-500">Level:</span> <span className="text-gray-900 dark:text-white">{jd.experience_level}</span></p>}
                  {jd.employment_type && <p><span className="text-gray-500">Type:</span> <span className="text-gray-900 dark:text-white">{jd.employment_type}</span></p>}
                  {jd.location && <p><span className="text-gray-500">Location:</span> <span className="text-gray-900 dark:text-white">{jd.location}</span></p>}
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
