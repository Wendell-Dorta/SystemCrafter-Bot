package unit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/config"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/handler"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

func TestConfigLoading(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("GEMINI_API_KEY", "test-key-123")
	os.Setenv("GEMINI_MODEL", "gemini-2.5-flash")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://app.example.com")

	cfg := config.LoadConfig()
	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}
	if cfg.GeminiAPIKey != "test-key-123" {
		t.Errorf("Expected api key test-key-123, got %s", cfg.GeminiAPIKey)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Errorf("Expected 2 allowed origins, got %d", len(cfg.AllowedOrigins))
	}
}

func TestTools_Comprehensive(t *testing.T) {
	// 1. Cost Estimator across all providers and tiers
	providers := []string{"AWS", "GCP", "Azure", "OCI", "Oracle"}
	tiers := []string{"Starter", "Growth", "Enterprise", "High"}

	for _, prov := range providers {
		for _, tier := range tiers {
			res := tools.EstimateCloudCosts(tools.EstimateCloudCostParams{
				Provider:       prov,
				WorkloadTier:   tier,
				ComputeType:    "Kubernetes",
				DatabaseType:   "Postgres",
				CacheEnabled:   true,
				StorageGB:      500,
				MonthlyReqsMil: 25.0,
			})
			if res.MonthlyTotalUSD <= 0 {
				t.Errorf("Expected positive monthly cost for %s / %s, got %v", prov, tier, res)
			}
		}
	}

	// 2. Pattern Catalog
	patterns := []string{"outbox", "cqrs", "saga", "event-sourcing", "circuit-breaker", "unknown-pattern"}
	for _, p := range patterns {
		res := tools.LookupArchitecturePattern(p)
		if res.Name == "" {
			t.Errorf("Expected non-empty pattern name for %s", p)
		}
	}

	// 3. Security & Compliance Auditor
	comps := []string{"LGPD", "PCI-DSS", "HIPAA", "SOC2"}
	for _, c := range comps {
		res := tools.AuditSecurityAndCompliance(tools.AuditSecurityParams{
			ArchitectureType:    "Microservices",
			DataClassification:  "Financial_PII",
			ComplianceStandards: []string{c},
			HasPublicDatabase:   false,
			HasMTLS:             true,
			HasEncryptionAtRest: true,
			HasRateLimiting:     true,
			HasMultiAZ:          true,
		})
		if res.RiskLevel == "" {
			t.Errorf("Expected risk level for %s", c)
		}
	}

	// 4. Tech Stack Matrix
	workloads := []string{"realtime", "ai_agent", "fintech", "ecommerce", "iot", "highthroughput"}
	for _, w := range workloads {
		matrix := tools.GenerateTechStackMatrix(tools.GenerateTechStackMatrixParams{
			WorkloadType:     w,
			TargetScale:      "High",
			TeamSkillPrimary: "Go",
		})
		if len(matrix.Recommended) == 0 {
			t.Errorf("Expected recommended tech stack for %s", w)
		}
	}

	// 5. Tool Service declarations & executions
	toolSvc := tools.NewToolService()
	decls := toolSvc.GetToolDeclarations()
	if decls == nil {
		t.Error("Expected non-nil tool declarations")
	}

	// Test executing valid tool
	res, dur, err := toolSvc.ExecuteTool(context.Background(), "estimate_cloud_costs", map[string]interface{}{
		"provider":       "AWS",
		"workloadTier":   "Growth",
		"monthlyReqsMil": float64(5.0),
		"storageGB":      float64(100),
	})
	if err != nil || res == nil || dur < 0 {
		t.Errorf("ExecuteTool estimate_cloud_costs failed: %v", err)
	}

	// Test executing unknown tool
	_, _, errUnknown := toolSvc.ExecuteTool(context.Background(), "non_existent_tool", nil)
	if errUnknown == nil {
		t.Error("Expected error for non-existent tool, got nil")
	}
}

func TestStorage_AuditReportPersistence(t *testing.T) {
	wsRepo, err := storage.NewWorkspaceRepository("")
	if err != nil {
		t.Fatalf("Failed creating workspace repo: %v", err)
	}

	report := entity.CodeAuditReport{
		RepositoryName:     "test-repo",
		RiskScore:          20,
		RiskLevel:          "LOW",
		AnalyzedFilesCount: 10,
	}

	err = wsRepo.SaveAuditReport(context.Background(), "session-audit-persist", report)
	if err != nil {
		t.Fatalf("Failed to save audit report: %v", err)
	}

	loaded, err := wsRepo.GetAuditReport(context.Background(), "session-audit-persist")
	if err != nil {
		t.Fatalf("Failed to get audit report: %v", err)
	}
	if loaded.RepositoryName != "test-repo" || loaded.RiskScore != 20 {
		t.Errorf("Report mismatch: got %+v", loaded)
	}
}

func TestAuditHandler(t *testing.T) {
	mockAI := &mockAIService{
		reply: `{"detectedStack":["Go","Docker"],"riskScore":15,"riskLevel":"LOW","architectureSummary":"Clean Go Arch","vulnerabilities":[],"recommendations":["Add linting"]}`,
	}
	wsRepo, _ := storage.NewWorkspaceRepository("")
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	auditUC := usecase.NewCodeAuditUseCase(mockAI, wsRepo, sessionRepo, profileRepo)

	h := handler.NewAuditHandler(auditUC)

	// Test empty / invalid URL
	reqBad := httptest.NewRequest(http.MethodPost, "/api/audit/github", nil)
	rrBad := httptest.NewRecorder()
	h.ServeHTTP(rrBad, reqBad)

	if rrBad.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty payload, got %d", rrBad.Code)
	}
}

func TestChatUseCase_StreamWithTools(t *testing.T) {
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	defer closeCache()

	toolService := tools.NewToolService()
	mockAI := &mockAIService{
		toolName: "estimate_cloud_costs",
		toolArgs: map[string]interface{}{
			"provider":       "AWS",
			"workloadTier":   "Growth",
			"monthlyReqsMil": float64(10.0),
		},
		streamToken: "Aqui está a estimativa de custos:",
	}
	chatUC := usecase.NewChatUseCase(mockAI, toolService, sessionRepo, cacheRepo, profileRepo)

	var emittedEvents []entity.StreamEvent
	emitter := &testEmitter{
		onEmit: func(e entity.StreamEvent) {
			emittedEvents = append(emittedEvents, e)
		},
	}

	err := chatUC.ExecuteStreamChat(context.Background(), usecase.ChatInput{
		Message:   "Quanto custa 10M de requisições na AWS?",
		SessionID: "session-tool-test",
	}, emitter)

	if err != nil {
		t.Fatalf("Unexpected error in ExecuteStreamChat: %v", err)
	}
	if len(emittedEvents) == 0 {
		t.Error("Expected stream events to be emitted")
	}
}

type testEmitter struct {
	onEmit func(e entity.StreamEvent)
}

func (e *testEmitter) Emit(event entity.StreamEvent) {
	if e.onEmit != nil {
		e.onEmit(event)
	}
}
