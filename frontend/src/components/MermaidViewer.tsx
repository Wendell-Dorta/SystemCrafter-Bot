'use client';

import React, { useEffect, useRef, useState } from 'react';
import mermaid from 'mermaid';
import { ZoomIn, ZoomOut, RotateCcw, Copy, Check, Maximize2, X, AlertTriangle } from 'lucide-react';

interface MermaidViewerProps {
  chart: string;
  id?: string;
  isStreaming?: boolean;
}

function sanitizeMermaid(rawChart: string): string {
  if (!rawChart) return '';
  let chart = rawChart.trim();

  // 1. Strip markdown code fences if accidentally included
  chart = chart.replace(/^```(?:mermaid)?\n?/i, '').replace(/\n?```$/i, '').trim();

  // 2. Fix AI hallucinations: 'subclass Name' -> 'subgraph "Name"'
  chart = chart.replace(/subclass\s+([^\n{]+)/gi, 'subgraph "$1"');

  // 3. Fix unquoted subgraphs: subgraph ACID Storage (Postgres) -> subgraph "ACID Storage (Postgres)"
  chart = chart.replace(/subgraph\s+([^\n"[{]+?\([^\n"[{]+?\))/gi, 'subgraph "$1"');

  // 4. Clean escaped quotes in labels: [\"Text\"] -> ["Text"]
  chart = chart.replace(/\[\\"(.*?)\\"\]/g, '["$1"]');

  // 5. Replace unquoted node labels containing parentheses: [Texto (Go)] -> ["Texto (Go)"]
  chart = chart.replace(/\[([^"\]\n]*\([^"\]\n]*\)[^"\]\n]*)\]/g, '["$1"]');

  // 6. Replace unquoted node labels containing slashes or ampersands: [Texto / Outro] -> ["Texto / Outro"]
  chart = chart.replace(/\[([^"\]\n]*[\\/&][^"\]\n]*)\]/g, '["$1"]');

  // 7. Ensure double quotes inside already quoted labels don't break syntax: ["Texto ("Sub")"] -> ["Texto (Sub)"]
  chart = chart.replace(/\["([^"\n]*)"([^"\n]*)"([^"\n]*)"\]/g, '["$1 $2 $3"]');

  return chart;
}

