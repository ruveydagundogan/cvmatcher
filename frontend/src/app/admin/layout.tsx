"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Sidebar } from "@/components/Sidebar";
import { Header } from "@/components/Header";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    const role = localStorage.getItem("userRole");
    if (!token) {
      router.replace("/");
    } else if (role !== "admin") {
      router.replace(role === "hr" ? "/ik" : "/dashboard");
    }
  }, [router]);

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 ml-64 flex flex-col">
        <Header />
        <main className="flex-1 p-6 overflow-auto">
          <div className="flex gap-2 mb-6">
            <Link href="/admin/adapters" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Adapters</Link>
            <Link href="/admin/prompts" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Prompts</Link>
            <Link href="/admin/settings" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Settings</Link>
            <Link href="/admin/logs" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Logs</Link>
          </div>
          {children}
        </main>
      </div>
    </div>
  );
}
