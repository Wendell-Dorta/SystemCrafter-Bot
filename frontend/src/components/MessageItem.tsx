'use client';

import React, { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Message } from '@/types';
import { MermaidViewer } from './MermaidViewer';
import { ToolExecutionCard } from './ToolExecutionCard';
import { Bot, User, Copy, Check } from 'lucide-react';
import { formatDate } from '@/lib/utils';

interface MessageItemProps {
  message: Message;
  isLastStreaming?: boolean;
}

export const MessageItem: React.FC<MessageItemProps> = ({ message, isLastStreaming }) => {
  const isUser = message.role === 'user';
  const [copiedCodeIndex, setCopiedCodeIndex] = useState<number | null>(null);

  const handleCopy = (text: string, idx: number) => {
    navigator.clipboard.writeText(text);
    setCopiedCodeIndex(idx);
    setTimeout(() => setCopiedCodeIndex(null), 2000);
  };

  return (
    <div className={`py-4 sm:py-6 w-full ${isUser ? 'bg-transparent' : 'bg-slate-900/30 border-y border-slate-800/40'}`}>
      <div className="mx-auto max-w-4xl lg:max-w-5xl px-3 sm:px-6 flex gap-2.5 sm:gap-4">
        {/* Avatar */}
        <div className="flex-shrink-0 pt-0.5">
          {isUser ? (
            <div className="flex h-7 w-7 sm:h-8 sm:w-8 items-center justify-center rounded-xl bg-slate-800 border border-slate-700 text-slate-300 shadow-sm">
              <User className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
            </div>
          ) : (
            <div className="flex h-7 w-7 sm:h-8 sm:w-8 items-center justify-center rounded-xl bg-gradient-to-tr from-sky-500 to-blue-600 text-white shadow-md shadow-sky-500/20">
              <Bot className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
            </div>
          )}
        </div>

        {/* Content Body */}
        <div className="flex-1 min-w-0 overflow-hidden space-y-2.5 sm:space-y-3">
          {/* Header Role & Time */}
          <div className="flex items-center gap-2 text-xs">
            <span className="font-semibold text-slate-200">
              {isUser ? 'Você' : 'ArchMind AI'}
            </span>
            <span className="text-[9px] sm:text-[10px] text-slate-500 font-mono">
              {formatDate(message.createdAt)}
            </span>
          </div>

          {/* User Attached Image (Multimodal Diagram) */}
          {message.imageBase64 && (
            <div className="inline-block rounded-xl border border-sky-500/30 overflow-hidden shadow-lg bg-slate-950/60 p-1 max-w-full">
              <img
                src={`data:${message.imageMimeType || 'image/png'};base64,${message.imageBase64}`}
                alt="Diagrama Anexado"
                className="max-h-56 sm:max-h-64 rounded-lg object-contain max-w-full"
              />
            </div>
          )}

          {/* Tool Execution Cards (Agentic Loop Results) */}
          {message.toolEvents && message.toolEvents.length > 0 && (
            <div className="space-y-2 w-full overflow-hidden">
              {message.toolEvents.map((toolEvt, idx) => (
                <ToolExecutionCard key={idx} toolEvent={toolEvt} />
              ))}
            </div>
          )}

          {/* Message Text / Markdown */}
          {isUser ? (
            <p className="text-xs sm:text-sm text-slate-200 whitespace-pre-wrap leading-relaxed break-words">
              {message.content}
            </p>
          ) : (
            <div className="text-xs sm:text-sm text-slate-200 leading-relaxed prose prose-invert max-w-none overflow-hidden break-words">
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                  code({ node, className, children, ...props }: any) {
                    const match = /language-(\w+)/.exec(className || '');
                    const lang = match ? match[1] : '';
                    const codeString = String(children).replace(/\n$/, '');

                    // Intercept Mermaid blocks and render interactive vector diagram
                    if (lang === 'mermaid') {
                      return <MermaidViewer chart={codeString} isStreaming={isLastStreaming} />;
                    }

                    // Multi-line code block with copy button
                    if (match) {
                      const codeId = Math.random();
                      return (
                        <div className="my-3 overflow-hidden rounded-xl border border-slate-800 bg-slate-950 shadow-md max-w-full">
                          <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/90 px-3 sm:px-4 py-1.5 text-xs text-slate-400 font-mono">
                            <span className="text-[11px] font-semibold text-sky-400">{lang}</span>
                            <button
                              onClick={() => handleCopy(codeString, codeId as any)}
                              className="flex items-center gap-1 hover:text-slate-200 transition-colors cursor-pointer text-[11px]"
                            >
                              {copiedCodeIndex === (codeId as any) ? (
                                <>
                                  <Check className="h-3 w-3 text-emerald-400" />
                                  <span className="text-emerald-400">Copiado!</span>
                                </>
                              ) : (
                                <>
                                  <Copy className="h-3 w-3" />
                                  <span>Copiar</span>
                                </>
                              )}
                            </button>
                          </div>
                          <pre className="overflow-x-auto p-3 sm:p-4 text-xs font-mono text-slate-300 leading-relaxed">
                            <code>{children}</code>
                          </pre>
                        </div>
                      );
                    }

                    // Inline code
                    return (
                      <code className="rounded bg-slate-800/80 px-1.5 py-0.5 font-mono text-[11px] text-sky-300 border border-slate-700/50 break-all">
                        {children}
                      </code>
                    );
                  },
                  table({ children }) {
                    return (
                      <div className="my-3 sm:my-4 overflow-x-auto rounded-xl border border-slate-800 bg-slate-950/60 shadow-sm max-w-full">
                        <table className="w-full text-left text-xs border-collapse min-w-[320px]">{children}</table>
                      </div>
                    );
                  },
                  th({ children }) {
                    return (
                      <th className="border-b border-slate-800 bg-slate-900/80 px-3 sm:px-4 py-2 font-semibold text-slate-200 whitespace-nowrap">
                        {children}
                      </th>
                    );
                  },
                  td({ children }) {
                    return (
                      <td className="border-b border-slate-800/60 px-3 sm:px-4 py-2 text-slate-300">
                        {children}
                      </td>
                    );
                  },
                  ul({ children }) {
                    return <ul className="my-2 space-y-1 list-disc list-inside text-slate-300">{children}</ul>;
                  },
                  ol({ children }) {
                    return <ol className="my-2 space-y-1 list-decimal list-inside text-slate-300">{children}</ol>;
                  },
                  h2({ children }) {
                    return <h2 className="text-sm sm:text-base font-bold text-slate-100 mt-4 mb-2 pb-1 border-b border-slate-800">{children}</h2>;
                  },
                  h3({ children }) {
                    return <h3 className="text-xs sm:text-sm font-bold text-slate-200 mt-3 mb-1.5">{children}</h3>;
                  },
                }}
              >
                {message.content}
              </ReactMarkdown>

              {/* Streaming Blinking Cursor */}
              {isLastStreaming && <span className="cursor-blink" />}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
