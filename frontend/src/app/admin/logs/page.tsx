"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

interface LogEntry {
  id: string;
  user_id: string;
  query: string;
  response: string;
  model: string;
  adapter: string;
  duration_ms: number;
  token_count: number;
  status: string;
  created_at: string;
}

export default function LogsPage() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<LogEntry | null>(null);

  useEffect(() => {
    api.get("/api/v1/admin/logs")
      .then((data) => {
        setLogs(data.items || []);
        setTotal(data.total || 0);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const formatDate = (d: string) => new Date(d).toLocaleString();

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Query Logs</h1>
      <p className="text-sm text-gray-500 mb-6">{total} total queries logged</p>

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full" />
        </div>
      ) : logs.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No logs yet</div>
      ) : (
        <div className="bg-white dark:bg-slate-900/50 rounded-2xl border border-gray-200 dark:border-white/10 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-white/10">
                <th className="text-left px-4 py-3 text-gray-500 font-medium">Time</th>
                <th className="text-left px-4 py-3 text-gray-500 font-medium">Model</th>
                <th className="text-left px-4 py-3 text-gray-500 font-medium">Query</th>
                <th className="text-right px-4 py-3 text-gray-500 font-medium">Duration</th>
                <th className="text-right px-4 py-3 text-gray-500 font-medium">Tokens</th>
                <th className="text-right px-4 py-3 text-gray-500 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id}
                  onClick={() => setSelected(selected?.id === log.id ? null : log)}
                  className="border-b border-gray-100 dark:border-white/5 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800/50">
                  <td className="px-4 py-3 text-gray-500 text-xs">{formatDate(log.created_at)}</td>
                  <td className="px-4 py-3 text-gray-900 dark:text-white">{log.model}</td>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-300 max-w-xs truncate">{log.query}</td>
                  <td className="px-4 py-3 text-right text-gray-500">{log.duration_ms}ms</td>
                  <td className="px-4 py-3 text-right text-gray-500">{log.token_count}</td>
                  <td className="px-4 py-3 text-right">
                    <span className={`text-xs px-2 py-0.5 rounded-full ${log.status === "success" ? "bg-green-500/10 text-green-400" : "bg-red-500/10 text-red-400"}`}>
                      {log.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={() => setSelected(null)}>
          <div className="bg-white dark:bg-slate-900 rounded-2xl p-6 max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Log Detail</h2>
              <button onClick={() => setSelected(null)} className="text-gray-500 hover:text-gray-700">&times;</button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="text-xs text-gray-500 uppercase tracking-wider">Query</label>
                <p className="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-slate-800 rounded-xl p-3 mt-1">{selected.query}</p>
              </div>
              <div>
                <label className="text-xs text-gray-500 uppercase tracking-wider">Response</label>
                <p className="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-slate-800 rounded-xl p-3 mt-1 whitespace-pre-wrap max-h-60 overflow-y-auto">{selected.response}</p>
              </div>
              <div className="grid grid-cols-3 gap-4 text-sm">
                <div><span className="text-gray-500">Model:</span> <span className="text-gray-900 dark:text-white">{selected.model}</span></div>
                <div><span className="text-gray-500">Duration:</span> <span className="text-gray-900 dark:text-white">{selected.duration_ms}ms</span></div>
                <div><span className="text-gray-500">Tokens:</span> <span className="text-gray-900 dark:text-white">{selected.token_count}</span></div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
