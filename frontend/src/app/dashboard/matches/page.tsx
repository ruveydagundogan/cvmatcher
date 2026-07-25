"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

interface MatchResult {
  id: string;
  cv_id: string;
  jd_id: string;
  cv_title: string;
  jd_title: string;
  overall_score: number;
  skill_match_score: number;
  experience_score: number;
  education_score: number;
  analysis: string;
  matched_skills: string[];
  missing_skills: string[];
  created_at: string;
}

interface CV {
  id: string;
  title: string;
}

interface JD {
  id: string;
  title: string;
}

export default function MatchListPage() {
  const [matches, setMatches] = useState<MatchResult[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showMatch, setShowMatch] = useState(false);
  const [cvs, setCVs] = useState<CV[]>([]);
  const [jds, setJDs] = useState<JD[]>([]);
  const [cvId, setCvId] = useState("");
  const [jdId, setJdId] = useState("");
  const [matching, setMatching] = useState(false);

  const fetchMatches = () => {
    setLoading(true);
    api.get("/api/v1/matches?page=1&limit=50")
      .then((data) => { setMatches(data.items || []); setTotal(data.total || 0); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchMatches(); }, []);

  const openMatchForm = async () => {
    setShowMatch(true);
    try {
      const [cvData, jdData] = await Promise.all([
        api.get("/api/v1/cvs?page=1&limit=100"),
        api.get("/api/v1/jds?page=1&limit=100"),
      ]);
      setCVs(cvData.items || []);
      setJDs(jdData.items || []);
    } catch (e: any) {
      setError(e.message);
    }
  };

  const handleMatch = async () => {
    if (!cvId || !jdId) return;
    setMatching(true);
    try {
      await api.post("/api/v1/matches", { cv_id: cvId, jd_id: jdId });
      setCvId("");
      setJdId("");
      setShowMatch(false);
      fetchMatches();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setMatching(false);
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 0.7) return "from-green-500 to-emerald-600 bg-green-500/10 border-green-500/20 text-green-500";
    if (score >= 0.4) return "from-yellow-500 to-orange-600 bg-yellow-500/10 border-yellow-500/20 text-yellow-500";
    return "from-red-500 to-rose-600 bg-red-500/10 border-red-500/20 text-red-500";
  };

  const removeMatch = async (id: string, e: React.MouseEvent) => {
    e.preventDefault();
    if (!confirm("Delete this match?")) return;
    try {
      await api.delete(`/api/v1/matches/${id}`);
      setMatches(matches.filter((m) => m.id !== id));
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Matches</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">{total} match{matches.length !== 1 ? "es" : ""}</p>
        </div>
        <button
          onClick={openMatchForm}
          className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-green-500 to-emerald-600 text-white font-medium hover:opacity-90 transition-all duration-200"
        >
          New Match
        </button>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      {showMatch && (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Run New Match</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="text-sm text-gray-500 mb-2 block">Select CV</label>
              <select
                value={cvId}
                onChange={(e) => setCvId(e.target.value)}
                className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                <option value="">Choose a CV...</option>
                {cvs.map((cv) => (
                  <option key={cv.id} value={cv.id}>{cv.title}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-sm text-gray-500 mb-2 block">Select Job Description</label>
              <select
                value={jdId}
                onChange={(e) => setJdId(e.target.value)}
                className="w-full px-4 py-3 rounded-xl bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                <option value="">Choose a JD...</option>
                {jds.map((jd) => (
                  <option key={jd.id} value={jd.id}>{jd.title}</option>
                ))}
              </select>
            </div>
          </div>
          <div className="flex gap-3">
            <button
              onClick={handleMatch}
              disabled={matching || !cvId || !jdId}
              className="px-5 py-2.5 rounded-xl bg-gradient-to-r from-green-500 to-emerald-600 text-white font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-50"
            >
              {matching ? "Matching..." : "Run Match"}
            </button>
            <button onClick={() => setShowMatch(false)} className="px-5 py-2.5 rounded-xl bg-white/10 text-gray-300 hover:bg-white/20">Cancel</button>
          </div>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : matches.length === 0 && !showMatch ? (
        <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 dark:border-white/10 p-20 text-center bg-white/5">
          <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center mb-6">
            <svg className="w-10 h-10 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">No matches yet</h2>
          <p className="text-gray-500 dark:text-gray-400 max-w-md">Run a match between a CV and a job description to see results.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {matches.map((m) => (
            <Link key={m.id} href={`/dashboard/matches/${m.id}`}>
              <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-5 hover:border-purple-500/30 transition-all duration-200 relative group">
                <button onClick={(e) => removeMatch(m.id, e)}
                  className="absolute top-3 right-3 w-7 h-7 rounded-full bg-red-500/10 text-red-400 hover:bg-red-500/20 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity text-xs">
                  ✕
                </button>
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                      <span className="text-white text-sm font-bold">M</span>
                    </div>
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">{m.cv_title} ↔ {m.jd_title}</p>
                      <p className="text-xs text-gray-500">{new Date(m.created_at).toLocaleDateString()}</p>
                    </div>
                  </div>
                  <div className={`px-4 py-2 rounded-xl border font-bold bg-gradient-to-r ${getScoreColor(m.overall_score)}`}>
                    <span className="text-white">{(m.overall_score * 100).toFixed(0)}</span>
                    <span className="text-white/60 text-sm">/100</span>
                  </div>
                </div>

                {m.analysis && (
                  <p className="text-sm text-gray-500 dark:text-gray-400 line-clamp-2 mb-3">{m.analysis}</p>
                )}

                <div className="flex flex-wrap gap-2">
                  {m.matched_skills && m.matched_skills.length > 0 && (
                    <div className="flex flex-wrap gap-1 items-center">
                      <span className="text-xs text-gray-500 mr-1">Matched:</span>
                      {m.matched_skills.slice(0, 3).map((s, i) => (
                        <span key={i} className="text-xs px-2 py-0.5 rounded-full bg-green-500/10 text-green-400">{s}</span>
                      ))}
                      {m.matched_skills.length > 3 && (
                        <span className="text-xs text-gray-500">+{m.matched_skills.length - 3}</span>
                      )}
                    </div>
                  )}
                  {m.missing_skills && m.missing_skills.length > 0 && (
                    <div className="flex flex-wrap gap-1 items-center">
                      <span className="text-xs text-gray-500 mr-1">Missing:</span>
                      {m.missing_skills.slice(0, 3).map((s, i) => (
                        <span key={i} className="text-xs px-2 py-0.5 rounded-full bg-red-500/10 text-red-400">{s}</span>
                      ))}
                      {m.missing_skills.length > 3 && (
                        <span className="text-xs text-gray-500">+{m.missing_skills.length - 3}</span>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
