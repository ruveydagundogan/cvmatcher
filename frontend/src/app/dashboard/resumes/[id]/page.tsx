"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

interface ParsedExperience {
  title: string;
  company: string;
  start_date: string;
  end_date: string;
  description: string;
}

interface ParsedEducation {
  degree: string;
  field: string;
  institution: string;
  start_year: string;
  end_year: string;
}

interface CVDetail {
  id: string;
  title: string;
  content: string;
  status: string;
  parsed_skills: string[];
  parsed_experience: ParsedExperience[];
  parsed_education: ParsedEducation[];
  parsed_summary: string;
  created_at: string;
  updated_at: string;
}

export default function CVDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [cv, setCV] = useState<CVDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [parsing, setParsing] = useState(false);

  const fetchCV = () => {
    setLoading(true);
    api.get(`/api/v1/cvs/${id}`)
      .then(setCV)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchCV(); }, [id]);

  const handleParse = async () => {
    setParsing(true);
    try {
      await api.post(`/api/v1/cvs/${id}/parse`);
      fetchCV();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setParsing(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  if (!cv) {
    return <div className="text-gray-500">CV not found</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <Link href="/dashboard/resumes" className="text-sm text-purple-400 hover:text-purple-300 mb-2 inline-block">&larr; Back to CVs</Link>
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{cv.title}</h1>
          <div className="flex items-center gap-3">
            <span className={`text-sm px-3 py-1 rounded-full font-medium ${
              cv.status === "completed" ? "bg-green-500/10 text-green-500" :
              cv.status === "pending" ? "bg-yellow-500/10 text-yellow-500" :
              "bg-red-500/10 text-red-500"
            }`}>{cv.status}</span>
            {cv.status !== "completed" && (
              <button
                onClick={handleParse}
                disabled={parsing}
                className="px-4 py-2 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-50"
              >
                {parsing ? "Parsing..." : "Parse with AI"}
              </button>
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
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Raw Content</h2>
            <pre className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap font-mono max-h-96 overflow-y-auto">{cv.content}</pre>
          </div>

          {cv.parsed_experience && cv.parsed_experience.length > 0 && (
            <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Experience</h2>
              <div className="space-y-4">
                {cv.parsed_experience.map((exp, i) => (
                  <div key={i} className="border-l-2 border-purple-500/30 pl-4">
                    <p className="font-medium text-gray-900 dark:text-white">{exp.title}</p>
                    <p className="text-sm text-purple-400">{exp.company}</p>
                    <p className="text-xs text-gray-500">{exp.start_date} - {exp.end_date}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">{exp.description}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {cv.parsed_education && cv.parsed_education.length > 0 && (
            <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Education</h2>
              <div className="space-y-4">
                {cv.parsed_education.map((edu, i) => (
                  <div key={i} className="border-l-2 border-blue-500/30 pl-4">
                    <p className="font-medium text-gray-900 dark:text-white">{edu.degree} in {edu.field}</p>
                    <p className="text-sm text-blue-400">{edu.institution}</p>
                    <p className="text-xs text-gray-500">{edu.start_year} - {edu.end_year}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {cv.parsed_summary && (
            <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Summary</h2>
              <p className="text-sm text-gray-700 dark:text-gray-300">{cv.parsed_summary}</p>
            </div>
          )}

          {cv.parsed_skills && cv.parsed_skills.length > 0 && (
            <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Skills ({cv.parsed_skills.length})</h2>
              <div className="flex flex-wrap gap-2">
                {cv.parsed_skills.map((s, i) => (
                  <span key={i} className="px-3 py-1.5 rounded-xl bg-purple-500/10 text-purple-400 text-sm border border-purple-500/20">{s}</span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
