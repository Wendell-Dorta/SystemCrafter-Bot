'use client';

import React, { useState, useEffect } from 'react';
import { ArchitectureTemplate } from '@/types';
import { fetchTemplates } from '@/lib/api';
import { MermaidViewer } from './MermaidViewer';
import { X, Layers, Sparkles, DollarSign, Tag, Check, ArrowRight, ShieldCheck } from 'lucide-react';

interface TemplateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelectTemplate: (template: ArchitectureTemplate) => void;
}

export const TemplateModal: React.FC<TemplateModalProps> = ({ isOpen, onClose, onSelectTemplate }) => {
  const [templates, setTemplates] = useState<ArchitectureTemplate[]>([]);
  const [selectedId, setSelectedId] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (isOpen) {
      setLoading(true);
      fetchTemplates().then((list) => {
        setTemplates(list || []);
        if (list && list.length > 0 && !selectedId) {
          setSelectedId(list[0].id);
        }
        setLoading(false);
      });
    }
  }, [isOpen, selectedId]);

  if (!isOpen) return null;

  const currentTemplate = templates.find((t) => t.id === selectedId) || templates[0];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-md p-2.5 sm:p-4 animate-in fade-in duration-200">
      <div className="relative flex h-[92dvh] sm:h-[88vh] w-full max-w-6xl flex-col rounded-2xl border border-slate-700/80 bg-slate-950 shadow-2xl overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/90 px-4 py-3 sm:px-6 sm:py-4 flex-shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-sky-500/10 text-sky-400 border border-sky-500/20 flex-shrink-0">
              <Layers className="h-4 w-4" />
            </div>
            <div>
              <h2 className="text-xs sm:text-sm font-semibold text-slate-100">Blueprints Arquiteturais de Referência</h2>
              <p className="text-[10px] sm:text-xs text-slate-400">Padrões de software distribuído e cloud native</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-800 hover:text-white transition-colors cursor-pointer"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Content Body */}
        <div className="flex flex-col md:flex-row flex-1 overflow-hidden">
          {/* Left Sidebar List */}
          <div className="w-full md:w-80 max-h-44 md:max-h-none border-b md:border-b-0 md:border-r border-slate-800/80 bg-slate-900/40 p-3 sm:p-4 overflow-y-auto space-y-2 flex-shrink-0">
            <span className="text-[10px] font-bold font-mono uppercase tracking-wider text-slate-500 px-1 block mb-1.5">
              Selecione um Modelo ({templates.length})
            </span>

            {loading ? (
              <div className="space-y-2">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-14 rounded-xl bg-slate-800/40 animate-pulse" />
                ))}
              </div>
            ) : (
              templates.map((tmpl) => {
                const isSelected = tmpl.id === selectedId;
                return (
                  <button
                    key={tmpl.id}
                    onClick={() => setSelectedId(tmpl.id)}
                    className={`w-full text-left p-2.5 sm:p-3 rounded-xl border transition-all cursor-pointer ${
                      isSelected
                        ? 'border-sky-500/60 bg-sky-950/30 text-slate-100 shadow-md ring-1 ring-sky-500/30'
                        : 'border-slate-800/80 bg-slate-950/40 text-slate-400 hover:border-slate-700 hover:bg-slate-900/60 hover:text-slate-200'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[10px] font-mono font-semibold text-sky-400">{tmpl.category}</span>
                      <span className="text-[8px] sm:text-[9px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-300 font-mono">
                        {tmpl.complexity}
                      </span>
                    </div>
                    <span className="font-semibold text-xs line-clamp-1">{tmpl.title}</span>
                  </button>
                );
              })
            )}
          </div>

          {/* Right Detail Pane */}
          {currentTemplate && (
            <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-4 sm:space-y-6">
              <div>
                <div className="flex items-center gap-2 mb-1.5">
                  <span className="text-xs font-mono font-semibold text-sky-400 px-2 py-0.5 rounded-full bg-sky-500/10 border border-sky-500/20">
                    {currentTemplate.category}
                  </span>
                  <span className="text-xs text-slate-400">• Complexidade: {currentTemplate.complexity}</span>
                </div>
                <h3 className="text-lg sm:text-xl font-bold text-slate-100">{currentTemplate.title}</h3>
                <p className="mt-1 text-xs text-slate-300 leading-relaxed">{currentTemplate.description}</p>
              </div>

              {/* Tags */}
              <div className="flex flex-wrap gap-1.5">
                {currentTemplate.tags?.map((t) => (
                  <span
                    key={t}
                    className="inline-flex items-center gap-1 rounded-md bg-slate-900 px-2 py-0.5 text-[10px] sm:text-[11px] font-mono text-slate-300 border border-slate-800"
                  >
                    <Tag className="h-2.5 w-2.5 sm:h-3 sm:w-3 text-slate-500" />
                    {t}
                  </span>
                ))}
              </div>

              {/* Recommended Stack & Cost */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
                <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-3.5 sm:p-4">
                  <span className="font-semibold text-xs text-slate-200 block mb-2 flex items-center gap-1.5">
                    <Layers className="h-3.5 w-3.5 text-sky-400" />
                    Stack Recomendada:
                  </span>
                  <ul className="space-y-1 text-xs text-slate-300 font-mono">
                    {currentTemplate.recommendedStack?.map((s, idx) => (
                      <li key={idx} className="flex items-center gap-1.5">
                        <span className="h-1.5 w-1.5 rounded-full bg-sky-400 flex-shrink-0" />
                        <span className="text-[11px]">{s}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-3.5 sm:p-4">
                  <span className="font-semibold text-xs text-slate-200 block mb-2 flex items-center gap-1.5">
                    <DollarSign className="h-3.5 w-3.5 text-emerald-400" />
                    Estimativa Cloud FinOps:
                  </span>
                  <p className="text-xs text-slate-300 leading-relaxed font-mono">
                    {currentTemplate.estimatedCost || 'Estimativa personalizada gerada sob demanda.'}
                  </p>
                </div>
              </div>

              {/* Diagram */}
              {currentTemplate.mermaidDiagram && (
                <div>
                  <span className="font-semibold text-xs text-slate-300 block mb-2 font-mono uppercase tracking-wider">
                    Diagrama Arquitetural de Referência:
                  </span>
                  <MermaidViewer chart={currentTemplate.mermaidDiagram} />
                </div>
              )}

              {/* Action Button */}
              <div className="pt-2">
                <button
                  onClick={() => {
                    onSelectTemplate(currentTemplate);
                    onClose();
                  }}
                  className="w-full flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-sky-500 to-blue-600 px-6 py-3 text-xs sm:text-sm font-semibold text-white shadow-lg shadow-sky-500/25 hover:from-sky-400 hover:to-blue-500 transition-all cursor-pointer"
                >
                  <Sparkles className="h-4 w-4" />
                  Iniciar Projeto com este Blueprint
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
