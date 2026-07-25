export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <div>
      <div className="flex gap-2 mb-6">
        <a href="/dashboard/admin/adapters" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Adapters</a>
        <a href="/dashboard/admin/prompts" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Prompts</a>
        <a href="/dashboard/admin/settings" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Settings</a>
        <a href="/dashboard/admin/logs" className="text-sm px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-purple-500/10 hover:text-purple-400">Logs</a>
      </div>
      {children}
    </div>
  );
}
