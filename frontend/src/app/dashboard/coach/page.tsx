"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { api, chatApi, friendlyError } from "@/lib/api";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  created_at: string;
}

interface Conversation {
  id: string;
  title: string;
  cv_id: string | null;
  jd_id: string | null;
  match_id: string | null;
  created_at: string;
  updated_at: string;
}

interface CVItem {
  id: string;
  title: string;
}

interface JDItem {
  id: string;
  title: string;
}

interface MatchItem {
  id: string;
  cv_title: string;
  jd_title: string;
  overall_score: number;
}

export default function CoachPage() {
  const router = useRouter();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [showContext, setShowContext] = useState(false);
  const [creating, setCreating] = useState(false);
  const [cvs, setCvs] = useState<CVItem[]>([]);
  const [jds, setJds] = useState<JDItem[]>([]);
  const [matches, setMatches] = useState<MatchItem[]>([]);
  const [selCv, setSelCv] = useState("");
  const [selJd, setSelJd] = useState("");
  const [selMatch, setSelMatch] = useState("");
  const [contextLabel, setContextLabel] = useState("");

  const bottomRef = useRef<HTMLDivElement>(null);
  const initializedRef = useRef(false);

  useEffect(() => {
    chatApi
      .listConversations()
      .then(setConversations)
      .catch((e) => setError(e.message));
    Promise.all([
      api.get("/api/v1/cvs").then((d) => setCvs(d.items || [])),
      api.get("/api/v1/jds").then((d) => setJds(d.items || [])),
      api.get("/api/v1/matches").then((d) => setMatches(d.items || [])),
    ]).catch(() => {});
  }, []);

  useEffect(() => {
    const convParam = new URLSearchParams(window.location.search).get("conv");
    if (convParam && !initializedRef.current) {
      initializedRef.current = true;
      setActiveId(convParam);
    }
  }, []);

  useEffect(() => {
    if (!activeId) return;
    chatApi
      .getConversation(activeId)
      .then((data) => {
        setMessages(data.messages || []);
        const conv: Conversation | undefined = data.conversation;
        if (conv) {
          setContextLabel(buildContextLabel(conv, cvs, jds, matches));
          setConversations((prev) => {
            if (prev.some((c) => c.id === conv.id)) {
              return prev.map((c) => (c.id === conv.id ? conv : c));
            }
            return [conv, ...prev];
          });
        }
      })
      .catch((e) => setError(e.message));
  }, [activeId, cvs, jds, matches]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, loading]);

  const openContextPicker = () => {
    setError("");
    setSelCv("");
    setSelJd("");
    setSelMatch("");
    setShowContext(true);
  };

  const startNewChat = async () => {
    if (creating) return;
    setError("");
    setCreating(true);
    try {
      const conv = selMatch
        ? await chatApi.createConversation(undefined, undefined, undefined, selMatch)
        : await chatApi.createConversation(undefined, selCv || undefined, selJd || undefined);
      setConversations((prev) => [conv, ...prev]);
      setActiveId(conv.id);
      setMessages([]);
      setContextLabel(selMatch ? "Match bağlamı eklendi" : [selCv ? "CV" : "", selJd ? "JD" : ""].filter(Boolean).join(" + "));
      setShowContext(false);
      router.replace("/dashboard/coach");
    } catch (e: any) {
      setError(e.message);
    } finally {
      setCreating(false);
    }
  };

  const openConversation = (id: string) => {
    setActiveId(id);
  };

  const deleteConversation = async (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    try {
      await chatApi.deleteConversation(id);
      setConversations((prev) => prev.filter((c) => c.id !== id));
      if (activeId === id) {
        setActiveId(null);
        setMessages([]);
        setContextLabel("");
      }
    } catch (err: any) {
      setError(err.message);
    }
  };

  const sendMessage = async () => {
    const content = input.trim();
    if (!content || !activeId || loading) return;
    setError("");
    setInput("");
    setMessages((prev) => [
      ...prev,
      { id: `tmp-${Date.now()}`, role: "user", content, created_at: new Date().toISOString() },
    ]);
    setLoading(true);
    try {
      const result = await chatApi.sendMessage(activeId, content);
      setMessages((prev) => [...prev, result.assistant_message]);
      if (result.conversation) {
        setConversations((prev) =>
          prev.map((c) => (c.id === result.conversation.id ? result.conversation : c))
        );
      }
    } catch (e: any) {
      setInput(content);
      setError(friendlyError(e).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex h-full">
      {/* Conversation list */}
      <div className="w-72 border-r border-gray-200 dark:border-gray-800 flex flex-col">
        <div className="p-4 border-b border-gray-200 dark:border-gray-800">
          <button
            onClick={openContextPicker}
            className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-opacity"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Yeni Sohbet
          </button>
        </div>

        {showContext && (
          <div className="p-4 border-b border-gray-200 dark:border-gray-800 space-y-3">
            <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">
              Bağlam Ekle (isteğe bağlı)
            </p>
            <div>
              <label className="block text-xs text-gray-500 mb-1">Match</label>
              <select
                value={selMatch}
                onChange={(e) => {
                  setSelMatch(e.target.value);
                  if (e.target.value) {
                    setSelCv("");
                    setSelJd("");
                  }
                }}
                className="w-full px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-purple-500"
              >
                <option value="">— Match seç —</option>
                {matches.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.cv_title} ↔ {m.jd_title} ({(m.overall_score * 100).toFixed(0)})
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">CV</label>
              <select
                value={selCv}
                onChange={(e) => {
                  setSelCv(e.target.value);
                  if (e.target.value) setSelMatch("");
                }}
                disabled={!!selMatch}
                className="w-full px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
              >
                <option value="">— CV seç —</option>
                {cvs.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.title}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">Job Description</label>
              <select
                value={selJd}
                onChange={(e) => {
                  setSelJd(e.target.value);
                  if (e.target.value) setSelMatch("");
                }}
                disabled={!!selMatch}
                className="w-full px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
              >
                <option value="">— JD seç —</option>
                {jds.map((j) => (
                  <option key={j.id} value={j.id}>
                    {j.title}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex gap-2 pt-1">
              <button
                onClick={startNewChat}
                disabled={creating}
                className="flex-1 px-3 py-2 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
              >
                {creating ? "Oluşturuluyor..." : "Sohbeti Başlat"}
              </button>
              <button
                onClick={() => setShowContext(false)}
                className="px-3 py-2 rounded-lg bg-white/10 text-gray-400 hover:bg-white/20 text-sm"
              >
                Vazgeç
              </button>
            </div>
          </div>
        )}

        <div className="flex-1 overflow-y-auto p-2 space-y-1">
          {conversations.length === 0 && !showContext && (
            <p className="text-sm text-gray-400 text-center mt-8 px-4">
              Henüz sohbet yok. "Yeni Sohbet" ile CV'n, bir JD veya match skoru üzerine konuşmaya başla.
            </p>
          )}
          {conversations.map((conv) => (
            <div
              key={conv.id}
              onClick={() => openConversation(conv.id)}
              className={`group flex items-center justify-between gap-2 px-3 py-2.5 rounded-xl cursor-pointer text-sm transition-colors ${
                activeId === conv.id
                  ? "bg-purple-500/10 text-purple-600 dark:text-purple-300 border border-purple-500/30"
                  : "text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800/60"
              }`}
            >
              <span className="truncate">{conv.title}</span>
              <button
                onClick={(e) => deleteConversation(e, conv.id)}
                className="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 transition-opacity"
                title="Sil"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Chat area */}
      <div className="flex-1 flex flex-col min-w-0">
        <div className="p-4 border-b border-gray-200 dark:border-gray-800">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">CV Coach</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            CV'ni iyileştir, mülakatlara hazırlan, match skorunu yükselt.
          </p>
          {contextLabel && (
            <span className="inline-flex items-center gap-1.5 mt-2 px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/30 text-xs text-purple-600 dark:text-purple-300">
              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
              </svg>
              {contextLabel}
            </span>
          )}
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-4">
          {messages.length === 0 && !loading && (
            <div className="h-full flex flex-col items-center justify-center text-center px-8">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center mb-4">
                <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
              </div>
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                Merhaba! Ben CV Coach
              </h2>
              <p className="text-sm text-gray-500 dark:text-gray-400 max-w-md">
                CV'nin özet bölümünü nasıl güçlendiririm? Match skorumu nasıl artırırım? Eksik
                becerilerimi nasıl kapatırım? gibi sorular sorabilirsin. Bağlam eklediysen
                (CV / JD / Match) kişisel ve skor bazlı öneriler veririm.
              </p>
            </div>
          )}

          {messages.map((m) => (
            <div key={m.id} className={`flex ${m.role === "user" ? "justify-end" : "justify-start"}`}>
              <div
                className={`max-w-[80%] px-4 py-3 rounded-2xl whitespace-pre-wrap text-sm leading-relaxed ${
                  m.role === "user"
                    ? "bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-br-md"
                    : "bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 border border-gray-200 dark:border-gray-700 rounded-bl-md"
                }`}
              >
                {m.content}
              </div>
            </div>
          ))}

          {loading && (
            <div className="flex justify-start">
              <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl rounded-bl-md px-4 py-3">
                <div className="flex gap-1.5">
                  <span className="w-2 h-2 rounded-full bg-purple-400 animate-bounce" />
                  <span className="w-2 h-2 rounded-full bg-purple-400 animate-bounce [animation-delay:0.15s]" />
                  <span className="w-2 h-2 rounded-full bg-purple-400 animate-bounce [animation-delay:0.3s]" />
                </div>
              </div>
            </div>
          )}

          {error && (
            <div className="text-center">
              <p className="text-sm text-red-500 inline-block bg-red-500/10 border border-red-500/20 rounded-xl px-4 py-2">
                {error}
              </p>
            </div>
          )}

          <div ref={bottomRef} />
        </div>

        <div className="p-4 border-t border-gray-200 dark:border-gray-800">
          <div className="flex gap-3">
            <input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && sendMessage()}
              placeholder={activeId ? "Mesajını yaz..." : "Önce yeni bir sohbet başlat"}
              disabled={!activeId || loading}
              className="flex-1 px-4 py-3 rounded-xl bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
            />
            <button
              onClick={sendMessage}
              disabled={!activeId || !input.trim() || loading}
              className="px-5 py-3 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-opacity disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Gönder
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

function buildContextLabel(
  conv: Conversation,
  cvs: CVItem[],
  jds: JDItem[],
  matches: MatchItem[]
): string {
  if (conv.match_id) {
    const m = matches.find((x) => x.id === conv.match_id);
    if (m) return `Match: ${m.cv_title} ↔ ${m.jd_title}`;
    return "Match";
  }
  const parts: string[] = [];
  if (conv.cv_id) {
    const c = cvs.find((x) => x.id === conv.cv_id);
    parts.push(`CV: ${c?.title ?? ""}`);
  }
  if (conv.jd_id) {
    const j = jds.find((x) => x.id === conv.jd_id);
    parts.push(`JD: ${j?.title ?? ""}`);
  }
  return parts.join(" + ");
}
