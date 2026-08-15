package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
	tooladapter "github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
	transporthttp "github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

type mockAIServiceForRouter struct{}

func (m *mockAIServiceForRouter) Generate(ctx context.Context, history []entity.Message, toolDecls interface{}) (string, string, map[string]interface{}, error) {
	return "Resposta simulada de arquitetura.", "", nil, nil
}

func (m *mockAIServiceForRouter) StreamGenerate(ctx context.Context, history []entity.Message, toolDecls interface{}, chunkHandler func(token string, fnName string, fnArgs map[string]interface{}) error) error {
	return chunkHandler("Token1", "", nil)
}

func setupTestRouter() (http.Handler, func()) {
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	templateRepo := memory.NewTemplateRepository()
	toolService := tooladapter.NewToolService()
	aiService := &mockAIServiceForRouter{}

	workspaceRepo, _ := storage.NewWorkspaceRepository("./tmp_integration_ws")
	auditUseCase := usecase.NewCodeAuditUseCase(aiService, workspaceRepo, sessionRepo, profileRepo)

	chatUseCase := usecase.NewChatUseCase(aiService, toolService, sessionRepo, cacheRepo, profileRepo)
	templateUseCase := usecase.NewTemplateUseCase(templateRepo)
	sessionUseCase := usecase.NewSessionUseCase(sessionRepo)
	profileUseCase := usecase.NewProfileUseCase(profileRepo)
	healthUseCase := usecase.NewHealthUseCase(cacheRepo)

	rateLimiter := middleware.NewRateLimiter(50.0, 100)
	securityOpts := middleware.SecurityOptions{
		AllowedOrigins: []string{"*"},
		MaxBodyBytes:   1024 * 1024,
	}

	router := transporthttp.NewRouter(transporthttp.RouterConfig{
		ChatUseCase:     chatUseCase,
		TemplateUseCase: templateUseCase,
		SessionUseCase:  sessionUseCase,
		ProfileUseCase:  profileUseCase,
		AuditUseCase:    auditUseCase,
		HealthUseCase:   healthUseCase,
		RateLimiter:     rateLimiter,
		SecurityOpts:    securityOpts,
	})

	return router, func() {
		closeCache()
		_ = os.RemoveAll("./tmp_integration_ws")
	}
}

func TestRouter_HealthEndpoint(t *testing.T) {
	router, cleanup := setupTestRouter()
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /health, got %d", rec.Code)
	}
}

func TestRouter_TemplatesEndpoint(t *testing.T) {
	router, cleanup := setupTestRouter()
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /api/templates, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode json: %v", err)
	}
	if resp["total"].(float64) == 0 {
		t.Errorf("Expected templates in response")
	}
}

func TestRouter_ProfileEndpoint(t *testing.T) {
	router, cleanup := setupTestRouter()
	defer cleanup()

	// GET Profile
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /api/profile, got %d", rec.Code)
	}

	var p entity.ArchitectProfile
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("Failed to decode profile json: %v", err)
	}
	if p.ID == "" {
		t.Errorf("Expected valid profile ID")
	}
}

func TestRouter_ChatEndpoint(t *testing.T) {
	router, cleanup := setupTestRouter()
	defer cleanup()

	body := bytes.NewReader([]byte(`{"message": "Como desenhar um sistema com CQRS?"}`))
	req := httptest.NewRequest(http.MethodPost, "/api/chat", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /api/chat, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestRouter_AuditUploadEndpoint(t *testing.T) {
	router, cleanup := setupTestRouter()
	defer cleanup()

	payload := `{
		"projectName": "TestApp",
		"sessionId": "test-sess",
		"files": [
			{"path": "main.go", "content": "package main\nfunc main() {}"},
			{"path": ".env", "content": "SECRET=123"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/audit/upload", bytes.NewReader([]byte(payload)))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /api/audit/upload, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	var report entity.CodeAuditReport
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse report: %v", err)
	}

	if report.AnalyzedFilesCount != 1 {
		t.Errorf("Expected 1 analyzed file, got %d", report.AnalyzedFilesCount)
	}
	if report.FilteredFilesCount != 1 {
		t.Errorf("Expected 1 filtered file (.env), got %d", report.FilteredFilesCount)
	}
}
