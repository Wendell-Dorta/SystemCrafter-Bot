# 🚀 SystemCrafter AI — Autonomous Software Architect & Cloud Designer

> **Copiloto Inteligente de Arquitetura de Software e Design de Sistemas em Nuvem** com suporte a **Visão Multimodal**, **Function Calling no Go (Tools)**, **Streaming em Tempo Real (SSE)**, **Renderização Dinâmica de Diagramas Mermaid.js** e **Túnel Público Gratuito Cloudflare**.

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-15%20App%20Router-black?style=flat&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Multi--stage-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Cloudflare](https://img.shields.io/badge/Cloudflare-Free%20Tunnel-F38020?style=flat&logo=cloudflare)](https://www.cloudflare.com/)
[![Gemini](https://img.shields.io/badge/AI-Gemini%202.0%20Flash-4285F4?style=flat&logo=google)](https://ai.google.dev/)
[![Clean Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Hexagonal-green?style=flat)]()
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 🌟 Principais Funcionalidades

1. **🤖 Ciclo Agêntico com Function Calling (Go Tools)**:
   - 💰 `estimate_cloud_costs`: Calcula custos detalhados em AWS, GCP ou Azure (Compute, DB, Cache, Storage e CDN).
   - 🔒 `audit_security_compliance`: Identifica vulnerabilidades, SPOF (Single Point of Failure) e conformidade com **LGPD, PCI-DSS e OWASP**.
   - ⚖️ `generate_tech_stack_matrix`: Produz matrizes de trade-offs entre linguagens, bancos e brokers de mensagens.
   - 📐 `lookup_architecture_patterns`: Consulta padrões de referência (*Transactional Outbox, CQRS, Saga Pattern, Event-Driven*).
2. **📐 Renderizador de Diagramas Mermaid Interativo**:
   - Transcreve arquiteturas em diagramas vetoriais renderizados on-the-fly com controles de Zoom In/Out, Reset, Cópia de Código e Modal Fullscreen.
3. **🖼️ Multimodalidade & Análise de Quadros Brancos (Gemini Vision)**:
   - Permite anexar prints de diagramas ou fotos de quadros brancos desenhados à mão para análise arquitetural instantânea.
4. **⚡ Streaming SSE (Server-Sent Events)**:
   - Respostas token a token com feedback visual em tempo real de ativação de ferramentas.
5. **☁️ Túnel Público Gratuito Cloudflare (TryCloudflare)**:
   - Gera automaticamente uma URL pública HTTPS segura (`https://*.trycloudflare.com`) sem necessidade de configurar portas no roteador ou contas pagas.
6. **🛡️ Segurança & Performance de Produção**:
   - **Rate Limiting por IP (Token Bucket)** integrado.
   - **Cache em Memória thread-safe com TTL** (sem necessidade de infraestrutura pesada externa).
   - **Headers de Segurança e CORS configurado**.
   - **Graceful Shutdown** e recuperação contra panics.
7. **🏛️ Clean Architecture & SOLID Estrito**:
   - 100% de separação entre Domínio, Casos de Uso, Infraestrutura e Apresentação, com suíte de testes isolada.

---

## 🏗️ Arquitetura do Projeto

```
SystemCrafter-Bot/
├── backend/                         # 🐹 Backend em Go (Clean Architecture Rígida)
│   ├── cmd/server/main.go           # Composition Root & Injeção de Dependências
│   ├── internal/                    # 100% Código Puro de Produção
│   │   ├── domain/                  # Entidades e Interfaces (Ports)
│   │   ├── usecase/                 # Casos de Uso (SOLID)
│   │   ├── infrastructure/          # Adaptadores (Gemini AI, Memory Cache, Tools)
│   │   └── transport/http/          # Handlers HTTP, Middlewares e Router
│   ├── tests/                       # 🧪 Suíte de Testes Isolada (Unit & Integration)
│   └── Dockerfile                   # Multi-stage Docker build para Go (< 25MB)
│
├── frontend/                        # 🌐 Next.js 15 (React 19, TypeScript, TailwindCSS)
│   ├── src/
│   │   ├── app/                     # App Router, Layout e Estilos
│   │   ├── components/              # ChatInterface, MermaidViewer, ToolCards, Header, Sidebar
│   │   ├── hooks/                   # useChatStream (SSE), useHealth
│   │   ├── lib/                     # API Client e Utilitários
│   │   └── types/                   # Tipos TypeScript
│   └── Dockerfile                   # Multi-stage Standalone Docker build para Next.js
│
├── nginx/
│   └── nginx.conf                   # Proxy reverso otimizado para streaming SSE
├── docker-compose.yml               # Orquestração completa de containers + Cloudflare Tunnel
├── .env.example                     # Modelo único de variáveis de ambiente na raiz
├── LICENSE                          # Licença MIT
└── README.md                        # Documentação completa
```

---

## 🐳 Executando com Docker Compose + Cloudflare Tunnel

```bash
# 1. Clone o repositório
git clone https://github.com/Wendell-Dorta/SystemCrafter-Bot.git
cd SystemCrafter-Bot

# 2. Crie o arquivo .env a partir do modelo na raiz:
cp .env.example .env
# Edite o .env colocando sua chave: GEMINI_API_KEY=AIzaSy...

# 3. Suba todos os containers
docker compose up --build -d
```

### 🌍 Como pegar a URL Pública Gratuita do Cloudflare para enviar aos amigos:

Basta rodar o comando:
```bash
docker logs systemcrafter-tunnel
```
> Você verá uma linha com o link HTTPS gerado automaticamente, por exemplo:
> `https://swift-brave-design.trycloudflare.com`

Envie essa URL para qualquer pessoa e ela poderá usar o chatbot imediatamente com HTTPS e certificado SSL gratuito!

---

### Portas Locais Disponíveis:
* **Gateway Unificado (Nginx)**: `http://localhost` (Porta 80)
* **Frontend Direto (Next.js)**: `http://localhost:3000`
* **Backend Direto (Go API)**: `http://localhost:8080`

---

## 💻 Executando Localmente (Sem Docker)

### 1. Backend (Go)
```bash
cd backend
go run cmd/server/main.go
```
*API rodando em `http://localhost:8080`*

### 2. Frontend (Next.js)
```bash
cd frontend
npm run dev
```
*Interface rodando em `http://localhost:3000`*

---

## 🧪 Como Rodar os Testes Automatizados

### Testes do Backend (Go)
```bash
cd backend
go test -v ./...
```

### Build do Frontend (Next.js)
```bash
cd frontend
npm run build
```

---

## 📚 Endpoints da API

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/chat/stream` | Streaming SSE em tempo real com eventos de ferramentas (`tool_start`, `tool_result`, `token`, `done`) |
| `POST` | `/api/chat` | Chat síncrono com cache automático em memória |
| `GET` | `/api/templates` | Lista blueprints pré-cadastrados (E-Commerce, Fintech, AI SaaS, IoT) |
| `GET` | `/api/templates/{id}` | Retorna detalhes e diagrama Mermaid de um blueprint específico |
| `GET` | `/api/sessions` | Lista todas as sessões e histórico de conversas gravadas |
| `GET` | `/health` | Telemetria do sistema, status do rate limit e taxa de cache hits |

---

## 📄 Licença

Este projeto está sob a licença **MIT** — veja o arquivo [LICENSE](LICENSE) para mais detalhes.
