package tools

import (
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// AuditSecurityParams contains input parameters for the security auditor.
type AuditSecurityParams struct {
	ArchitectureType    string   `json:"architectureType"`
	DataClassification  string   `json:"dataClassification"`
	ComplianceStandards []string `json:"complianceStandards"`
	HasPublicDatabase   bool     `json:"hasPublicDatabase"`
	HasMTLS             bool     `json:"hasMTLS"`
	HasEncryptionAtRest bool     `json:"hasEncryptionAtRest"`
	HasRateLimiting     bool     `json:"hasRateLimiting"`
	HasMultiAZ          bool     `json:"hasMultiAZ"`
}

// AuditSecurityAndCompliance performs security and SPOF analysis.
func AuditSecurityAndCompliance(params AuditSecurityParams) entity.SecurityAuditResult {
	score := 10
	findings := make([]entity.SecurityFinding, 0)
	complianceMap := make(map[string]string)
	recommendations := make([]string, 0)

	// 1. Check Public Database
	if params.HasPublicDatabase {
		score += 35
		findings = append(findings, entity.SecurityFinding{
			ID:          "SEC-CRIT-001",
			Severity:    "CRITICAL",
			Category:    "Network",
			Description: "Banco de dados exposto diretamente na internet pública.",
			Remediation: "Mover o banco para Subnets Privadas / VPC Isolada acessível somente via VPN / mTLS.",
		})
	}

	// 2. Encryption at Rest
	if !params.HasEncryptionAtRest {
		score += 20
		findings = append(findings, entity.SecurityFinding{
			ID:          "SEC-HIGH-002",
			Severity:    "HIGH",
			Category:    "Encryption",
			Description: "Falta de criptografia em repouso nos volumes de dados.",
			Remediation: "Habilitar criptografia KMS (AES-256) em todos os discos e buckets.",
		})
	}

	// 3. mTLS
	archType := strings.ToLower(params.ArchitectureType)
	if (archType == "microservices" || archType == "event-driven") && !params.HasMTLS {
		score += 15
		findings = append(findings, entity.SecurityFinding{
			ID:          "SEC-MED-003",
			Severity:    "MEDIUM",
			Category:    "Auth",
			Description: "Comunicação entre serviços sem autenticação mTLS (Zero Trust).",
			Remediation: "Implementar mTLS via Service Mesh ou tokens JWT assimétricos.",
		})
	}

	// 4. Rate Limiting
	if !params.HasRateLimiting {
		score += 15
		findings = append(findings, entity.SecurityFinding{
			ID:          "SEC-MED-004",
			Severity:    "MEDIUM",
			Category:    "Availability",
			Description: "Ausência de controle de vazão (Rate Limiting) na borda da API.",
			Remediation: "Adicionar rate limiting por IP/Token no API Gateway.",
		})
	}

	// 5. Single Point of Failure
	if !params.HasMultiAZ {
		score += 15
		findings = append(findings, entity.SecurityFinding{
			ID:          "SPOF-HIGH-005",
			Severity:    "HIGH",
			Category:    "SPOF",
			Description: "Infraestrutura concentrada em uma única Zona de Disponibilidade (Single-AZ).",
			Remediation: "Configurar deploys Multi-AZ com Load Balancers gerenciados.",
		})
	}

	if score > 100 {
		score = 100
	}

	riskLevel := "LOW"
	if score >= 70 {
		riskLevel = "CRITICAL"
	} else if score >= 45 {
		riskLevel = "HIGH"
	} else if score >= 25 {
		riskLevel = "MEDIUM"
	}

	for _, std := range params.ComplianceStandards {
		stdNorm := strings.ToUpper(strings.TrimSpace(std))
		switch stdNorm {
		case "LGPD", "GDPR":
			if params.HasEncryptionAtRest && !params.HasPublicDatabase {
				complianceMap[stdNorm] = "EM CONFORMIDADE: Dados cifrados e isolados da internet pública."
			} else {
				complianceMap[stdNorm] = "NÃO CONFORME: Requer criptografia (Art. 46 LGPD) e trilha de auditoria."
			}
		case "PCI-DSS":
			if params.HasEncryptionAtRest && params.HasMTLS && !params.HasPublicDatabase {
				complianceMap[stdNorm] = "EM CONFORMIDADE: Segmentação de rede e criptografia ativadas."
			} else {
				complianceMap[stdNorm] = "NÃO CONFORME: Requer tokenização de PAN e mTLS."
			}
		case "OWASP":
			complianceMap[stdNorm] = "PARCIALMENTE CONFORME: Necessário sanitização de input e rate limiting."
		}
	}

	recommendations = append(recommendations,
		"Implementar política de Least Privilege (IAM Roles granulares).",
		"Centralizar logs com mascaramento automático de PII.",
		"Executar varredura periódica de vulnerabilidades (SAST/DAST).",
	)

	return entity.SecurityAuditResult{
		RiskScore:        score,
		RiskLevel:        riskLevel,
		Vulnerabilities:  findings,
		ComplianceStatus: complianceMap,
		Recommendations:  recommendations,
	}
}
