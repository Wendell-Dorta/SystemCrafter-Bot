'use client';

import React from 'react';
import { HealthBadge } from './HealthBadge';
import { Layers, Plus, Menu, History, Shield, Cpu, Brain } from 'lucide-react';

interface HeaderProps {
  onNewSession: () => void;
  onOpenHistory: () => void;
  onOpenTemplates: () => void;
  onOpenProfile: () => void;
  onOpenCodeAudit: () => void;
  onToggleSidebar: () => void;
  isSidebarOpen: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  onNewSession,
  onOpenHistory,
  onOpenTemplates,
  onOpenProfile,
  onOpenCodeAudit,
  onToggleSidebar,
  isSidebarOpen,
}) => {
  return (
    <header className="sticky top-0 z-40 flex h-14 sm:h-16 w-full flex-shrink-0 items-center justify-between border-b border-slate-800/80 bg-slate-950/85 px-3 sm:px-6 backdrop-blur-xl transition-all">
      {/* Left: Brand & Sidebar Toggle */}
      <div className="flex items-center gap-2 sm:gap-3">
        <button
          onClick={onToggleSidebar}
          title="Abrir menu de navegação"
          aria-label="Abrir menu"
          className="flex h-9 w-9 items-center justify-center rounded-xl text-slate-400 hover:bg-slate-800/80 hover:text-white transition-colors cursor-pointer"
        >
          <Menu className="h-5 w-5" />
        </button>

        <div className="flex items-center gap-2 sm:gap-2.5">
          <div className="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-xl bg-gradient-to-tr from-sky-500 to-blue-600 text-white shadow-md shadow-sky-500/20 flex-shrink-0">
            <Cpu className="h-4 w-4 sm:h-5 sm:w-5" />
          </div>
          <div>
            <div className="flex items-center gap-1.5">
              <span className="font-bold text-xs sm:text-sm text-slate-100 tracking-tight whitespace-nowrap">
                SystemCrafter AI
              </span>
              <span className="hidden xs:inline rounded bg-sky-500/10 px-1.5 py-0.2 text-[8px] sm:text-[9px] font-mono font-semibold text-sky-400 border border-sky-500/20 whitespace-nowrap">
                Go Core
              </span>
            </div>
            <span className="text-[9px] sm:text-[10px] text-slate-400 hidden md:block">
              AI Software Architect & System Designer
            </span>
          </div>
        </div>
      </div>

      {/* Center: Live Health Telemetry (Desktop/TV) */}
      <div className="hidden lg:flex items-center gap-3">
        <HealthBadge />
      </div>

      {/* Right: Actions (Adaptive Mobile-First) */}
      <div className="flex items-center gap-1 sm:gap-2">
        {/* History Button (Always Accessible) */}
        <button
          onClick={onOpenHistory}
          title="Histórico de conversas"
          className="inline-flex items-center gap-1 sm:gap-1.5 rounded-xl border border-sky-500/30 bg-sky-500/10 px-2.5 py-1.5 text-xs font-medium text-sky-300 hover:border-sky-500/50 hover:bg-sky-500/20 transition-all cursor-pointer"
        >
          <History className="h-3.5 w-3.5 text-sky-400 flex-shrink-0" />
          <span className="hidden sm:inline">Histórico</span>
        </button>

        {/* Code Audit (Hidden on extra small mobile, accessible via sidebar) */}
        <button
          onClick={onOpenCodeAudit}
          title="Auditar repositório público do GitHub"
          className="hidden sm:inline-flex items-center gap-1.5 rounded-xl border border-sky-500/30 bg-sky-500/10 px-2.5 py-1.5 text-xs font-medium text-sky-300 hover:border-sky-500/50 hover:bg-sky-500/20 transition-all cursor-pointer"
        >
          <Shield className="h-3.5 w-3.5 text-sky-400 flex-shrink-0" />
          <span className="hidden md:inline">Auditar GitHub</span>
        </button>

        {/* Memory Profile */}
        <button
          onClick={onOpenProfile}
          title="Memória do arquiteto e preferências"
          className="hidden sm:inline-flex items-center gap-1.5 rounded-xl border border-amber-500/30 bg-amber-500/10 px-2.5 py-1.5 text-xs font-medium text-amber-300 hover:border-amber-500/50 hover:bg-amber-500/20 transition-all cursor-pointer"
        >
          <Brain className="h-3.5 w-3.5 text-amber-400 flex-shrink-0" />
          <span className="hidden md:inline">Memória</span>
        </button>

        {/* Blueprints (Tablet & Desktop) */}
        <button
          onClick={onOpenTemplates}
          title="Catálogo de blueprints de arquitetura"
          className="hidden md:inline-flex items-center gap-1.5 rounded-xl border border-slate-800 bg-slate-900/60 px-2.5 py-1.5 text-xs font-medium text-slate-300 hover:border-slate-700 hover:bg-slate-800 hover:text-white transition-all cursor-pointer"
        >
          <Layers className="h-3.5 w-3.5 text-sky-400 flex-shrink-0" />
          <span>Blueprints</span>
        </button>

        {/* New Session Button */}
        <button
          onClick={onNewSession}
          title="Iniciar nova conversa"
          className="inline-flex items-center gap-1 sm:gap-1.5 rounded-xl bg-slate-900 hover:bg-slate-800 border border-slate-700/80 px-2.5 sm:px-3 py-1.5 text-xs font-medium text-slate-200 transition-all shadow-sm cursor-pointer"
        >
          <Plus className="h-3.5 w-3.5 flex-shrink-0 text-sky-400" />
          <span className="hidden sm:inline">Novo Chat</span>
        </button>

        {/* GitHub Link (Large screen/Desktop) */}
        <a
          href="https://github.com/Wendell-Dorta/SystemCrafter-Bot"
          target="_blank"
          rel="noopener noreferrer"
          title="Ver código-fonte no GitHub"
          className="hidden xl:flex h-9 w-9 items-center justify-center rounded-xl border border-slate-800 bg-slate-900/60 text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
        >
          <svg className="h-4 w-4 fill-current" viewBox="0 0 24 24">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
          </svg>
        </a>
      </div>
    </header>
  );
};
