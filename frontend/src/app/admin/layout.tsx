"use client";

import { Sidebar } from "@/components/Sidebar";
import { Header } from "@/components/Header";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 ml-64 flex flex-col">
        <Header />
        <main className="flex-1 p-6 overflow-auto">
          <div className="flex gap-2 mb-6">
            <a href="/admin/adapters" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Adapters</a>
            <a href="/admin/prompts" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Prompts</a>
            <a href="/admin/settings" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Settings</a>
            <a href="/admin/logs" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Logs</a>
          </div>
          {children}
        </main>
      </div>
    </div>
  );
}
