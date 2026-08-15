'use client';

import React from 'react';
import { ToolEvent } from '@/types';
import {
  DollarSign,
  ShieldAlert,
  ShieldCheck,
  Cpu,
  Layers,
  CheckCircle2,
  AlertTriangle,
  Server,
  Zap,
  Sparkles,
} from 'lucide-react';
import { formatCurrency } from '@/lib/utils';

interface ToolExecutionCardProps {
  toolEvent: ToolEvent;
}

export const ToolExecutionCard: React.FC<ToolExecutionCardProps> = ({ toolEvent }) => {
  const { toolName, result, durationMs, error } = toolEvent;

  if (error) {
    return (
      <div className="my-3 rounded-xl border border-red-500/30 bg-red-950/20 p-3.5 text-xs text-red-300">
        <div className="flex items-center gap-2 font-medium">
          <AlertTriangle className="h-4 w-4 text-red-400" />
          <span>Erro na execução da ferramenta: {toolName}</span>
        </div>
        <p className="mt-1 text-slate-400">{error}</p>
      </div>
    );
  }

  // 1. Cloud Cost Tool Card
  if (toolName === 'estimate_cloud_costs' && result) {
    const cost = result;
    return (
      <div className="my-3 overflow-hidden rounded-xl border border-emerald-500/30 bg-slate-900/90 p-4 shadow-lg backdrop-blur-md">
        <div className="flex items-center justify-between border-b border-slate-800 pb-3">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              <DollarSign className="h-4 w-4" />
            </div>
            <div>
              <span className="font-semibold text-xs text-slate-200">
                Estimativa de Custos em Nuvem ({cost.provider} - {cost.workloadTier})
              </span>
              <span className="block text-[10px] text-slate-400">Calculado em {durationMs}ms</span>
            </div>
          </div>
          <div className="text-right">
            <span className="text-xs text-slate-400 block font-mono">Total Estimado</span>
            <span className="text-base font-bold text-emerald-400 font-mono">
              {formatCurrency(cost.monthlyTotalUSD || 0)} / mês
            </span>
          </div>
        </div>

        {/* Breakdown */}
        {cost.breakdownUSD && (
          <div className="mt-3">
            <span className="text-[10px] uppercase font-mono tracking-wider text-slate-500 font-semibold block mb-1.5">
              Decomposição de Serviços:
            </span>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 text-xs">
              {Object.entries(cost.breakdownUSD).map(([key, val]: [string, any]) => (
                <div key={key} className="rounded-lg bg-slate-950/60 p-2 border border-slate-800/80">
                  <span className="text-[10px] text-slate-400 block">{key}</span>
                  <span className="font-semibold text-slate-200 font-mono">{formatCurrency(val)}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Multi-Cloud Comparative Pricing */}
        {cost.comparativePricing && Object.keys(cost.comparativePricing).length > 0 && (
          <div className="mt-3 pt-2.5 border-t border-slate-800/80">
            <span className="text-[10px] uppercase font-mono tracking-wider text-sky-400 font-semibold block mb-2">
              Comparativo Multi-Cloud (Mesma Carga):
            </span>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
              {Object.entries(cost.comparativePricing).map(([pName, pTotal]: [string, any]) => {
                const isSelected = pName.toUpperCase() === (cost.provider || '').toUpperCase();
                const isOCI = pName.toUpperCase() === 'ORACLE';
                return (
                  <div
                    key={pName}
                    className={`rounded-lg p-2.5 border transition-all ${
                      isSelected
                        ? 'border-emerald-500/50 bg-emerald-500/10'
                        : 'border-slate-800 bg-slate-950/40 hover:border-slate-700'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[11px] font-bold text-slate-200">
                        {pName === 'ORACLE' ? 'Oracle (OCI)' : pName}
                      </span>
                      {isSelected && (
                        <span className="rounded bg-emerald-500/20 px-1 py-0.2 text-[8px] font-mono text-emerald-400">
                          Ativo
                        </span>
                      )}
                      {!isSelected && isOCI && (
                        <span className="rounded bg-sky-500/20 px-1 py-0.2 text-[8px] font-mono text-sky-400">
                          Econômico
                        </span>
                      )}
                    </div>
                    <span className="font-bold text-slate-100 font-mono text-xs block">
                      {formatCurrency(pTotal)}
                    </span>
                    <span className="text-[9px] text-slate-500">/ mês</span>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Optimizations */}
        {cost.optimizations && cost.optimizations.length > 0 && (
          <div className="mt-3 pt-2.5 border-t border-slate-800 text-[11px] text-slate-300">
            <span className="font-medium text-emerald-400 flex items-center gap-1 mb-1">
              <Sparkles className="h-3 w-3" /> Estratégias de Economia FinOps:
            </span>
            <ul className="list-disc list-inside space-y-0.5 text-slate-400">
              {cost.optimizations.slice(0, 3).map((opt: string, idx: number) => (
                <li key={idx}>{opt}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    );
  }

  // 2. Security Auditor Tool Card
  if (toolName === 'audit_security_compliance' && result) {
    const audit = result;
    const isCritical = audit.riskLevel === 'CRITICAL' || audit.riskScore >= 70;
    const isHigh = audit.riskLevel === 'HIGH' || (audit.riskScore >= 45 && audit.riskScore < 70);
    const isLow = audit.riskLevel === 'LOW';

    const badgeColor = isCritical
      ? 'border-red-500/40 bg-red-950/30 text-red-400'
      : isHigh
      ? 'border-amber-500/40 bg-amber-950/30 text-amber-400'
      : 'border-emerald-500/40 bg-emerald-950/30 text-emerald-400';

    return (
      <div className="my-3 overflow-hidden rounded-xl border border-slate-700/80 bg-slate-900/90 p-4 shadow-lg backdrop-blur-md">
        <div className="flex items-center justify-between border-b border-slate-800 pb-3">
          <div className="flex items-center gap-2.5">
            <div className={`flex h-8 w-8 items-center justify-center rounded-lg border ${badgeColor}`}>
              {isLow ? <ShieldCheck className="h-4 w-4" /> : <ShieldAlert className="h-4 w-4" />}
            </div>
            <div>
              <span className="font-semibold text-xs text-slate-200">Auditoria de Segurança & Conformidade</span>
              <span className="block text-[10px] text-slate-400">Auditoria executada em {durationMs}ms</span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-[10px] text-slate-400 font-mono">Risco:</span>
            <span className={`rounded-full px-2.5 py-0.5 text-xs font-semibold font-mono border ${badgeColor}`}>
              {audit.riskScore}/100 ({audit.riskLevel})
            </span>
          </div>
        </div>

        {/* Vulnerabilities */}
        {audit.vulnerabilities && audit.vulnerabilities.length > 0 && (
          <div className="mt-3 space-y-1.5">
            <span className="text-[11px] font-medium text-slate-300 block">Vulnerabilidades e Riscos Detectados:</span>
            {audit.vulnerabilities.map((v: any, idx: number) => (
              <div key={idx} className="rounded-lg bg-slate-950/60 p-2.5 border border-slate-800/80 text-xs">
                <div className="flex items-center gap-2 mb-1">
                  <span
                    className={`rounded px-1.5 py-0.2 text-[9px] font-bold ${
                      v.severity === 'CRITICAL'
                        ? 'bg-red-500/20 text-red-400'
                        : v.severity === 'HIGH'
                        ? 'bg-amber-500/20 text-amber-400'
                        : 'bg-sky-500/20 text-sky-400'
                    }`}
                  >
                    {v.severity}
                  </span>
                  <span className="font-medium text-slate-200">{v.category}</span>
                </div>
                <p className="text-[11px] text-slate-400 mb-1">{v.description}</p>
                <p className="text-[10px] text-emerald-400">👉 Remediação: {v.remediation}</p>
              </div>
            ))}
          </div>
        )}

        {/* Compliance Status */}
        {audit.complianceStatus && Object.keys(audit.complianceStatus).length > 0 && (
          <div className="mt-3 pt-2.5 border-t border-slate-800 grid grid-cols-1 sm:grid-cols-2 gap-1.5 text-[10px]">
            {Object.entries(audit.complianceStatus).map(([std, status]: [string, any]) => (
              <div key={std} className="flex items-center gap-1.5 text-slate-300">
                <CheckCircle2 className="h-3 w-3 text-sky-400 flex-shrink-0" />
                <span className="font-semibold text-sky-300">{std}:</span>
                <span className="text-slate-400 truncate">{status}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  // 3. Tech Stack Matrix Tool Card
  if (toolName === 'generate_tech_stack_matrix' && result) {
    const matrix = result;
    return (
      <div className="my-3 overflow-hidden rounded-xl border border-sky-500/30 bg-slate-900/90 p-4 shadow-lg backdrop-blur-md">
        <div className="flex items-center gap-2.5 border-b border-slate-800 pb-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-sky-500/10 text-sky-400 border border-sky-500/20">
            <Cpu className="h-4 w-4" />
          </div>
          <div>
            <span className="font-semibold text-xs text-slate-200">
              Matriz de Tecnologias & Trade-offs ({matrix.workloadType})
            </span>
            <span className="block text-[10px] text-slate-400">Processado em {durationMs}ms</span>
          </div>
        </div>

        {matrix.recommended && (
          <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
            {Object.entries(matrix.recommended).map(([layer, tech]: [string, any]) => (
              <div key={layer} className="rounded-lg bg-slate-950/60 p-2.5 border border-slate-800/80">
                <span className="text-[10px] text-sky-400 font-semibold block">{layer}</span>
                <span className="text-slate-200 text-[11px] font-mono">{tech}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  // 4. Default Tool Card
  return (
    <div className="my-2 rounded-lg border border-slate-800 bg-slate-950/50 p-2.5 text-xs text-slate-400 flex items-center justify-between">
      <div className="flex items-center gap-2">
        <Layers className="h-3.5 w-3.5 text-sky-400" />
        <span>Ferramenta executada: <strong className="text-slate-200">{toolName}</strong></span>
      </div>
      <span className="font-mono text-[10px] text-slate-500">{durationMs}ms</span>
    </div>
  );
};
