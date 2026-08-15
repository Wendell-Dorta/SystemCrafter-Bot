'use client';

import React, { useEffect, useRef } from 'react';
import { Message } from '@/types';
import { MessageItem } from './MessageItem';
import { Loader2, AlertCircle } from 'lucide-react';

interface MessageListProps {
  messages: Message[];
  isStreaming: boolean;
  currentTool: { toolName: string; status: string } | null;
  error: string | null;
}

export const MessageList: React.FC<MessageListProps> = ({
  messages,
  isStreaming,
  currentTool,
  error,
}) => {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, currentTool, isStreaming]);

  return (
    <div className="flex-1 overflow-y-auto overscroll-contain pb-6 w-full">
      <div className="w-full">
        {messages.map((msg, idx) => {
          const isLast = idx === messages.length - 1;
          const isLastStreaming = isLast && isStreaming && msg.role === 'model';
          return <MessageItem key={msg.id || idx} message={msg} isLastStreaming={isLastStreaming} />;
        })}
      </div>

      {/* Active Tool Execution Indicator */}
      {currentTool && (
        <div className="mx-auto max-w-4xl lg:max-w-5xl px-3 sm:px-6 py-3">
          <div className="flex items-center gap-3 rounded-2xl border border-sky-500/40 bg-sky-950/30 p-3.5 sm:p-4 shadow-lg backdrop-blur-md animate-pulse">
            <div className="flex h-8 w-8 items-center justify-center rounded-xl bg-sky-500/20 text-sky-400 flex-shrink-0">
              <Loader2 className="h-4 w-4 animate-spin" />
            </div>
            <div className="overflow-hidden">
              <span className="font-semibold text-xs sm:text-sm text-sky-300 block truncate">
                Executando {currentTool.toolName}...
              </span>
              <span className="text-[10px] sm:text-xs text-slate-400 block truncate">{currentTool.status}</span>
            </div>
          </div>
        </div>
      )}

      {/* Error Alert */}
      {error && (
        <div className="mx-auto max-w-4xl lg:max-w-5xl px-3 sm:px-6 py-3">
          <div className="flex items-start gap-3 rounded-2xl border border-red-500/30 bg-red-950/30 p-3.5 sm:p-4 text-xs text-red-300 shadow-md">
            <AlertCircle className="h-4 w-4 text-red-400 flex-shrink-0 mt-0.5" />
            <div className="overflow-hidden">
              <span className="font-semibold block mb-0.5">Ocorreu um erro:</span>
              <p className="text-slate-400 break-words">{error}</p>
            </div>
          </div>
        </div>
      )}

      <div ref={bottomRef} className="h-2" />
    </div>
  );
};
