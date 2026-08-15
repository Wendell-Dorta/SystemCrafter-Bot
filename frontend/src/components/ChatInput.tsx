'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Send, Image as ImageIcon, Square, X, Sparkles } from 'lucide-react';

interface ChatInputProps {
  onSendMessage: (message: string, imageBase64?: string, imageMimeType?: string) => void;
  onStopStream: () => void;
  isStreaming: boolean;
  disabled?: boolean;
}

export const ChatInput: React.FC<ChatInputProps> = ({
  onSendMessage,
  onStopStream,
  isStreaming,
  disabled,
}) => {
  const [text, setText] = useState<string>('');
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [imageMimeType, setImageMimeType] = useState<string>('image/png');
  const [imageBase64Raw, setImageBase64Raw] = useState<string>('');

  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Auto-resize textarea smoothly up to 160px max
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 160)}px`;
    }
  }, [text]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Use HTML5 Canvas to downscale high-res camera photos/screenshots to max 1280px
    // This prevents token volume blowup, eliminates 429 RESOURCE_EXHAUSTED errors, and makes uploads instant
    const reader = new FileReader();
    reader.onload = (event) => {
      const originalDataUrl = event.target?.result as string;
      const img = new Image();
      img.onload = () => {
        const maxDim = 1280;
        let width = img.width;
        let height = img.height;

        if (width > maxDim || height > maxDim) {
          if (width > height) {
            height = Math.round((height * maxDim) / width);
            width = maxDim;
          } else {
            width = Math.round((width * maxDim) / height);
            height = maxDim;
          }
        }

        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext('2d');
        if (ctx) {
          ctx.drawImage(img, 0, 0, width, height);
          const mime = 'image/jpeg';
          const compressedDataUrl = canvas.toDataURL(mime, 0.85);
          setImagePreview(compressedDataUrl);
          setImageMimeType(mime);
          const rawBase64 = compressedDataUrl.split(',')[1] || compressedDataUrl;
          setImageBase64Raw(rawBase64);
        } else {
          // Fallback if canvas is unavailable
          setImagePreview(originalDataUrl);
          setImageMimeType(file.type || 'image/png');
          const rawBase64 = originalDataUrl.split(',')[1] || originalDataUrl;
          setImageBase64Raw(rawBase64);
        }
      };
      img.onerror = () => {
        setImagePreview(originalDataUrl);
        setImageMimeType(file.type || 'image/png');
        const rawBase64 = originalDataUrl.split(',')[1] || originalDataUrl;
        setImageBase64Raw(rawBase64);
      };
      img.src = originalDataUrl;
    };
    reader.readAsDataURL(file);
  };

  const removeImage = () => {
    setImagePreview(null);
    setImageBase64Raw('');
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = () => {
    if ((!text.trim() && !imageBase64Raw) || isStreaming || disabled) return;

    onSendMessage(text.trim(), imageBase64Raw || undefined, imageMimeType || undefined);
    setText('');
    removeImage();

    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
    }
  };

  const quickActions = [
    { label: '💰 Custos AWS vs OCI', query: 'Estimar custos comparativos entre AWS e Oracle OCI para 15M requisições/mês com Kubernetes e PostgreSQL' },
    { label: '🛡️ SPOF & LGPD', query: 'Auditar conformidade LGPD e pontos únicos de falha nesta arquitetura de microsserviços' },
    { label: '📐 Outbox Pattern', query: 'Gerar diagrama Mermaid e código Go do padrão Transactional Outbox com Kafka' },
  ];

  return (
    <div className="relative border-t border-slate-800/80 bg-slate-950/90 p-2.5 sm:p-4 backdrop-blur-xl flex-shrink-0">
      <div className="mx-auto max-w-4xl space-y-2">
        {/* Quick Action Chips (Mobile-First touch scroll) */}
        <div className="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs no-scrollbar touch-pan-x">
          <span className="text-[10px] font-mono text-slate-500 flex items-center gap-1 flex-shrink-0 hidden xs:flex">
            <Sparkles className="h-3 w-3 text-sky-400" /> Ações:
          </span>
          {quickActions.map((qa, i) => (
            <button
              key={i}
              type="button"
              disabled={isStreaming}
              onClick={() => setText(qa.query)}
              className="inline-flex flex-shrink-0 items-center gap-1 rounded-full border border-slate-800 bg-slate-900/70 px-2.5 py-1 text-[10px] sm:text-[11px] text-slate-300 hover:border-sky-500/40 hover:bg-slate-800 hover:text-white transition-colors cursor-pointer disabled:opacity-50 active:scale-95"
            >
              {qa.label}
            </button>
          ))}
        </div>

        {/* Image Attachment Preview */}
        {imagePreview && (
          <div className="relative inline-flex items-center gap-2 rounded-xl border border-sky-500/40 bg-sky-950/40 p-2 shadow-md">
            <img
              src={imagePreview}
              alt="Preview"
              className="h-12 w-12 sm:h-14 sm:w-14 rounded-lg object-cover border border-sky-500/30 flex-shrink-0"
            />
            <div className="pr-6 text-xs overflow-hidden">
              <span className="font-semibold text-slate-200 block text-[11px] truncate">
                Rascunho de Arquitetura Anexado
              </span>
              <span className="text-[10px] text-sky-400 block truncate">
                Visão multimodal ativa
              </span>
            </div>
            <button
              onClick={removeImage}
              className="absolute top-1.5 right-1.5 rounded-full bg-slate-800 p-1 text-slate-400 hover:bg-slate-700 hover:text-white cursor-pointer"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        )}

        {/* Input Text Box */}
        <div className="relative flex items-end gap-1.5 sm:gap-2 rounded-2xl border border-slate-700/80 bg-slate-900/95 p-1.5 sm:p-2 shadow-inner focus-within:border-sky-500/60 focus-within:ring-1 focus-within:ring-sky-500/30 transition-all">
          <input
            type="file"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept="image/png,image/jpeg,image/webp,image/jpg"
            className="hidden"
          />

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            title="Anexar foto ou print de arquitetura"
            className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-xl text-slate-400 hover:bg-slate-800 hover:text-sky-400 transition-colors flex-shrink-0 cursor-pointer"
          >
            <ImageIcon className="h-4 w-4 sm:h-5 sm:w-5" />
          </button>

          <textarea
            ref={textareaRef}
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            rows={1}
            disabled={disabled}
            placeholder="Descreva seu sistema ou faça perguntas de arquitetura..."
            className="flex-1 min-h-[36px] max-h-36 resize-none bg-transparent py-2 px-1 text-xs sm:text-sm text-slate-100 placeholder-slate-500 focus:outline-none leading-relaxed"
          />

          {isStreaming ? (
            <button
              type="button"
              onClick={onStopStream}
              title="Interromper geração"
              className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-xl bg-amber-500/20 text-amber-400 hover:bg-amber-500/30 border border-amber-500/40 transition-colors flex-shrink-0 cursor-pointer"
            >
              <Square className="h-4 w-4 fill-current" />
            </button>
          ) : (
            <button
              type="button"
              onClick={handleSubmit}
              disabled={(!text.trim() && !imageBase64Raw) || disabled}
              title="Enviar mensagem (Enter)"
              className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-xl bg-gradient-to-r from-sky-500 to-blue-600 text-white shadow-md shadow-sky-500/20 hover:from-sky-400 hover:to-blue-500 transition-all disabled:opacity-40 disabled:cursor-not-allowed flex-shrink-0 cursor-pointer"
            >
              <Send className="h-4 w-4" />
            </button>
          )}
        </div>

        <div className="flex items-center justify-between px-1 text-[9px] sm:text-[10px] text-slate-500 font-mono">
          <span className="truncate">Enter envia • Shift+Enter nova linha</span>
          <span className="hidden sm:inline">Gemini Flash • Go Clean Arch</span>
        </div>
      </div>
    </div>
  );
};
