"use client";

import { useEffect, useRef } from "react";

import { ChatInput } from "@/components/chat/ChatInput";
import { ChatMessage } from "@/components/chat/ChatMessage";
import { useAppStore } from "@/lib/store";

export function ChatPanel() {
  const { messages, sendMessage, isStreaming, tokenStats } = useAppStore();
  const endRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <section className="flex h-full min-w-0 flex-1 flex-col gap-6">
      <div className="panel flex items-center justify-between rounded-[24px] px-6 py-4">
        <div className="flex-1">
          <p className="text-xs uppercase tracking-wider text-[var(--color-ink-soft)]">
            Conversation
          </p>
          <h1 className="text-2xl font-bold tracking-tight">
            实时对话与工具回放
          </h1>
        </div>
        <div className="mono text-sm text-[var(--color-ink-soft)] flex-shrink-0">
          {tokenStats ? `${tokenStats.total_tokens} tokens` : "No metrics yet"}
        </div>
      </div>

      <div className="flex-1 flex flex-col">
        <div className="flex-1 overflow-y-auto pr-4">
          {!messages.length && (
            <div className="flex flex-col items-center justify-center h-full px-4">
              <div className="rounded-[24px] border border-dashed border-[var(--color-line)] bg-white/30 p-12 text-center">
                <p className="mb-2 text-xs uppercase tracking-wider text-[var(--color-ink-soft)]">
                  Ready to chat
                </p>
                <h2 className="mt-4 text-xl font-semibold tracking-tight text-[var(--color-ink)]">
                  一个本地、透明、文件驱动的 Agent 工作台
                </h2>
                <p className="mt-6 max-w-xl text-[var(--color-ink-soft)] leading-relaxed">
                  你可以直接提问，也可以在右侧编辑 Memory、Skills 和 Workspace 文件。
                  所有系统提示、会话和工具执行都可追踪。
                </p>
              </div>
            </div>
          )}

          {messages.map((message) => (
            <ChatMessage
              content={message.content}
              key={message.id}
              retrievals={message.retrievals}
              role={message.role}
              toolCalls={message.toolCalls}
            />
          ))}
          <div ref={endRef} />
        </div>
      </div>

      <ChatInput disabled={isStreaming} onSend={sendMessage} />
    </section>
  );
}
