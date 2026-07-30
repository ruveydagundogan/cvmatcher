"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

interface Adapter {
  id: string;
  name: string;
  description: string;
  file_path: string;
  active: boolean;
  model_name: string;
}

export default function AdaptersPage() {
  const [adapters, setAdapters] = useState<Adapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState("");
  const [desc, setDesc] = useState("");
  const [filePath, setFilePath] = useState("");
  const [modelName, setModelName] = useState("qwen2.5:1.5b-instruct");

  const load = async () => {
    try {
      setAdapters(await api.get("/api/v1/admin/adapters") || []);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const create = async () => {
    if (!name.trim()) return;
    try {
      await api.post("/api/v1/admin/adapters", { name, description: desc, file_path: filePath, model_name: modelName });
      setName(""); setDesc(""); setFilePath(""); setModelName("qwen2.5:1.5b-instruct");
      load();
    } catch (e) {
      console.error(e);
    }
  };

  const remove = async (id: string) => {
    try {
      await api.delete(`/api/v1/admin/adapters/${id}`);
      load();
    } catch (e) {
      console.error(e);
    }
  };

  const activate = async (id: string) => {
    try {
      await api.post(`/api/v1/admin/adapters/${id}/activate`, {});
      load();
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Adapter Management</h1>

      <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 p-6 mb-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Load New Adapter</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input type="text" value={name} onChange={(e) => setName(e.target.value)}
            placeholder="Adapter name" className="px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <input type="text" value={desc} onChange={(e) => setDesc(e.target.value)}
            placeholder="Description" className="px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <input type="text" value={filePath} onChange={(e) => setFilePath(e.target.value)}
            placeholder="File path (e.g. /adapters/lora-v1)" className="px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
          <input type="text" value={modelName} onChange={(e) => setModelName(e.target.value)}
            placeholder="Model name" className="px-4 py-2.5 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-gray-900 dark:text-white" />
        </div>
        <button onClick={create}
          className="mt-4 px-6 py-2.5 rounded-xl bg-purple-600 text-white font-medium hover:bg-purple-700">
          Load Adapter
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : adapters.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No adapters loaded</div>
      ) : (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-white/10">
                <th className="text-left px-6 py-3 text-gray-500 font-medium">Name</th>
                <th className="text-left px-6 py-3 text-gray-500 font-medium">Model</th>
                <th className="text-left px-6 py-3 text-gray-500 font-medium">Status</th>
                <th className="text-right px-6 py-3 text-gray-500 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {adapters.map((a) => (
                <tr key={a.id} className="border-b border-gray-100 dark:border-white/5">
                  <td className="px-6 py-4 text-gray-900 dark:text-white">{a.name}</td>
                  <td className="px-6 py-4 text-gray-500">{a.model_name}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs px-2 py-0.5 rounded-full ${a.active ? "bg-green-500/10 text-green-400" : "bg-gray-500/10 text-gray-400"}`}>
                      {a.active ? "Active" : "Inactive"}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right space-x-3">
                    {!a.active && (
                      <button onClick={() => activate(a.id)} className="text-xs text-green-400 hover:text-green-300">Activate</button>
                    )}
                    <button onClick={() => remove(a.id)} className="text-xs text-red-400 hover:text-red-300">Remove</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
