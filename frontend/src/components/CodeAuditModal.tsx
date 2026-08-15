'use client';

import React, { useState } from 'react';
import {
  ShieldAlert,
  ShieldCheck,
  AlertTriangle,
  X,
  Loader2,
  Lock,
  ArrowRight,
  RefreshCw,
  Sparkles,
  GitBranch,
} from 'lucide-react';
import { CodeAuditReport } from '@/types';
import { auditGitHubRepo } from '@/lib/api';
import { MermaidViewer } from './MermaidViewer';

interface CodeAuditModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentSessionId: string;
  onAuditCompleted: (report: CodeAuditReport) => void;
}

const GitHubIcon: React.FC<{ className?: string }> = ({ className = 'h-4 w-4' }) => (
  <svg className={`${className} fill-current`} viewBox="0 0 24 24">
    <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
  </svg>
);

export const CodeAuditModal: React.FC<CodeAuditModalProps> = ({
  isOpen,
  onClose,
  currentSessionId,
  onAuditCompleted,
}) => {
  const [githubUrl, setGithubUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [report, setReport] = useState<CodeAuditReport | null>(null);

  if (!isOpen) return null;

  const handleGitHubSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!githubUrl.trim() || isLoading) return;

    setIsLoading(true);
    setError(null);
    try {
      const res = await auditGitHubRepo(githubUrl.trim(), currentSessionId);
      setReport(res);
      onAuditCompleted(res);
    } catch (err: any) {
      setError(err.message || 'Erro ao auditar repositório GitHub');
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    setReport(null);
    setError(null);
    setGithubUrl('');
  };

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'CRITICAL':
        return 'text-rose-400 bg-rose-950/60 border-rose-800';
      case 'HIGH':
        return 'text-amber-400 bg-amber-950/60 border-amber-800';
      case 'MEDIUM':
        return 'text-yellow-400 bg-yellow-950/60 border-yellow-800';
      default:
        return 'text-emerald-400 bg-emerald-950/60 border-emerald-800';
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-md p-2.5 sm:p-4 animate-in fade-in duration-200">
      <div className="relative flex max-h-[92dvh] sm:max-h-[90vh] w-full max-w-3xl flex-col rounded-2xl border border-slate-700/70 bg-slate-950 shadow-2xl overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/90 px-4 py-3 sm:px-6 sm:py-4 flex-shrink-0">
          <div className="flex items-center gap-2.5 sm:gap-3 overflow-hidden">
            <div className="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-xl bg-gradient-to-tr from-sky-500 to-indigo-600 text-white shadow-md flex-shrink-0">
              <ShieldAlert className="h-4 w-4 sm:h-5 sm:w-5" />
            </div>
            <div className="overflow-hidden">
              <h2 className="text-sm sm:text-base font-bold text-slate-100 flex items-center gap-1.5 sm:gap-2 truncate">
                Auditar Repositório GitHub
                <span className="rounded-full bg-emerald-950/80 px-2 py-0.2 text-[9px] sm:text-[10px] font-medium text-emerald-400 border border-emerald-700/50 hidden xs:inline">
                  Grounded AI
                </span>
              </h2>
              <p className="text-[10px] sm:text-xs text-slate-400 truncate">
                Inspeção estática da arquitetura real, detecção de stack e vulnerabilidades
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            disabled={isLoading}
            className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-800 hover:text-white transition-colors cursor-pointer flex-shrink-0 disabled:opacity-40"
          >
            <X className="h-4 w-4 sm:h-5 sm:w-5" />
          </button>
        </div>

        {/* Body Content */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-4 sm:space-y-6">
          {error && (
            <div className="rounded-xl border border-rose-500/40 bg-rose-950/30 p-3.5 sm:p-4 text-xs text-rose-300 flex items-start gap-3">
              <AlertTriangle className="h-4 w-4 sm:h-5 sm:w-5 text-rose-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-semibold mb-0.5">Erro na Auditoria</p>
                <p>{error}</p>
              </div>
            </div>
          )}

          {!report ? (
            <form onSubmit={handleGitHubSubmit} className="space-y-4 pt-1">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1.5">
                  URL do Repositório Público no GitHub
                </label>
                <div className="relative">
                  <div className="absolute left-3.5 top-3 text-slate-500">
                    <GitHubIcon className="h-4 w-4" />
                  </div>
                  <input
                    type="text"
                    value={githubUrl}
                    onChange={(e) => setGithubUrl(e.target.value)}
                    placeholder="https://github.com/zen-browser/desktop"
                    disabled={isLoading}
                    className="w-full rounded-xl border border-slate-800 bg-slate-900/90 pl-10 pr-4 py-2.5 text-xs text-slate-100 placeholder-slate-500 focus:border-sky-500 focus:outline-none focus:ring-1 focus:ring-sky-500 disabled:opacity-50"
                  />
                </div>
                <div className="mt-2 flex items-center gap-1.5 text-[11px] text-slate-400">
                  <GitBranch className="h-3.5 w-3.5 text-sky-400 flex-shrink-0" />
                  <span>
                    Detecta automaticamente a branch padrão (<strong>dev</strong>, <strong>main</strong>, <strong>master</strong>, etc.) ou branches específicas via URL.
                  </span>
                </div>
              </div>

              {/* Security Notice */}
              <div className="rounded-xl border border-slate-800 bg-slate-900/40 p-3.5 sm:p-4 text-xs text-slate-400 space-y-2">
                <div className="flex items-center gap-2 text-emerald-400 font-semibold">
                  <Lock className="h-4 w-4" />
                  Filtro Automático de Segurança & Mascaramento
                </div>
                <p className="text-[11px] leading-relaxed text-slate-400">
                  Variáveis <code className="text-amber-300">.env</code>, chaves privadas <code className="text-amber-300">*.key</code>, <code className="text-amber-300">credentials.json</code> e segredos são <strong>automaticamente filtrados e protegidos</strong> antes do envio para a IA.
                </p>
              </div>

              <button
                type="submit"
                disabled={isLoading || !githubUrl.trim()}
                className="flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-600 py-3 text-xs font-semibold text-white shadow-lg hover:from-sky-400 hover:to-indigo-500 transition-all disabled:opacity-50 cursor-pointer"
              >
                {isLoading ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Inspecionando estrutura & Auditando código com IA...
                  </>
                ) : (
                  <>
                    <ShieldCheck className="h-4 w-4" />
                    Iniciar Auditoria no GitHub
                  </>
                )}
              </button>
            </form>
          ) : (
            /* Audit Report View */
            <div className="space-y-6">
              {/* Score Header */}
              <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 rounded-xl border border-slate-800 bg-slate-900/60 p-4">
                <div>
                  <h3 className="text-base font-bold text-slate-100">
                    Relatório de Auditoria: {report.repositoryName}
                  </h3>
                  <p className="text-xs text-slate-400 mt-0.5">
                    {report.analyzedFilesCount} arquivos inspecionados • {report.filteredFilesCount} arquivos sensíveis filtrados
                  </p>
                  <div className="flex flex-wrap gap-1.5 mt-2">
                    {report.detectedStack?.map((tech, idx) => (
                      <span
                        key={idx}
                        className="rounded-md bg-slate-800 px-2 py-0.5 text-[10px] font-mono text-sky-300 border border-slate-700"
                      >
                        {tech}
                      </span>
                    ))}
                  </div>
                </div>

                <div className={`rounded-xl border px-4 py-2 text-center ${getRiskColor(report.riskLevel)}`}>
                  <div className="text-2xl font-black">{report.riskScore}/100</div>
                  <div className="text-[10px] font-bold tracking-wider uppercase">
                    Risco {report.riskLevel}
                  </div>
                </div>
              </div>

              {/* Architecture Summary */}
              <div>
                <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                  Visão Geral da Arquitetura Detectada
                </h4>
                <p className="text-xs leading-relaxed text-slate-300 rounded-xl border border-slate-800 bg-slate-900/40 p-4 whitespace-pre-wrap">
                  {report.architectureSummary}
                </p>
              </div>

              {/* Vulnerabilities & Risks */}
              {report.vulnerabilities && report.vulnerabilities.length > 0 && (
                <div>
                  <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                    Vulnerabilidades & Pontos Críticos Encontrados ({report.vulnerabilities.length})
                  </h4>
                  <div className="space-y-2">
                    {report.vulnerabilities.map((vuln, idx) => (
                      <div
                        key={idx}
                        className="rounded-xl border border-slate-800 bg-slate-900/40 p-3.5 space-y-1.5"
                      >
                        <div className="flex items-center justify-between">
                          <span className="font-semibold text-xs text-slate-200">
                            {vuln.id}: {vuln.category}
                          </span>
                          <span
                            className={`rounded-full px-2 py-0.5 text-[10px] font-bold font-mono border ${getRiskColor(
                              vuln.severity
                            )}`}
                          >
                            {vuln.severity}
                          </span>
                        </div>
                        <p className="text-xs text-slate-300 leading-relaxed">{vuln.description}</p>
                        <div className="rounded-lg bg-emerald-950/30 border border-emerald-800/40 p-2.5 text-xs text-emerald-300">
                          <strong>👉 Remediação Recomendada:</strong> {vuln.remediation}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Architecture Diagram */}
              {report.mermaidDiagram && (
                <div>
                  <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                    Diagrama Arquitetural Extraído do Código
                  </h4>
                  <MermaidViewer chart={report.mermaidDiagram} />
                </div>
              )}

              {/* Recommendations */}
              {report.recommendations && report.recommendations.length > 0 && (
                <div>
                  <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                    Recomendações Práticas de Engenharia
                  </h4>
                  <ul className="space-y-1.5 text-xs text-slate-300 list-disc list-inside">
                    {report.recommendations.map((rec, idx) => (
                      <li key={idx}>{rec}</li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Actions Footer */}
              <div className="flex items-center justify-end gap-3 pt-4 border-t border-slate-800">
                <button
                  onClick={handleReset}
                  className="inline-flex items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800 hover:bg-slate-700 px-4 py-2 text-xs font-medium text-slate-200 transition-colors cursor-pointer"
                >
                  <RefreshCw className="h-3.5 w-3.5" />
                  Nova Auditoria
                </button>

                <button
                  onClick={onClose}
                  className="inline-flex items-center gap-1.5 rounded-xl bg-sky-600 hover:bg-sky-500 px-5 py-2 text-xs font-semibold text-white shadow-md transition-all cursor-pointer"
                >
                  Fechar & Visualizar no Chat
                  <ArrowRight className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
