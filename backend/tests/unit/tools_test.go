package unit_test

import (
	"context"
	"testing"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
)

func TestEstimateCloudCosts(t *testing.T) {
	est := tools.EstimateCloudCosts(tools.EstimateCloudCostParams{
		Provider:       "AWS",
		WorkloadTier:   "Growth",
		ComputeType:    "Kubernetes",
		DatabaseType:   "Postgres",
		CacheEnabled:   true,
		StorageGB:      250,
		MonthlyReqsMil: 10.0,
	})

	if est.MonthlyTotalUSD <= 0 {
		t.Fatalf("Expected total cost > 0, got %f", est.MonthlyTotalUSD)
	}

	if est.Provider != "AWS" {
		t.Errorf("Expected provider AWS, got %s", est.Provider)
	}
}

func TestAuditSecurity(t *testing.T) {
	audit := tools.AuditSecurityAndCompliance(tools.AuditSecurityParams{
		ArchitectureType:    "Microservices",
		DataClassification:  "Financial_PII",
		ComplianceStandards: []string{"LGPD", "PCI-DSS"},
		HasPublicDatabase:   true,
		HasMTLS:             false,
		HasEncryptionAtRest: false,
		HasRateLimiting:     false,
		HasMultiAZ:          false,
	})

	if audit.RiskScore < 50 {
		t.Errorf("Expected critical risk score >= 50, got %d", audit.RiskScore)
	}

	if len(audit.Vulnerabilities) == 0 {
		t.Errorf("Expected vulnerabilities to be flagged")
	}
}

func TestToolService_ExecuteTool(t *testing.T) {
	svc := tools.NewToolService()
	res, duration, err := svc.ExecuteTool(context.Background(), "lookup_architecture_patterns", map[string]interface{}{
		"patternName": "outbox",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("Expected non-nil result")
	}
	if duration < 0 {
		t.Fatalf("Expected duration >= 0, got %d", duration)
	}
}
