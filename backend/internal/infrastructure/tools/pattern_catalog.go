package tools

import (
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// LookupArchitecturePattern returns structured details and Mermaid diagrams for a given pattern.
func LookupArchitecturePattern(patternName string) entity.ArchitecturePattern {
	norm := strings.ToLower(strings.TrimSpace(patternName))

	switch {
	case strings.Contains(norm, "outbox"):
		return entity.ArchitecturePattern{
			Name:     "Transactional Outbox Pattern",
			Category: "Data Consistency & Messaging",
			Summary:  "Garante consistência eventual atômica entre a alteração no banco de dados e a publicação de eventos no broker de mensagens.",
			WhenToUse: []string{
				"Quando um microsserviço precisa persistir dados em banco e disparar evento para Kafka/RabbitMQ sem risco de perder mensagens.",
				"Para evitar 'Dual Write' inconsistente.",
			},
			Pitfalls: []string{
				"Consumidores devem ser idempotentes.",
				"Necessita de processo de polling ou Change Data Capture (Debezium).",
			},
			MermaidDiagram: `sequenceDiagram
    participant App as Aplicação (Go)
    participant DB as Banco de Dados (PostgreSQL)
    participant Relay as Message Relay / CDC
    participant Broker as Message Broker (Kafka)

    App->>DB: BEGIN Transaction
    App->>DB: INSERT into orders (status: 'CREATED')
    App->>DB: INSERT into outbox_events (type: 'ORDER_CREATED')
    App->>DB: COMMIT Transaction
    
    Relay->>DB: Lê eventos pendentes
    Relay->>Broker: Publica mensagem
    Relay->>DB: Marca como processado`,
		}

	case strings.Contains(norm, "cqrs"):
		return entity.ArchitecturePattern{
			Name:     "CQRS (Command Query Responsibility Segregation)",
			Category: "Data Architecture",
			Summary:  "Separa os modelos de escrita dos modelos de leitura otimizados.",
			WhenToUse: []string{
				"Aplicações com assimetria severa de tráfego (ex: 95% leitura / 5% escrita).",
			},
			Pitfalls: []string{
				"Complexidade de sincronização entre bancos.",
			},
			MermaidDiagram: `graph LR
    User[Cliente Web] -->|Command| WriteAPI[Write Model]
    WriteAPI -->|ACID| WriteDB[(PostgreSQL Master)]
    WriteDB -->|CDC| SyncWorker[Sync Engine]
    SyncWorker -->|Read Model| ReadDB[(Elasticsearch / Redis)]
    User -->|Query| ReadAPI[Read Model]
    ReadAPI --> ReadDB`,
		}

	case strings.Contains(norm, "saga"):
		return entity.ArchitecturePattern{
			Name:     "Saga Pattern",
			Category: "Distributed Transactions",
			Summary:  "Gerencia transações distribuídas através de sequência de transações locais e compensações.",
			WhenToUse: []string{
				"Processamento de pedidos em múltiplos microsserviços.",
			},
			Pitfalls: []string{
				"Complexidade de rollback compensatório.",
			},
			MermaidDiagram: `graph TD
    Client[Cliente] --> Orchestrator[Saga Orchestrator]
    Orchestrator -->|1. Reservar| StockSvc[Stock Service]
    StockSvc -->|Sucesso| Orchestrator
    Orchestrator -->|2. Cobrar| PaymentSvc[Payment Service]
    PaymentSvc -->|Falha| Orchestrator
    Orchestrator -.->|Compensar| StockSvc`,
		}

	default:
		return entity.ArchitecturePattern{
			Name:     "Event-Driven Microservices Architecture",
			Category: "Asynchronous Systems",
			Summary:  "Arquitetura desacoplada baseada na emissão e consumo de eventos assíncronos.",
			WhenToUse: []string{
				"Sistemas que exigem escalabilidade horizontal desacoplada.",
			},
			Pitfalls: []string{
				"Rastreabilidade exige Distributed Tracing (OpenTelemetry).",
			},
			MermaidDiagram: `graph LR
    Producer[Produtor] -->|Publica| Broker{{Event Broker}}
    Broker -->|Assina| ConsumerA[Serviço A]
    Broker -->|Assina| ConsumerB[Serviço B]`,
		}
	}
}
