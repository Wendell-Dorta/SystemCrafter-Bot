package unit_test

import (
	"context"
	"os"
	"testing"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

func TestSecuritySanitizer_BlocksSensitiveFiles(t *testing.T) {
	sensitivePaths := []string{
		".env",
		".env.local",
		"config/prod.env",
		"id_rsa",
		"certs/server.key",
		"certs/tls.pem",
		"credentials.json",
		"secrets.yaml",
		"node_modules/express/index.js",
		".git/config",
		"bin/server.exe",
	}

	for _, p := range sensitivePaths {
		sensitive, reason := storage.IsSensitivePath(p)
		if !sensitive {
			t.Errorf("Expected path '%s' to be detected as sensitive (got reason: %s)", p, reason)
		}
	}
}

func TestSecuritySanitizer_AllowsValidCodeFiles(t *testing.T) {
	validPaths := []string{
		"cmd/server/main.go",
		"src/components/Header.tsx",
		"api/routes.js",
		"Dockerfile",
		"docker-compose.yml",
		"go.mod",
		"package.json",
	}

	for _, p := range validPaths {
		sensitive, _ := storage.IsSensitivePath(p)
		if sensitive {
			t.Errorf("Expected valid code file '%s' NOT to be marked sensitive", p)
		}
		allowed, _ := storage.IsAllowedCodeFile(p)
		if !allowed {
			t.Errorf("Expected code file '%s' to be allowed", p)
		}
	}
}

func TestSecuritySanitizer_RedactsSecretTokens(t *testing.T) {
	raw := `
apiKey := "AIzaSyD-1234567890abcdefghijklmnopqrstuvw"
dbURL := "postgres://admin:supersecretpassword123@db.internal:5432/orders"
token := "ghp_1234567890abcdefghijklmnopqrstuvwxyz"
`
	sanitized := storage.SanitizeContent(raw)
	isSens, _ := storage.IsSensitivePath("main.go")
	if isSens {
		t.Fail()
	}
	if string(sanitized) == raw {
		t.Errorf("Expected secret tokens to be redacted from content")
	}
}

func TestWorkspaceRepository_SaveAndLoad(t *testing.T) {
	tmpDir := "./tmp_test_workspaces"
	defer os.RemoveAll(tmpDir)

	repo, err := storage.NewWorkspaceRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to init workspace repo: %v", err)
	}

	files := []entity.CodeFile{
		{
			Path:     "cmd/main.go",
			Content:  "package main\nfunc main() {}",
			Language: "Go",
		},
		{
			Path:     "Dockerfile",
			Content:  "FROM alpine",
			Language: "Dockerfile",
		},
	}

	wsID := "test-ws-1"
	_, err = repo.SaveWorkspace(context.Background(), wsID, files)
	if err != nil {
		t.Fatalf("Failed to save workspace: %v", err)
	}

	loaded, err := repo.LoadWorkspace(context.Background(), wsID)
	if err != nil || len(loaded) != 2 {
		t.Fatalf("Expected 2 loaded files, got %d (err: %v)", len(loaded), err)
	}
}

func TestCodeAuditUseCase_AuditUploadedFiles(t *testing.T) {
	tmpDir := "./tmp_audit_usecase_ws"
	defer os.RemoveAll(tmpDir)

	wsRepo, _ := storage.NewWorkspaceRepository(tmpDir)
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	mockAI := &mockAIService{
		reply: `{
			"detectedStack": ["Go 1.24", "Docker"],
			"riskScore": 15,
			"riskLevel": "LOW",
			"architectureSummary": "Aplicação modular em Go",
			"vulnerabilities": [],
			"mermaidDiagram": "graph TD\n  Client --> Server",
			"recommendations": ["Adicionar testes"]
		}`,
	}

	auditUC := usecase.NewCodeAuditUseCase(mockAI, wsRepo, sessionRepo, profileRepo)

	files := []entity.CodeFile{
		{
			Path:    "cmd/main.go",
			Content: "package main\nfunc main() {}",
		},
		{
			Path:    ".env", // Should be discarded
			Content: "SECRET_KEY=12345",
		},
	}

	report, err := auditUC.AuditUploadedFiles(context.Background(), files, "sess-test-audit", "MyService")
	if err != nil {
		t.Fatalf("Failed to audit uploaded files: %v", err)
	}

	if report.AnalyzedFilesCount != 1 {
		t.Errorf("Expected 1 analyzed file (excluding .env), got %d", report.AnalyzedFilesCount)
	}
	if report.FilteredFilesCount != 1 {
		t.Errorf("Expected 1 filtered sensitive file (.env), got %d", report.FilteredFilesCount)
	}
	if report.RiskLevel != "LOW" {
		t.Errorf("Expected risk level 'LOW', got '%s'", report.RiskLevel)
	}
}