export const MermaidViewer: React.FC<MermaidViewerProps> = ({ chart, id, isStreaming }) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [svgContent, setSvgContent] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [zoom, setZoom] = useState<number>(1);
  const [copied, setCopied] = useState<boolean>(false);
  const [isFullscreen, setIsFullscreen] = useState<boolean>(false);

  const uniqueId = useRef(`mermaid-${Math.random().toString(36).substring(2, 9)}`).current;

  useEffect(() => {
    mermaid.initialize({
      startOnLoad: false,
      suppressErrorRendering: true,
      theme: 'dark',
      themeVariables: {
        darkMode: true,
        background: '#0f172a',
        primaryColor: '#1e293b',
        primaryTextColor: '#f8fafc',
        primaryBorderColor: '#38bdf8',
        lineColor: '#64748b',
        secondaryColor: '#334155',
        tertiaryColor: '#1e293b',
      },
      fontFamily: 'inherit',
      securityLevel: 'loose',
    });

    let isMounted = true;

    const renderChart = async () => {
      if (!chart || !chart.trim()) return;
      const cleanChart = sanitizeMermaid(chart);

      try {
        const renderId = `${uniqueId}-${Date.now()}`;
        const { svg } = await mermaid.render(renderId, cleanChart);
        if (isMounted) {
          setSvgContent(svg);
          setError(null);
        }
      } catch (err: any) {
        if (isMounted) {
          // If still streaming tokens, it's expected that incomplete code blocks fail momentarily
          if (isStreaming) {
            setError(null);
          } else {
            console.warn('Mermaid syntax issue auto-handled:', err);
            setError('Diagrama gerado com sintaxe específica.');
          }
        }
      } finally {
        // Clean up any stray error elements inserted into DOM by Mermaid library
        const strayElements = document.querySelectorAll(`[id^="d${uniqueId}"], [id^="mermaid-"]`);
        strayElements.forEach((el) => {
          if (el.parentElement === document.body) {
            el.remove();
          }
        });
      }
    };

    renderChart();

    return () => {
      isMounted = false;
    };
  }, [chart, isStreaming, uniqueId]);

  const handleCopyCode = () => {
    navigator.clipboard.writeText(chart);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleZoomIn = () => setZoom((prev) => Math.min(prev + 0.2, 2.5));
  const handleZoomOut = () => setZoom((prev) => Math.max(prev - 0.2, 0.5));
  const handleResetZoom = () => setZoom(1);

  if (error && !svgContent) {
    return (
      <div className="my-3 overflow-hidden rounded-xl border border-slate-800 bg-slate-950 shadow-md">
        <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/90 px-4 py-2 text-xs text-amber-400 font-mono">
          <span className="flex items-center gap-1.5">
            <AlertTriangle className="w-3.5 h-3.5" />
            Estrutura do Diagrama Mermaid
          </span>
          <button
            onClick={handleCopyCode}
            className="flex items-center gap-1 hover:text-slate-200 transition-colors text-slate-400 cursor-pointer"
          >
            {copied ? <Check className="h-3 w-3 text-emerald-400" /> : <Copy className="h-3 w-3" />}
            <span>{copied ? 'Copiado' : 'Copiar'}</span>
          </button>
        </div>
        <pre className="overflow-x-auto p-4 text-xs font-mono text-sky-300 bg-slate-950 leading-relaxed">
          <code>{chart}</code>
        </pre>
      </div>
    );
  }

  return (
    <>
      <div className="my-4 overflow-hidden rounded-xl border border-slate-700/60 bg-slate-950/80 shadow-xl backdrop-blur-md">
        {/* Header toolbar */}
        <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/80 px-4 py-2 text-xs text-slate-400">
          <div className="flex items-center gap-2 font-mono font-medium text-sky-400">
            <span className="h-2 w-2 rounded-full bg-sky-400 animate-pulse" />
            Diagrama de Arquitetura Interativo
          </div>

          <div className="flex items-center gap-1">
            <button
              onClick={handleZoomOut}
              title="Diminuir Zoom"
              className="rounded p-1.5 hover:bg-slate-800 hover:text-slate-200 transition-colors cursor-pointer"
            >
              <ZoomOut className="h-3.5 w-3.5" />
            </button>
            <span className="text-[10px] font-mono w-8 text-center">{Math.round(zoom * 100)}%</span>
            <button
              onClick={handleZoomIn}
              title="Aumentar Zoom"
              className="rounded p-1.5 hover:bg-slate-800 hover:text-slate-200 transition-colors cursor-pointer"
            >
              <ZoomIn className="h-3.5 w-3.5" />
            </button>
            <button
              onClick={handleResetZoom}
              title="Resetar Zoom"
              className="rounded p-1.5 hover:bg-slate-800 hover:text-slate-200 transition-colors cursor-pointer"
            >
              <RotateCcw className="h-3.5 w-3.5" />
            </button>
            <div className="h-3.5 w-[1px] bg-slate-800 mx-1" />
            <button
              onClick={handleCopyCode}
              title="Copiar Código Mermaid"
              className="flex items-center gap-1 rounded px-2 py-1 hover:bg-slate-800 hover:text-slate-200 transition-colors cursor-pointer"
            >
              {copied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
              <span>{copied ? 'Copiado' : 'Código'}</span>
            </button>
            <button
              onClick={() => setIsFullscreen(true)}
              title="Expandir Diagrama"
              className="rounded p-1.5 hover:bg-slate-800 hover:text-slate-200 transition-colors cursor-pointer"
            >
              <Maximize2 className="h-3.5 w-3.5" />
            </button>
          </div>
        </div>

        {/* Diagram Canvas */}
        <div
          ref={containerRef}
          className="relative min-h-[160px] overflow-auto p-6 flex items-center justify-center transition-all bg-gradient-to-b from-slate-950 to-slate-900/60"
        >
          {svgContent ? (
            <div
              style={{ transform: `scale(${zoom})`, transformOrigin: 'top center', transition: 'transform 0.15s ease-out' }}
              dangerouslySetInnerHTML={{ __html: svgContent }}
              className="mermaid max-w-full"
            />
          ) : (
            <div className="flex items-center gap-2 text-xs text-slate-500 animate-pulse">
              <div className="h-3 w-3 rounded-full border-2 border-sky-400 border-t-transparent animate-spin" />
              Renderizando diagrama vetorial...
            </div>
          )}
        </div>
      </div>

      {/* Fullscreen Modal View */}
      {isFullscreen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-md p-4">
          <div className="relative flex h-[90vh] w-[95vw] flex-col rounded-2xl border border-slate-700 bg-slate-950 shadow-2xl overflow-hidden">
            <div className="flex items-center justify-between border-b border-slate-800 bg-slate-900/90 px-6 py-3">
              <span className="font-semibold text-sm text-slate-200 flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-sky-400" />
                Visualização Expandida do Diagrama de Arquitetura
              </span>
              <button
                onClick={() => setIsFullscreen(false)}
                className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-800 hover:text-white cursor-pointer"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="flex-1 overflow-auto p-8 flex items-center justify-center">
              <div
                style={{ transform: `scale(${zoom * 1.3})`, transformOrigin: 'center' }}
                dangerouslySetInnerHTML={{ __html: svgContent }}
                className="mermaid"
              />
            </div>
          </div>
        </div>
      )}
    </>
  );
};
