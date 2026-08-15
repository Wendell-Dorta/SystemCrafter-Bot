package memory

import "github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"

// GetPreloadedTemplates returns rich, pre-configured architectural templates.
func GetPreloadedTemplates() []entity.ArchitectureTemplate {
	return []entity.ArchitectureTemplate{
		{
			ID:          "ecommerce-event-driven",
			Title:       "High-Scale Event-Driven E-Commerce",
			Category:    "E-Commerce / Retail",
			Description: "Resilient microservices architecture with CQRS, Event Sourcing, Outbox Pattern, and Redis caching for handling black friday traffic spikes.",
			Tags:        []string{"Microservices", "Event-Driven", "CQRS", "Outbox Pattern", "Redis", "PostgreSQL", "Kafka/RabbitMQ"},
			Complexity:  "High",
			RecommendedStack: []string{
				"Frontend: Next.js (App Router, Edge Caching)",
				"Backend: Go (High-throughput checkout) + Node.js (Catalog)",
				"Database: PostgreSQL (ACID Orders) + MongoDB (Product Catalog)",
				"Messaging: Kafka / Cloud PubSub (Outbox Worker)",
				"Cache: Redis Cluster (Sessions & Carts)",
			},
			MermaidDiagram: `graph TD
    Client[Next.js Client] --> Cloudflare[Cloudflare CDN / WAF]
    Cloudflare --> APIGateway[API Gateway / Envoy]
    
    APIGateway --> AuthSvc[Auth Service - Go]
    APIGateway --> CatalogSvc[Catalog Service - Node.js]
    APIGateway --> OrderSvc[Order Service - Go]
    APIGateway --> PaymentSvc[Payment Service - Go]
    
    CatalogSvc --> RedisCache[(Redis Cluster)]
    CatalogSvc --> CatalogDB[(MongoDB Read Replica)]
    
    OrderSvc --> OrderDB[(PostgreSQL Main)]
    OrderSvc --> OutboxTable[(Outbox Table)]
    OutboxTable --> Debezium[Debezium / CDC]
    Debezium --> KafkaTopic{{Kafka: order-events}}
    
    KafkaTopic --> PaymentSvc
    KafkaTopic --> InventorySvc[Inventory Service]
    KafkaTopic --> NotificationSvc[Notification Service - Push/Email]`,
			EstimatedCost: "$450 - $1,200 / mês (Workload médio a alto)",
			KeyTradeoffs: []string{
				"Eventual consistency between orders and inventory requires idempotency keys.",
				"Higher operational overhead for Kafka/CDC compared to monolithic setups.",
				"Sub-millisecond catalog reads via Redis layer.",
			},
			PromptStarter: "Gostaria de avaliar a arquitetura de E-Commerce Event-Driven com foco em suportar 50.000 requisições/minuto na Black Friday. Quais os pontos de falha e estimativa de custo na AWS?",
		},
		{
			ID:          "fintech-payment-ledger",
			Title:       "Fintech Real-Time Payment Gateway & Immutable Ledger",
			Category:    "Fintech / Banking",
			Description: "Ultra-secure, double-entry bookkeeping ledger with distributed locks, strict idempotency, and PCI-DSS / LGPD compliance.",
			Tags:        []string{"Fintech", "Ledger", "Idempotency", "PostgreSQL", "PCI-DSS", "Zero-Trust", "Go"},
			Complexity:  "Very High",
			RecommendedStack: []string{
				"Backend: Go (Strict typing, deterministic memory, low latency)",
				"Database: PostgreSQL with Row-Level Locking & Serializable Isolation",
				"Key Management: AWS KMS / HashiCorp Vault (Envelope Encryption)",
				"Observability: OpenTelemetry + Prometheus + Grafana",
				"Network: Private VPC + Mutual TLS (mTLS)",
			},
			MermaidDiagram: `graph TD
    App[Mobile/Web App] --> WAF[AWS WAF / Shield]
    WAF --> APIGW[API Gateway - mTLS]
    
    APIGW --> IdempotencyFilter[Idempotency Check - Redis TTL]
    IdempotencyFilter --> PaymentCore[Payment Engine - Go]
    
    PaymentCore --> KMS[AWS KMS / Vault - Tokenization]
    PaymentCore --> AntiFraud[Anti-Fraud ML Scorer]
    
    PaymentCore --> DBTransaction[Serializable Transaction]
    subgraph ACID Storage
        DBTransaction --> LedgerTable[(Double-Entry Ledger)]
        DBTransaction --> AuditLog[(Immutable Audit Log - WORM)]
    end
    
    PaymentCore --> WebhookWorker[Webhook Dispatcher - Worker Pool]
    WebhookWorker --> Merchant[Merchant Endpoint]`,
			EstimatedCost: "$600 - $1,800 / mês (Alta redundância Multi-AZ)",
			KeyTradeoffs: []string{
				"Serializable isolation reduces transaction throughput in favor of zero race conditions.",
				"Strict cryptographic signing adds ~15ms latency per request.",
				"Zero-loss SLA with synchronous Multi-AZ DB replication.",
			},
			PromptStarter: "Poderia revisar os requisitos de resiliência e conformidade de uma arquitetura de Gateway de Pagamentos e Ledger Imutável em Go? Preciso garantir conformidade com PCI-DSS e proteção contra double-spending.",
		},
		{
			ID:          "ai-saas-rag-engine",
			Title:       "Modern AI-Powered SaaS with RAG & Multimodal Engine",
			Category:    "AI / SaaS",
			Description: "Full-stack AI SaaS with vector similarity search, document parsing pipeline, background embeddings, and Gemini 2.0 streaming.",
			Tags:        []string{"Generative AI", "RAG", "Vector DB", "Gemini", "Go", "Next.js", "Async Workers"},
			Complexity:  "Intermediate",
			RecommendedStack: []string{
				"Frontend: Next.js 15 (React 19, Server-Sent Events, TailwindCSS)",
				"Backend API: Go (REST, SSE Streaming, Rate Limiting)",
				"Vector Storage: PostgreSQL with pgvector or Qdrant",
				"Object Storage: S3 / Cloud Storage (Documents & OCR assets)",
				"AI Core: Google Gemini 2.0 Flash + Text Embeddings",
			},
			MermaidDiagram: `graph TD
    User[Web Client - Next.js] --> Edge[Cloudflare Edge / SSE]
    Edge --> GoBackend[Go Agentic Core]
    
    GoBackend --> RateLimiter[In-Memory Token Bucket Limiter]
    RateLimiter --> Cache[Thread-Safe In-Memory Cache]
    
    subgraph Ingestion Pipeline
        User -->|Upload PDF/Docs| S3Bucket[(S3 Object Storage)]
        S3Bucket --> WorkerPool[Go Document Parser Worker]
        WorkerPool --> GeminiEmbed[Gemini Embedding API]
        GeminiEmbed --> VectorDB[(PGVector / Qdrant)]
    end
    
    subgraph Retrieval & Inference
        GoBackend --> VectorDB
        GoBackend --> GeminiLLM[Google Gemini 2.0 Flash]
        GeminiLLM -->|SSE Stream| GoBackend
        GoBackend -->|SSE Stream Tokens| User
    end`,
			EstimatedCost: "$120 - $350 / mês (Baseado em tokens e compute serverless)",
			KeyTradeoffs: []string{
				"pgvector keeps tech stack unified in PostgreSQL without dedicated vector infra.",
				"Async embedding processing prevents request timeouts on large file uploads.",
			},
			PromptStarter: "Quero projetar um SaaS de IA com RAG e busca vetorial. Como estruturar o backend em Go para processar documentos em background e fazer streaming via SSE?",
		},
		{
			ID:          "realtime-iot-telemetry",
			Title:       "High-Throughput IoT Telemetry & Analytics Platform",
			Category:    "IoT / Big Data",
			Description: "Ingestion pipeline capable of processing 100,000+ telemetry events/sec from smart devices with time-series analysis and anomaly detection.",
			Tags:        []string{"IoT", "Time-Series", "MQTT", "ClickHouse", "TimescaleDB", "Go", "WebSockets"},
			Complexity:  "High",
			RecommendedStack: []string{
				"Ingestion Gateway: Go MQTT Broker / EMQX",
				"Stream Processing: Go Worker Pool with Ring Buffers",
				"Time-Series DB: ClickHouse (Analytics) + TimescaleDB (Metrics)",
				"Realtime Dashboard: Next.js with WebSocket Feeds",
			},
			MermaidDiagram: `graph TD
    Sensors[100k IoT Devices] -->|MQTT / gRPC| LB[Network Load Balancer]
    LB --> Ingestion[Go Ingestion Nodes]
    
    Ingestion --> RingBuffer[Go In-Memory Ring Buffer]
    RingBuffer --> BatchWriter[Batch Bulk Inserter]
    
    BatchWriter --> ClickHouse[(ClickHouse Columnar DB)]
    Ingestion --> StreamAlert[Anomaly Detection Filter]
    StreamAlert --> WebSocketSvc[WebSocket Server - Go]
    WebSocketSvc --> LiveDashboard[Next.js Realtime Dashboard]`,
			EstimatedCost: "$300 - $800 / mês",
			KeyTradeoffs: []string{
				"Micro-batching inserts (e.g. every 500ms) optimizes ClickHouse compression.",
				"WebSocket fan-out handled via Go goroutines with low memory footprint.",
			},
			PromptStarter: "Como desenhar uma arquitetura de telemetria IoT em tempo real com Go e ClickHouse para aguentar 100k requisições/segundo com baixo custo?",
		},
	}
}
