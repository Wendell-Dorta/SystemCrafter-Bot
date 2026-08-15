package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/service"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/ai"
)

type toolServiceAdapter struct{}

// NewToolService creates an instance of ToolService.
func NewToolService() service.ToolService {
	return &toolServiceAdapter{}
}

func (s *toolServiceAdapter) GetToolDeclarations() interface{} {
	return []ai.FunctionDeclaration{
		{
			Name:        "estimate_cloud_costs",
			Description: "Calcula a estimativa detalhada e realista de custos de infraestrutura em nuvem (suportando AWS, GCP, Azure e Oracle Cloud OCI) com cálculo comparativo multi-cloud e estratégias de FinOps.",
			Parameters: ai.Schema{
				Type: "OBJECT",
				Properties: map[string]ai.SchemaProperty{
					"provider":            {Type: "STRING", Description: "Provedor de nuvem principal: 'AWS', 'GCP', 'AZURE', ou 'ORACLE' (OCI)."},
					"workloadTier":        {Type: "STRING", Description: "Porte da carga: 'Startup', 'Growth', ou 'Enterprise'."},
					"computeType":         {Type: "STRING", Description: "Tipo de computação: 'Serverless' (Lambda/CloudRun/ContainerApps/Functions), 'Kubernetes' (EKS/GKE/AKS/OKE), ou 'VMs'."},
					"computeArchitecture": {Type: "STRING", Description: "Arquitetura do processador: 'ARM' (Graviton/Ampere) ou 'x86'."},
					"databaseType":        {Type: "STRING", Description: "Tipo de banco de dados: 'PostgreSQL', 'MySQL', 'MongoDB', 'DynamoDB', 'Redis'."},
					"databaseRedundancy":  {Type: "STRING", Description: "Redundância do banco: 'Single-AZ', 'Multi-AZ', ou 'Multi-Region'."},
					"cacheEnabled":        {Type: "BOOLEAN", Description: "Se deve incluir camada de cache Redis gerenciado."},
					"storageGB":           {Type: "INTEGER", Description: "Volume de armazenamento em objetos e discos (GB)."},
					"monthlyReqsMil":      {Type: "NUMBER", Description: "Volume estimado de requisições mensais em milhões."},
					"dataEgressGB":        {Type: "INTEGER", Description: "Tráfego de saída de dados para a internet (GB/mês)."},
				},
				Required: []string{"provider", "workloadTier", "computeType", "databaseType"},
			},
		},
		{
			Name:        "audit_security_compliance",
			Description: "Realiza uma auditoria de segurança arquitetural, verificando riscos de Ponto Único de Falha (SPOF), conformidade com LGPD/GDPR, OWASP Top 10 e PCI-DSS.",
			Parameters: ai.Schema{
				Type: "OBJECT",
				Properties: map[string]ai.SchemaProperty{
					"architectureType":    {Type: "STRING", Description: "Tipo da arquitetura: 'Monolith', 'Microservices', 'Serverless', ou 'Event-Driven'."},
					"dataClassification":  {Type: "STRING", Description: "Sensibilidade dos dados: 'Public', 'Internal', 'Confidential', ou 'Financial_PII'."},
					"complianceStandards": {Type: "ARRAY", Items: &ai.SchemaProperty{Type: "STRING"}, Description: "Lista de normas a validar: ['LGPD', 'OWASP', 'PCI-DSS', 'SOC2']."},
					"hasPublicDatabase":   {Type: "BOOLEAN", Description: "Indica se o banco de dados possui IP público."},
					"hasMTLS":             {Type: "BOOLEAN", Description: "Indica se há autenticação mTLS / Zero-Trust."},
					"hasEncryptionAtRest": {Type: "BOOLEAN", Description: "Indica se os dados e backups são criptografados com KMS."},
					"hasRateLimiting":     {Type: "BOOLEAN", Description: "Indica se há controle de taxa no API Gateway."},
					"hasMultiAZ":          {Type: "BOOLEAN", Description: "Indica se a infraestrutura está em múltiplas Zonas de Disponibilidade."},
				},
				Required: []string{"architectureType", "dataClassification"},
			},
		},
		{
			Name:        "generate_tech_stack_matrix",
			Description: "Gera uma matriz comparativa de tecnologias e trade-offs arquiteturais com recomendações precisas.",
			Parameters: ai.Schema{
				Type: "OBJECT",
				Properties: map[string]ai.SchemaProperty{
					"workloadType":     {Type: "STRING", Description: "Tipo de carga: 'Realtime', 'CRUD', 'HighThroughput', 'DataAnalytics', ou 'AI_Agent'."},
					"targetScale":      {Type: "STRING", Description: "Escala pretendida: 'MVP', 'Growth', ou 'Hyperscale'."},
					"teamSkillPrimary": {Type: "STRING", Description: "Linguagem principal da equipe: 'Go', 'TypeScript', 'Python', 'Java', ou 'CSharp'."},
				},
				Required: []string{"workloadType"},
			},
		},
		{
			Name:        "lookup_architecture_patterns",
			Description: "Consulta o catálogo de padrões de arquitetura de software de referência (Outbox, CQRS, Saga, Event-Driven) com diagramas Mermaid.",
			Parameters: ai.Schema{
				Type: "OBJECT",
				Properties: map[string]ai.SchemaProperty{
					"patternName": {Type: "STRING", Description: "Nome do padrão: 'outbox', 'cqrs', 'saga', ou 'event_driven'."},
				},
				Required: []string{"patternName"},
			},
		},
	}
}

func (s *toolServiceAdapter) ExecuteTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, int64, error) {
	start := time.Now()
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal tool args: %w", err)
	}

	var result interface{}
	var toolErr error

	switch name {
	case "estimate_cloud_costs":
		var params EstimateCloudCostParams
		if err := json.Unmarshal(argsJSON, &params); err != nil {
			toolErr = fmt.Errorf("invalid args: %w", err)
		} else {
			result = EstimateCloudCosts(params)
		}

	case "audit_security_compliance":
		var params AuditSecurityParams
		if err := json.Unmarshal(argsJSON, &params); err != nil {
			toolErr = fmt.Errorf("invalid args: %w", err)
		} else {
			result = AuditSecurityAndCompliance(params)
		}

	case "generate_tech_stack_matrix":
		var params GenerateTechStackMatrixParams
		if err := json.Unmarshal(argsJSON, &params); err != nil {
			toolErr = fmt.Errorf("invalid args: %w", err)
		} else {
			result = GenerateTechStackMatrix(params)
		}

	case "lookup_architecture_patterns":
		var p struct {
			PatternName string `json:"patternName"`
		}
		if err := json.Unmarshal(argsJSON, &p); err != nil {
			toolErr = fmt.Errorf("invalid args: %w", err)
		} else {
			result = LookupArchitecturePattern(p.PatternName)
		}

	default:
		toolErr = fmt.Errorf("unknown tool name: %s", name)
	}

	duration := time.Since(start).Milliseconds()
	return result, duration, toolErr
}
