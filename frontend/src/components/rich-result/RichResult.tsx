"use client";

interface RichResultSection {
  title: string;
  content: string;
  type: string;
}

interface RichResultData {
  text: string;
  type: string;
  sections?: RichResultSection[];
  metadata?: Record<string, unknown>;
}

export default function RichResult({ result }: { result: RichResultData }) {
  if (result.sections && result.sections.length > 0) {
    return (
      <div className="space-y-6">
        {result.sections.map((section, i) => (
          <div key={i} className="bg-white dark:bg-slate-900/50 rounded-xl border border-gray-200 dark:border-white/10 p-5">
            {section.title && (
              <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2 uppercase tracking-wider">{section.title}</h3>
            )}
            <div className="prose prose-sm dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
              {section.content}
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (result.type === "code" && result.text.includes("```")) {
    const parts = result.text.split("```");
    return (
      <div className="space-y-4">
        {parts.map((part, i) => {
          if (i % 2 === 0) {
            return part.trim() ? (
              <p key={i} className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{part}</p>
            ) : null;
          }
          const lines = part.split("\n");
          const lang = lines[0] || "";
          const code = lines.slice(1).join("\n");
          return (
            <pre key={i} className="bg-gray-900 text-green-400 rounded-xl p-4 overflow-x-auto text-sm">
              <code>{code || part}</code>
            </pre>
          );
        })}
      </div>
    );
  }

  return (
    <div className="prose prose-sm dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
      {result.text}
    </div>
  );
}
