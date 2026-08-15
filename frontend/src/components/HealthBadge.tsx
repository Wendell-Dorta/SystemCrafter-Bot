'use client';

import React from 'react';
import { useHealth } from '@/hooks/useHealth';
import { Activity, Zap, Database, Shield } from 'lucide-react';

export const HealthBadge: React.FC = () => {
  const { health, isOnline, latencyMs, loading } = useHealth();

  return (
    <div className="flex items-center gap-2 rounded-full border border-slate-800 bg-slate-900/90 px-3 py-1 text-xs text-slate-300 backdrop-blur-md">
      <div className="flex items-center gap-1.5">
        <span
          className={`h-2 w-2 rounded-full ${
            loading
              ? 'bg-amber-400 animate-pulse'
              : isOnline
              ? 'bg-emerald-400 ring-2 ring-emerald-500/20 animate-pulse'
              : 'bg-red-400'
          }`}
        />
        <span className="font-mono text-[11px] font-medium">
          {loading ? 'Conectando...' : isOnline ? 'Go Core Online' : 'Backend Offline'}
        </span>
      </div>

      {isOnline && latencyMs !== null && (
        <>
          <span className="text-slate-600">•</span>
          <span className="font-mono text-[10px] text-emerald-400 flex items-center gap-0.5">
            <Zap className="h-3 w-3" />
            {latencyMs}ms
          </span>
        </>
      )}

      {health?.cacheStats && (
        <>
          <span className="text-slate-600 hidden sm:inline">•</span>
          <span className="font-mono text-[10px] text-sky-400 hidden sm:flex items-center gap-0.5" title="Cache em Memória Ativo">
            <Database className="h-3 w-3" />
            {Math.round(health.cacheStats.hitRatio || 0)}% hit
          </span>
        </>
      )}
    </div>
  );
};
