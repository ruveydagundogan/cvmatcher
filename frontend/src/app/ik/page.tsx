"use client";

import Link from "next/link";

export default function IKDashboardPage() {
  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-white">IK Dashboard</h1>
        <p className="text-gray-400 mt-1">Manage candidates, job descriptions and hiring pipeline</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white/5 backdrop-blur-xl rounded-3xl border border-white/10 p-6">
          <h3 className="text-lg font-semibold text-white mb-2">Candidates</h3>
          <p className="text-gray-400 text-sm mb-4">
            Add candidate CVs and let AI parse their skills, experience and education.
          </p>
          <span className="inline-flex px-3 py-1 rounded-full bg-blue-500/20 text-blue-300 text-xs font-medium">
            Coming soon (Faz 2)
          </span>
        </div>

        <div className="bg-white/5 backdrop-blur-xl rounded-3xl border border-white/10 p-6">
          <h3 className="text-lg font-semibold text-white mb-2">Job Descriptions</h3>
          <p className="text-gray-400 text-sm mb-4">
            Create job postings and find the best matching candidates automatically.
          </p>
          <span className="inline-flex px-3 py-1 rounded-full bg-blue-500/20 text-blue-300 text-xs font-medium">
            Coming soon (Faz 2)
          </span>
        </div>

        <div className="bg-white/5 backdrop-blur-xl rounded-3xl border border-white/10 p-6">
          <h3 className="text-lg font-semibold text-white mb-2">Pipeline</h3>
          <p className="text-gray-400 text-sm mb-4">
            Track candidates through screening, interviews and hiring decisions.
          </p>
          <span className="inline-flex px-3 py-1 rounded-full bg-blue-500/20 text-blue-300 text-xs font-medium">
            Coming soon (Faz 2)
          </span>
        </div>
      </div>

      <div className="mt-8 p-6 bg-white/5 backdrop-blur-xl rounded-3xl border border-white/10">
        <h3 className="text-lg font-semibold text-white mb-2">Need help?</h3>
        <p className="text-gray-400 text-sm">
          The IK Portal features are planned for the next phase. Meanwhile, you can use the{" "}
          <Link href="/dashboard" className="text-blue-400 hover:text-blue-300 font-medium">
            personal dashboard
          </Link>{" "}
          for CV management.
        </p>
      </div>
    </div>
  );
}
