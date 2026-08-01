"use client";

import { useEffect, useRef, useState } from "react";
import { chatApi } from "@/lib/api";

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
  created_at: string;
  updated_at: string;
}

export default function CoachPage() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    chatApi
      .listConversations()
      .then(setConversations)
      .catch((e) => setError(e.message));
  }, []);

  useEffect(() => {
    if (!activeId) return;
    chatApi
      .getConversation(activeId)
      .then((data) => {
        setMessages(data.messages || []);
      })
      .catch((e) => setError(e.message));
  }, [activeId]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, loading]);

  const startNewChat = async () => {
    setError("");
    try {
      const conv = await chatApi.createConversation();
      setConversations((prev) => [conv, ...prev]);
      setActiveId(conv.id);
      setMessages([]);
    } catch (e: any) {
      setError(e.message);
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
      setError(e.message);
      setMessages((prev) => [
        ...prev,
        { id: `err-${Date.now()}`, role: "assistant", content: `⚠️ ${e.message}`, created_at: new Date().toISOString() },
      ]);
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
            onClick={startNewChat}
            className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-medium hover:opacity-90 transition-opacity"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Yeni Sohbet
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1">
          {conversations.length === 0 && (
            <p className="text-sm text-gray-400 text-center mt-8 px-4">
              Henüz sohbet yok. CV'n hakkında soru sorarak başla.
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
            CV'ni iyileştir, mülakatlara hazırlan, başvurularını güçlendir.
          </p>
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
                CV'nin özet bölümünü nasıl güçlendiririm? Hangi becerileri öne çıkarmalıyım?
                Mülakatta "güçlü yönlerin neler?" sorusuna nasıl cevap vermeliyim? gibi sorular
                sorabilirsin. CV'lerini yüklediysen kişisel öneriler veririm.
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
