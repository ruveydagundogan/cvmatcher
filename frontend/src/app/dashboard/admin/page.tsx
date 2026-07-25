"use client";

import Link from "next/link";

export default function AdminPage() {
  const cards = [
    { href: "/dashboard/admin/adapters", title: "Adapters", desc: "Manage PEFT/LoRA adapters for hot-swap model tuning" },
    { href: "/dashboard/admin/prompts", title: "System Prompts", desc: "Create and activate system prompt templates" },
    { href: "/dashboard/admin/settings", title: "LLM Settings", desc: "Configure temperature, tokens, context length" },
    { href: "/dashboard/admin/logs", title: "Query Logs", desc: "Monitor LLM query history and latency" },
  ];

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Admin Panel</h1>
      <p className="text-gray-500 dark:text-gray-400 mb-8">Manage your LLM engine — adapters, prompts, settings, and monitoring</p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {cards.map((c) => (
          <Link key={c.href} href={c.href}
            className="block bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 hover:border-purple-500/50 transition-colors">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">{c.title}</h2>
            <p className="text-sm text-gray-500">{c.desc}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
