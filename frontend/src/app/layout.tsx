import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'SystemCrafter AI — Autonomous Software Architect & Cloud Designer',
  description:
    'Copiloto Inteligente de Arquitetura de Software e Design de Sistemas em Nuvem. Integração com Google Gemini 2.0 Flash, Backend em Go (Clean Architecture) e Diagramas Mermaid.js.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR" className="dark">
      <body className="bg-[#0b0f19] text-slate-100 antialiased selection:bg-sky-500/30 selection:text-sky-200">
        {children}
      </body>
    </html>
  );
}
