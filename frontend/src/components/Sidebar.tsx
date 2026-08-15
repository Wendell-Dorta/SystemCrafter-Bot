'use client';

import React, { useState, useEffect } from 'react';
import { Session, ArchitectureTemplate } from '@/types';
import { fetchSessions, fetchTemplates, deleteSession } from '@/lib/api';
import {
  MessageSquare,
  Layers,
  Plus,
  X,
  Brain,
  Shield,
  Trash2,
  Cpu,
  History,
  ArrowRight,
  Search,
} from 'lucide-react';
import { formatDate } from '@/lib/utils';

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
  currentSessionId: string;
  onSelectSession: (sessionId: string) => void;
  onNewSession: () => void;
  onSelectTemplate: (template: ArchitectureTemplate) => void;
  onOpenHistoryModal: () => void;
  onOpenTemplatesModal: () => void;
  onOpenProfileModal: () => void;
  onOpenCodeAuditModal: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  isOpen,
  onClose,
  currentSessionId,
  onSelectSession,
  onNewSession,
  onSelectTemplate,
  onOpenHistoryModal,
  onOpenTemplatesModal,
  onOpenProfileModal,
  onOpenCodeAuditModal,
}) => {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [templates, setTemplates] = useState<ArchitectureTemplate[]>([]);

  const loadData = () => {
    fetchSessions().then((list) => setSessions(list || []));
    fetchTemplates().then((list) => setTemplates(list || []));
  };

  useEffect(() => {
    if (isOpen) {
      loadData();
    }
  }, [isOpen, currentSessionId]);

  const handleDeleteSession = async (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (confirm('Deseja excluir esta conversa do histórico?')) {
      const ok = await deleteSession(id);
      if (ok) {
        setSessions((prev) => prev.filter((s) => s.id !== id));
        if (id === currentSessionId) {
          onNewSession();
        }
      }
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Mobile backdrop */}
      <div
        onClick={onClose}
        className="fixed inset-0 z-40 bg-black/70 backdrop-blur-sm lg:hidden animate-in fade-in"
      />

      {/* Sidebar panel */}
      <aside className="fixed top-0 left-0 z-50 flex h-full w-80 flex-col border-r border-slate-800 bg-slate-950 p-4 shadow-2xl transition-transform animate-in slide-in-from-left duration-200">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800/80 pb-3">
          <div className="flex items-center gap-2">
            <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-sky-500/10 text-sky-400 border border-sky-500/20">
              <Cpu className="h-4 w-4" />
            </div>
            <span className="font-bold text-xs text-slate-200">Menu Principal</span>
          </div>
          <button
            onClick={onClose}
            className="rounded-lg p-1 text-slate-400 hover:bg-slate-800 hover:text-white cursor-pointer"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Quick Hub Action Buttons */}
        <div className="py-3 space-y-2">
          <button
            onClick={() => {
              onNewSession();
              onClose();
            }}
            className="w-full flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-sky-500 to-blue-600 px-4 py-2 text-xs font-semibold text-white shadow-md shadow-sky-500/20 hover:from-sky-400 hover:to-blue-500 transition-all cursor-pointer"
          >
            <Plus className="h-4 w-4" />
            Novo Chat / Sessão
          </button>

          <button
            onClick={() => {
              onOpenHistoryModal();
              onClose();
            }}
            className="w-full flex items-center justify-between rounded-xl bg-slate-900 border border-sky-500/30 px-3.5 py-2 text-xs font-medium text-sky-300 hover:bg-sky-500/10 transition-all cursor-pointer"
          >
            <div className="flex items-center gap-2">
              <History className="h-4 w-4 text-sky-400" />
              <span>Histórico de Chats</span>
            </div>
            <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-slate-800 text-sky-300 border border-slate-700">
              {sessions.length}
            </span>
          </button>

          <button
            onClick={() => {
              onOpenCodeAuditModal();
              onClose();
            }}
            className="w-full flex items-center justify-start gap-2 rounded-xl bg-slate-900 border border-sky-500/30 px-3.5 py-2 text-xs font-medium text-sky-300 hover:bg-sky-500/10 transition-all cursor-pointer"
          >
            <Shield className="h-4 w-4 text-sky-400" />
            <span>Auditar Repositório GitHub</span>
          </button>

          <button
            onClick={() => {
              onOpenProfileModal();
              onClose();
            }}
            className="w-full flex items-center justify-start gap-2 rounded-xl bg-slate-900 border border-amber-500/30 px-3.5 py-2 text-xs font-medium text-amber-300 hover:bg-amber-500/10 transition-all cursor-pointer"
          >
            <Brain className="h-4 w-4 text-amber-400" />
            <span>Memória do Arquiteto</span>
          </button>
        </div>

        {/* Scrollable sections */}
        <div className="flex-1 overflow-y-auto space-y-4 pr-1">
          {/* Section: Recent Sessions */}
          <div>
            <div className="flex items-center justify-between px-1 mb-2">
              <span className="text-[10px] font-mono uppercase tracking-wider text-slate-400 font-bold">
                Conversas Recentes
              </span>
              <button
                onClick={() => {
                  onOpenHistoryModal();
                  onClose();
                }}
                className="text-[10px] text-sky-400 hover:underline flex items-center gap-0.5 cursor-pointer"
              >
                Buscar todas <ArrowRight className="w-2.5 h-2.5" />
              </button>
            </div>

            {sessions.length === 0 ? (
              <div className="rounded-xl border border-slate-800/60 bg-slate-900/20 p-3 text-center text-xs text-slate-500">
                Nenhuma conversa gravada
              </div>
            ) : (
              <div className="space-y-1">
                {sessions.slice(0, 3).map((sess) => {
                  const isActive = sess.id === currentSessionId;
                  return (
                    <div
                      key={sess.id}
                      onClick={() => {
                        onSelectSession(sess.id);
                        onClose();
                      }}
                      className={`group w-full text-left p-2.5 rounded-xl border transition-all flex items-center justify-between gap-2 cursor-pointer ${
                        isActive
                          ? 'border-sky-500/50 bg-sky-950/30 text-slate-100 shadow-sm'
                          : 'border-transparent text-slate-400 hover:border-slate-800 hover:bg-slate-900/50 hover:text-slate-200'
                      }`}
                    >
                      <div className="flex items-center gap-2 overflow-hidden flex-1">
                        <MessageSquare
                          className={`h-3.5 w-3.5 flex-shrink-0 ${isActive ? 'text-sky-400' : 'text-slate-500'}`}
                        />
                        <div className="flex-1 overflow-hidden">
                          <span className="font-medium text-xs truncate block">
                            {sess.title || 'Conversa de Arquitetura'}
                          </span>
                          <span className="text-[9px] text-slate-500 font-mono block">
                            {formatDate(sess.updatedAt)}
                          </span>
                        </div>
                      </div>
                      <button
                        onClick={(e) => handleDeleteSession(e, sess.id)}
                        title="Excluir chat"
                        className="opacity-0 group-hover:opacity-100 text-slate-500 hover:text-red-400 p-1 rounded transition-all"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          {/* Section: Architecture Blueprints */}
          <div>
            <div className="flex items-center justify-between mb-2 px-1">
              <span className="text-[10px] font-mono uppercase tracking-wider text-slate-500 font-bold">
                Blueprints Arquiteturais
              </span>
              <button
                onClick={onOpenTemplatesModal}
                className="text-[10px] text-sky-400 hover:underline flex items-center gap-0.5 cursor-pointer"
              >
                Catálogo
              </button>
            </div>
            <div className="space-y-1.5">
              {templates.slice(0, 3).map((tmpl) => (
                <button
                  key={tmpl.id}
                  onClick={() => {
                    onSelectTemplate(tmpl);
                    onClose();
                  }}
                  className="w-full text-left p-2.5 rounded-xl border border-slate-800/80 bg-slate-900/40 hover:border-sky-500/40 hover:bg-slate-900/80 transition-all group cursor-pointer"
                >
                  <div className="flex items-center justify-between mb-0.5">
                    <span className="text-[9px] font-mono text-sky-400 font-medium">{tmpl.category}</span>
                    <span className="text-[8px] px-1 py-0.2 rounded bg-slate-800 text-slate-400 font-mono">
                      {tmpl.complexity}
                    </span>
                  </div>
                  <span className="font-medium text-[11px] text-slate-300 group-hover:text-slate-100 line-clamp-1">
                    {tmpl.title}
                  </span>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Footer: Architecture Specs */}
        <div className="border-t border-slate-800/80 pt-3 text-[10px] text-slate-500 font-mono space-y-1">
          <div className="flex items-center justify-between">
            <span>Backend Core:</span>
            <span className="text-emerald-400 font-semibold">Go (Clean Arch)</span>
          </div>
          <div className="flex items-center justify-between">
            <span>AI Model:</span>
            <span className="text-sky-400 font-semibold">Gemini 3.5 Flash Lite</span>
          </div>
          <div className="flex items-center justify-between">
            <span>FinOps & Cloud:</span>
            <span className="text-amber-400 font-semibold">AWS • GCP • Azure • OCI</span>
          </div>
        </div>
      </aside>
    </>
  );
};
