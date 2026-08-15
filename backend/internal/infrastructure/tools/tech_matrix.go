package tools

import (
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// GenerateTechStackMatrixParams parameters for matrix comparison.
type GenerateTechStackMatrixParams struct {
	WorkloadType     string `json:"workloadType"`
	TargetScale      string `json:"targetScale"`
	TeamSkillPrimary string `json:"teamSkillPrimary"`
}

// GenerateTechStackMatrix produces an architectural tradeoff matrix.
func GenerateTechStackMatrix(params GenerateTechStackMatrixParams) entity.TechStackMatrix {
	workload := strings.ToLower(strings.TrimSpace(params.WorkloadType))
	if workload == "" {
		workload = "highthroughput"
	}

	recommended := make(map[string]string)
	tradeoffs := make(map[string]string)
	alternatives := make(map[string]string)

	switch workload {
	case "realtime", "chat", "iot":
		recommended["Frontend"] = "Next.js 15 (React 19 + Server-Sent Events / WebSockets)"
		recommended["Backend"] = "Go (Goroutines nativas com baixo overhead de memória por conexão)"
		recommended["Database"] = "PostgreSQL (Metadados) + ClickHouse / ScyllaDB (Eventos em tempo real)"
		recommended["Cache_PubSub"] = "Redis Streams / Dragonfly (Sub-milissegundo)"
		recommended["Protocol"] = "HTTP/2 SSE para streaming unidirecional ou WebSockets para duplex"

		tradeoffs["Go vs Node.js"] = "Go consome ~10x menos memória RAM em concorrência pesada (100k conexões simultâneas)."
		tradeoffs["Postgres vs DynamoDB"] = "Postgres oferece flexibilidade relacional e ACID completo; DynamoDB escala horizontalmente sem esforço, mas custa mais em queries complexas."

		alternatives["MessageBroker"] = "RabbitMQ (AMQP) para roteamento flexível de tópicos."

	case "ai_agent", "llm", "rag":
		recommended["Frontend"] = "Next.js 15 com App Router e Streaming UI"
		recommended["Backend"] = "Go (Gateway de alta performance, SSE e rate-limiting) ou Python (FastAPI/LangChain para ML pipelines)"
		recommended["VectorDB"] = "PostgreSQL com extensão pgvector (unifica dados relacionais e vetores)"
		recommended["LLM_Provider"] = "Google Gemini 2.0 Flash (Janela de contexto gigante de 1M tokens + Multimodalidade nativa)"
		recommended["Storage"] = "AWS S3 / Cloud Storage para arquivos de referência e PDFs"

		tradeoffs["pgvector vs Pinecone"] = "pgvector reduz custos e complexidade mantendo tudo no Postgres; Pinecone é puramente gerenciado com busca vetorial otimizada para bilhões de vetores."
		tradeoffs["Gemini vs GPT-4o"] = "Gemini 2.0 Flash possui latência extremamente baixa de primeiro token (TTFT) e custo-benefício imbatível para streaming interativo."

		alternatives["VectorDB"] = "Qdrant / Milvus para grafos de conhecimento de ultra-escala."

	default: // CRUD / General
		recommended["Frontend"] = "Next.js com TailwindCSS e Shadcn UI"
		recommended["Backend"] = "Go com arquitetura limpa (Clean Architecture / Hexagonal)"
		recommended["Database"] = "PostgreSQL com Read Replicas e Connection Pooling (PgBouncer)"
		recommended["Cache"] = "Redis / In-Memory Token-Bucket para taxa de requisições"
		recommended["CI_CD"] = "GitHub Actions + Docker + Deploy Serverless / Kubernetes"

		tradeoffs["Monólito vs Microsserviços"] = "Inicie com Monólito Modular em Go para alta velocidade de entrega; extraia microsserviços apenas quando gargalos de time/domínio surgirem."
		tradeoffs["REST vs gRPC"] = "REST/JSON para clientes web públicos; gRPC com Protobuf para comunicação interna inter-serviços."

		alternatives["Database"] = "CockroachDB para distribuição geográfica global."
	}

	return entity.TechStackMatrix{
		WorkloadType:  params.WorkloadType,
		Recommended:   recommended,
		TradeoffNotes: tradeoffs,
		Alternatives:  alternatives,
	}
}
