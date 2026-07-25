"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

interface MatchDetail {
  id: string;
  cv_id: string;
  jd_id: string;
  cv_title: string;
  jd_title: string;
  cv_file_name: string;
  overall_score: number;
  skill_match_score: number;
  experience_score: number;
  education_score: number;
  analysis: string;
  matched_skills: string[];
  missing_skills: string[];
  created_at: string;
}

function Gauge({ value, label, color }: { value: number; label: string; color: string }) {
  const pct = Math.round(value * 100);
  const radius = 54;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (pct / 100) * circumference;

  return (
    <div className="flex flex-col items-center">
      <div className="relative flex items-center justify-center" style={{ width: 130, height: 130 }}>
        <svg width="130" height="130" className="absolute transform -rotate-90">
          <circle cx="65" cy="65" r={radius} fill="none" stroke="currentColor" strokeWidth="8" className="text-gray-200 dark:text-gray-700" />
          <circle cx="65" cy="65" r={radius} fill="none" stroke="currentColor" strokeWidth="8"
            strokeDasharray={circumference} strokeDashoffset={offset}
            strokeLinecap="round" className={`${color} transition-all duration-1000`}
          />
        </svg>
        <span className="text-2xl font-bold text-gray-900 dark:text-white z-10">{pct}</span>
      </div>
      <span className="text-xs text-gray-500 mt-2">{label}</span>
    </div>
  );
}

export default function MatchDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;
  const [match, setMatch] = useState<MatchDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    api.get(`/api/v1/matches/${id}`)
      .then(setMatch)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  const remove = async () => {
    if (!confirm("Delete this match?")) return;
    try {
      await api.delete(`/api/v1/matches/${id}`);
      router.push("/dashboard/matches");
    } catch (e) {
      console.error(e);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  if (!match) {
    return <div className="text-gray-500">Match not found</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <Link href="/dashboard/matches" className="text-sm text-purple-400 hover:text-purple-300 mb-2 inline-block">&larr; Back to Matches</Link>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Match Result</h1>
            <p className="text-gray-500 dark:text-gray-400 mt-1">
              {match.cv_title} ↔ {match.jd_title}
            </p>
          </div>
          <button onClick={remove}
            className="px-4 py-2 rounded-xl bg-red-500/10 text-red-400 hover:bg-red-500/20 text-sm font-medium">
            Delete Match
          </button>
        </div>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-8 flex items-center justify-center">
          <Gauge value={match.overall_score} label="Overall Match" color="text-purple-500" />
        </div>

        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">Score Breakdown</h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-500">Skill Match</span>
                <span className="text-gray-900 dark:text-white font-medium">{(match.skill_match_score * 100).toFixed(0)}%</span>
              </div>
              <div className="h-2 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
                <div className="h-full rounded-full bg-blue-500 transition-all duration-500" style={{ width: `${match.skill_match_score * 100}%` }} />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-500">Experience</span>
                <span className="text-gray-900 dark:text-white font-medium">{(match.experience_score * 100).toFixed(0)}%</span>
              </div>
              <div className="h-2 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
                <div className="h-full rounded-full bg-green-500 transition-all duration-500" style={{ width: `${match.experience_score * 100}%` }} />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-500">Education</span>
                <span className="text-gray-900 dark:text-white font-medium">{(match.education_score * 100).toFixed(0)}%</span>
              </div>
              <div className="h-2 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
                <div className="h-full rounded-full bg-orange-500 transition-all duration-500" style={{ width: `${match.education_score * 100}%` }} />
              </div>
            </div>
          </div>
        </div>
      </div>

      {match.analysis && (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">AI Analysis</h2>
          <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{match.analysis}</p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {match.matched_skills && match.matched_skills.length > 0 && (
          <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Matched Skills ({match.matched_skills.length})</h2>
            <div className="flex flex-wrap gap-2">
              {match.matched_skills.map((s, i) => (
                <span key={i} className="px-3 py-1.5 rounded-xl bg-green-500/10 text-green-400 text-sm border border-green-500/20">{s}</span>
              ))}
            </div>
          </div>
        )}

        {match.missing_skills && match.missing_skills.length > 0 && (
          <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Missing Skills ({match.missing_skills.length})</h2>
            <div className="flex flex-wrap gap-2">
              {match.missing_skills.map((s, i) => (
                <span key={i} className="px-3 py-1.5 rounded-xl bg-red-500/10 text-red-400 text-sm border border-red-500/20">{s}</span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
