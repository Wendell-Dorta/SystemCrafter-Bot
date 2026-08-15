'use client';

import React from 'react';
import { DollarSign, Shield, Cpu, GitBranch, Sparkles, Layers, ArrowRight } from 'lucide-react';

interface QuickPromptsProps {
  onSelectPrompt: (prompt: string) => void;
  onOpenTemplates: () => void;
}

export const QuickPrompts: React.FC<QuickPromptsProps> = ({ onSelectPrompt, onOpenTemplates }) => {
  const prompts = [
    {
      title: 'Custos Multi-Cloud (AWS/GCP/Azure/OCI)',
      category: 'FinOps & Cloud',
      icon: DollarSign,
      color: 'text-emerald-400 border-emerald-500/20 bg-emerald-500/10',
      prompt:
        'Poderia estimar os custos mensais de uma arquitetura em nuvem para um SaaS com 15 milhões de requisições/mês (comparando AWS, GCP, Azure e Oracle Cloud OCI) usando Kubernetes, PostgreSQL e Redis?',
    },
    {
      title: 'Modelagem de Ameaças & Compliance',
      category: 'Security & Compliance',
      icon: Shield,
      color: 'text-amber-400 border-amber-500/20 bg-amber-500/10',
      prompt:
        'Gostaria de fazer uma modelagem de ameaças e revisão de conformidade de uma arquitetura de microsserviços financeiros, avaliando regras de LGPD, PCI-DSS, mTLS e pontos únicos de falha (SPOF).',
    },
    {
      title: 'Transactional Outbox & CQRS',
      category: 'Distributed Patterns',
      icon: GitBranch,
      color: 'text-sky-400 border-sky-500/20 bg-sky-500/10',
      prompt:
        'Como estruturar o padrão Transactional Outbox com Debezium e Kafka em Go para garantir que mensagens nunca sejam perdidas em falhas do broker? Gere o diagrama Mermaid.',
    },
    {
      title: 'Matriz de Stack: Go vs Node',
      category: 'Tech Trade-offs',
      icon: Cpu,
      color: 'text-indigo-400 border-indigo-500/20 bg-indigo-500/10',
      prompt:
        'Poderia gerar uma matriz de trade-offs comparando Go vs Node.js vs Python para um gateway de alta concorrência com 100k conexões simultâneas?',
    },
  ];

  return (
    <div className="mx-auto w-full max-w-4xl px-3.5 sm:px-6 pt-6 pb-12 sm:pt-10 sm:pb-16 space-y-4 sm:space-y-6 animate-in fade-in duration-500">
      {/* Hero Banner */}
      <div className="text-center space-y-2 sm:space-y-3">
        <div className="inline-flex items-center gap-1.5 sm:gap-2 rounded-full border border-sky-500/30 bg-sky-950/40 px-3 py-1 text-[11px] sm:text-xs font-mono text-sky-400 shadow-sm backdrop-blur-md">
          <Sparkles className="h-3 w-3 sm:h-3.5 sm:w-3.5 animate-pulse" />
          <span>Arquiteto de Software & Cloud Advisor</span>
        </div>
        <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-slate-100 leading-tight">
          Projete e Audite Sistemas de <span className="text-sky-400">Alta Performance</span>
        </h1>
        <p className="text-xs sm:text-sm text-slate-400 max-w-xl mx-auto leading-relaxed">
          Desenhe diagramas interativos em Mermaid, calcule custos multi-cloud e audite
          vulnerabilidades com inteligência artificial e backend em Clean Architecture.
        </p>
      </div>

      {/* Blueprints Banner */}
      <div className="rounded-2xl border border-slate-800 bg-gradient-to-r from-slate-900/90 via-slate-900/60 to-slate-950 p-3 sm:p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 sm:gap-4 shadow-xl">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-xl bg-sky-500/10 text-sky-400 border border-sky-500/20 flex-shrink-0">
            <Layers className="h-4 w-4 sm:h-5 sm:w-5" />
          </div>
          <div>
            <span className="font-semibold text-xs sm:text-sm text-slate-200 block">Blueprints Arquiteturais</span>
            <span className="text-[10px] sm:text-xs text-slate-400 block">Fintech, E-Commerce, SaaS RAG e Telemetria</span>
          </div>
        </div>
        <button
          onClick={onOpenTemplates}
          className="inline-flex items-center gap-2 rounded-xl bg-slate-800 hover:bg-slate-700 px-3.5 py-1.5 sm:px-4 sm:py-2 text-xs font-medium text-slate-200 transition-colors border border-slate-700 w-full sm:w-auto justify-center cursor-pointer"
        >
          Ver Catálogo
          <ArrowRight className="h-3.5 w-3.5" />
        </button>
      </div>

      {/* Grid of Starter Prompts */}
      <div>
        <span className="text-[10px] sm:text-[11px] font-mono uppercase tracking-wider text-slate-500 font-semibold block mb-2.5 px-1">
          Sugestões Rápidas de Arquitetura:
        </span>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 sm:gap-3">
          {prompts.map((p, idx) => {
            const Icon = p.icon;
            return (
              <button
                key={idx}
                onClick={() => onSelectPrompt(p.prompt)}
                className="group relative flex flex-col justify-between rounded-xl border border-slate-800/80 bg-slate-900/50 p-3.5 sm:p-4 text-left shadow-sm transition-all hover:border-sky-500/50 hover:bg-slate-900/90 hover:shadow-md cursor-pointer active:scale-[0.98]"
              >
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-[9px] sm:text-[10px] font-mono text-slate-400 font-medium">{p.category}</span>
                    <div className={`flex h-6 w-6 items-center justify-center rounded-md border ${p.color}`}>
                      <Icon className="h-3.5 w-3.5" />
                    </div>
                  </div>
                  <h4 className="font-semibold text-xs sm:text-sm text-slate-200 group-hover:text-sky-300 transition-colors">
                    {p.title}
                  </h4>
                  <p className="mt-1 text-[11px] text-slate-400 line-clamp-2 leading-relaxed">{p.prompt}</p>
                </div>
                <div className="mt-2.5 flex items-center gap-1 text-[10px] sm:text-[11px] font-medium text-sky-400 group-hover:translate-x-0.5 transition-transform">
                  <span>Executar</span>
                  <ArrowRight className="h-3 w-3" />
                </div>
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
};
