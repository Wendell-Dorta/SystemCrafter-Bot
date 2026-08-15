'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { Session } from '@/types';
import { fetchSessions, deleteSession } from '@/lib/api';
import {
  MessageSquare,
  Search,
  Plus,
  Trash2,
  X,
  RefreshCw,
  Calendar,
  Clock,
  ArrowRight,
  History,
  CheckCircle2,
} from 'lucide-react';
import { formatDate } from '@/lib/utils';

interface HistoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentSessionId: string;
  onSelectSession: (sessionId: string) => void;
  onNewSession: () => void;
}

export const HistoryModal: React.FC<HistoryModalProps> = ({
  isOpen,
  onClose,
  currentSessionId,
  onSelectSession,
  onNewSession,
}) => {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const loadSessions = async () => {
    setLoading(true);
    try {
      const list = await fetchSessions();
      setSessions(list || []);
    } catch (err) {
      console.error('Failed to load chat history:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      loadSessions();
      setSearchQuery('');
    }
  }, [isOpen]);

  const filteredSessions = useMemo(() => {
    if (!searchQuery.trim()) return sessions;
    const query = searchQuery.toLowerCase().trim();
    return sessions.filter((s) => {
      const titleMatch = (s.title || '').toLowerCase().includes(query);
      const idMatch = (s.id || '').toLowerCase().includes(query);
      return titleMatch || idMatch;
    });
  }, [sessions, searchQuery]);

  const handleDelete = async (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (!confirm('Deseja realmente excluir esta conversa do histórico?')) return;

    setDeletingId(id);
    try {
      const ok = await deleteSession(id);
      if (ok) {
        setSessions((prev) => prev.filter((s) => s.id !== id));
        if (id === currentSessionId) {
          onNewSession();
        }
      }
    } catch (err) {
      console.error('Failed to delete session:', err);
    } finally {
      setDeletingId(null);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-2.5 sm:p-4 bg-black/80 backdrop-blur-md animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-3xl max-h-[92dvh] sm:max-h-[85vh] flex flex-col shadow-2xl overflow-hidden text-slate-200">
        {/* Header */}
        <div className="px-4 py-3 sm:px-6 sm:py-4 border-b border-slate-800 flex items-center justify-between bg-slate-950/70">
          <div className="flex items-center gap-2.5 sm:gap-3 overflow-hidden">
            <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-xl bg-sky-500/10 border border-sky-500/30 flex items-center justify-center text-sky-400 flex-shrink-0">
              <History className="w-4 h-4 sm:w-5 sm:h-5" />
            </div>
            <div className="overflow-hidden">
              <h2 className="text-sm sm:text-lg font-bold text-white flex items-center gap-1.5 sm:gap-2 truncate">
                Histórico de Chats
                <span className="text-[10px] sm:text-xs px-2 py-0.2 rounded-full bg-slate-800 text-sky-300 font-mono font-medium border border-slate-700">
                  {sessions.length}
                </span>
              </h2>
              <p className="text-[10px] sm:text-xs text-slate-400 truncate">
                Pesquise e gerencie suas conversas
              </p>
            </div>
          </div>

          <div className="flex items-center gap-1.5 sm:gap-2 flex-shrink-0">
            <button
              onClick={() => {
                onNewSession();
                onClose();
              }}
              className="inline-flex items-center gap-1 rounded-xl bg-sky-600 hover:bg-sky-500 px-2.5 py-1.5 sm:px-3.5 text-xs font-semibold text-white shadow-md transition-all cursor-pointer"
            >
              <Plus className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">Nova Conversa</span>
            </button>

            <button
              onClick={onClose}
              className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors cursor-pointer"
            >
              <X className="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
          </div>
        </div>

        {/* Search Bar */}
        <div className="p-3 sm:p-4 border-b border-slate-800 bg-slate-950/40">
          <div className="relative flex items-center">
            <Search className="absolute left-3 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Pesquisar por assunto ou palavras-chave..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-slate-800/80 border border-slate-700/80 rounded-xl pl-9 pr-9 py-2 text-xs text-white placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500 transition-all"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 text-slate-400 hover:text-white"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            )}
          </div>
        </div>

        {/* Sessions List */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-2.5">
          {loading ? (
            <div className="py-16 flex flex-col items-center justify-center text-slate-400 gap-2">
              <RefreshCw className="w-6 h-6 animate-spin text-sky-400" />
              <p className="text-xs">Carregando histórico de conversas...</p>
            </div>
          ) : filteredSessions.length === 0 ? (
            <div className="py-16 text-center space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-slate-800/50 border border-slate-700 flex items-center justify-center mx-auto text-slate-500">
                <MessageSquare className="w-6 h-6" />
              </div>
              {searchQuery ? (
                <>
                  <p className="text-sm font-medium text-slate-300">Nenhuma conversa encontrada</p>
                  <p className="text-xs text-slate-500">
                    Nenhum resultado para o termo &quot;<span className="text-sky-400">{searchQuery}</span>&quot;.
                  </p>
                  <button
                    onClick={() => setSearchQuery('')}
                    className="text-xs text-sky-400 hover:underline pt-1"
                  >
                    Limpar busca
                  </button>
                </>
              ) : (
                <>
                  <p className="text-sm font-medium text-slate-300">Seu histórico está vazio</p>
                  <p className="text-xs text-slate-500 max-w-sm mx-auto">
                    Inicie uma nova conversa para projetar sistemas, calcular custos ou inspecionar códigos.
                  </p>
                  <button
                    onClick={() => {
                      onNewSession();
                      onClose();
                    }}
                    className="inline-flex items-center gap-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 px-4 py-2 text-xs font-medium text-slate-200 transition-colors"
                  >
                    <Plus className="w-4 h-4 text-sky-400" /> Criar Primeiro Chat
                  </button>
                </>
              )}
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-2.5">
              {filteredSessions.map((sess) => {
                const isActive = sess.id === currentSessionId;
                const msgCount = sess.messages ? sess.messages.length : 0;
                return (
                  <div
                    key={sess.id}
                    onClick={() => {
                      onSelectSession(sess.id);
                      onClose();
                    }}
                    className={`group relative flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-4 rounded-xl border transition-all cursor-pointer ${
                      isActive
                        ? 'border-sky-500/60 bg-sky-950/30 text-white shadow-md shadow-sky-950/50'
                        : 'border-slate-800/80 bg-slate-900/50 hover:border-slate-700 hover:bg-slate-900/90 text-slate-300'
                    }`}
                  >
                    <div className="flex items-start sm:items-center gap-3 overflow-hidden flex-1">
                      <div
                        className={`w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 border ${
                          isActive
                            ? 'bg-sky-500/20 border-sky-500/40 text-sky-300'
                            : 'bg-slate-800/60 border-slate-700/60 text-slate-400 group-hover:text-slate-200'
                        }`}
                      >
                        <MessageSquare className="w-4 h-4" />
                      </div>

                      <div className="overflow-hidden flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <h4
                            className={`font-semibold text-xs sm:text-sm truncate ${
                              isActive ? 'text-sky-300' : 'text-slate-100 group-hover:text-white'
                            }`}
                          >
                            {sess.title || 'Conversa de Arquitetura'}
                          </h4>
                          {isActive && (
                            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.2 text-[9px] font-mono font-medium text-emerald-400 border border-emerald-500/20">
                              <CheckCircle2 className="w-2.5 h-2.5" /> Ativo
                            </span>
                          )}
                        </div>

                        <div className="flex items-center gap-3 text-[11px] text-slate-400 font-mono">
                          <span className="flex items-center gap-1">
                            <Clock className="w-3 h-3 text-slate-500" />
                            {formatDate(sess.updatedAt)}
                          </span>
                          {msgCount > 0 && (
                            <span>
                              {msgCount} {msgCount === 1 ? 'mensagem' : 'mensagens'}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Action buttons */}
                    <div className="flex items-center justify-end gap-2 border-t sm:border-t-0 border-slate-800/60 pt-2 sm:pt-0">
                      <button
                        onClick={(e) => handleDelete(e, sess.id)}
                        disabled={deletingId === sess.id}
                        title="Excluir conversa"
                        className="p-2 rounded-lg text-slate-500 hover:text-red-400 hover:bg-slate-800 transition-colors cursor-pointer"
                      >
                        {deletingId === sess.id ? (
                          <RefreshCw className="w-4 h-4 animate-spin text-red-400" />
                        ) : (
                          <Trash2 className="w-4 h-4" />
                        )}
                      </button>

                      <div className="flex items-center gap-1 text-xs font-semibold text-sky-400 group-hover:translate-x-1 transition-transform">
                        <span className="hidden sm:inline">{isActive ? 'Continuar' : 'Abrir'}</span>
                        <ArrowRight className="w-3.5 h-3.5" />
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-3 border-t border-slate-800 bg-slate-950/60 flex items-center justify-between text-xs text-slate-400">
          <span className="font-mono text-[11px]">
            {filteredSessions.length} de {sessions.length} exibidos
          </span>
          <button
            onClick={onClose}
            className="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium transition-colors"
          >
            Fechar
          </button>
        </div>
      </div>
    </div>
  );
};
