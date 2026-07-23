"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

interface DashboardStats {
  total_cvs: number;
  total_jds: number;
  total_matches: number;
  average_score: number;
  match_rate: number;
  recent_matches: Array<{
    id: string;
    cv_title: string;
    jd_title: string;
    overall_score: number;
    created_at: string;
  }>;
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get("/api/v1/dashboard/stats")
      .then(setStats)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const statCards = [
    { label: "CVs", value: stats?.total_cvs ?? "-", href: "/dashboard/cvs", color: "from-blue-500 to-cyan-600", desc: "Upload & parse resumes" },
    { label: "Job Descriptions", value: stats?.total_jds ?? "-", href: "/dashboard/jds", color: "from-purple-500 to-pink-600", desc: "Add job postings" },
    { label: "Matches", value: stats?.total_matches ?? "-", href: "/dashboard/matches", color: "from-green-500 to-emerald-600", desc: "Run AI matching" },
    { label: "Avg Score", value: stats ? (stats.average_score ? stats.average_score.toFixed(1) : "—") : "-", href: "/dashboard/matches", color: "from-orange-500 to-red-600", desc: "Average match score" },
  ];

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-1">CV–Job Description AI Matching Platform</p>
      </div>

      {error && (
        <div className="rounded-xl bg-red-500/10 border border-red-500/20 p-4 mb-6">
          <p className="text-sm text-red-400">{error}</p>
          <p className="text-xs text-red-400/70 mt-1">
            Backend'e bağlanılamadı. Vercel'de <code className="bg-red-500/10 px-1 rounded">NEXT_PUBLIC_API_URL</code> environment variable'ını Render URL'ine ayarla.
          </p>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-32">
          <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {statCards.map((card) => (
            <Link key={card.label} href={card.href}>
              <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 hover:border-purple-500/30 transition-all duration-200 h-full">
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">{card.desc}</p>
                <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">{card.label}</p>
                <p className={`text-4xl font-bold bg-gradient-to-r ${card.color} bg-clip-text text-transparent`}>
                  {card.value}
                </p>
              </div>
            </Link>
          ))}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <Link href="/dashboard/cvs" className="bg-gradient-to-br from-blue-500/10 to-cyan-500/10 rounded-2xl border border-blue-500/20 p-6 hover:border-blue-500/40 transition-all duration-200">
          <h3 className="font-semibold text-gray-900 dark:text-white mb-1">CV Ekle</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">CV metnini yapıştır ve AI ile parse et</p>
        </Link>
        <Link href="/dashboard/jds" className="bg-gradient-to-br from-purple-500/10 to-pink-500/10 rounded-2xl border border-purple-500/20 p-6 hover:border-purple-500/40 transition-all duration-200">
          <h3 className="font-semibold text-gray-900 dark:text-white mb-1">İş Tanımı Ekle</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">Job description ekle ve AI ile analiz et</p>
        </Link>
        <Link href="/dashboard/matches" className="bg-gradient-to-br from-green-500/10 to-emerald-500/10 rounded-2xl border border-green-500/20 p-6 hover:border-green-500/40 transition-all duration-200">
          <h3 className="font-semibold text-gray-900 dark:text-white mb-1">Match Çalıştır</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">CV ile iş tanımını eşleştir ve skorla</p>
        </Link>
      </div>

      {stats && stats.recent_matches && stats.recent_matches.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Son Match'ler</h2>
          <div className="space-y-3">
            {stats.recent_matches.map((m) => (
              <Link key={m.id} href={`/dashboard/matches/${m.id}`}>
                <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-4 hover:border-purple-500/30 transition-all duration-200 flex items-center justify-between">
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">{m.cv_title} ↔ {m.jd_title}</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">{new Date(m.created_at).toLocaleDateString()}</p>
                  </div>
                  <div className={`px-4 py-2 rounded-xl font-bold ${
                    m.overall_score >= 70 ? "bg-green-500/10 text-green-500" :
                    m.overall_score >= 40 ? "bg-yellow-500/10 text-yellow-500" :
                    "bg-red-500/10 text-red-500"
                  }`}>
                    {m.overall_score.toFixed(0)}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {!loading && !error && (!stats || stats.total_cvs === 0) && (
        <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 dark:border-white/10 p-16 text-center bg-white/5">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Hoş geldin!</h2>
          <p className="text-gray-500 dark:text-gray-400 max-w-md">
            İlk CV'ni ekleyerek başla, ardından bir iş tanımı gir ve AI ile match skorunu gör.
          </p>
        </div>
      )}
    </div>
  );
}
