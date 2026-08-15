package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

type mockAIService struct {
	reply       string
	toolName    string
	toolArgs    map[string]interface{}
	streamToken string
}

func (m *mockAIService) Generate(ctx context.Context, history []entity.Message, toolDecls interface{}) (string, string, map[string]interface{}, error) {
	return m.reply, m.toolName, m.toolArgs, nil
}

func (m *mockAIService) StreamGenerate(ctx context.Context, history []entity.Message, toolDecls interface{}, chunkHandler func(token string, fnName string, fnArgs map[string]interface{}) error) error {
	if m.toolName != "" {
		if err := chunkHandler("", m.toolName, m.toolArgs); err != nil {
			return err
		}
	}
	if m.streamToken != "" {
		if err := chunkHandler(m.streamToken, "", nil); err != nil {
			return err
		}
	}
	return nil
}

func TestTemplateUseCase(t *testing.T) {
	repo := memory.NewTemplateRepository()
	uc := usecase.NewTemplateUseCase(repo)

	// List
	templates, err := uc.ListTemplates(context.Background())
	if err != nil || len(templates) == 0 {
		t.Fatalf("Expected templates list, got err: %v", err)
	}

	// Get by ID
	tmpl, err := uc.GetTemplateByID(context.Background(), "ecommerce-event-driven")
	if err != nil || tmpl.ID != "ecommerce-event-driven" {
		t.Fatalf("Expected template 'ecommerce-event-driven', got err: %v", err)
	}

	// Get non-existent
	_, err = uc.GetTemplateByID(context.Background(), "non-existent-id")
	if err == nil {
		t.Fatalf("Expected error for non-existent template ID")
	}
}

func TestSessionUseCase(t *testing.T) {
	repo := memory.NewSessionRepository()
	uc := usecase.NewSessionUseCase(repo)

	_ = repo.AppendMessage(context.Background(), "sess-abc", entity.Message{
		Role:    entity.RoleUser,
		Content: "Pergunta sobre arquitetura",
	})

	sessions, err := uc.ListSessions(context.Background())
	if err != nil || len(sessions) != 1 {
		t.Fatalf("Expected 1 session, got %d (err: %v)", len(sessions), err)
	}

	history, err := uc.GetSessionHistory(context.Background(), "sess-abc")
	if err != nil || len(history) != 1 {
		t.Fatalf("Expected 1 message in history, got %d (err: %v)", len(history), err)
	}

	// Test deletion
	_ = uc.DeleteSession(context.Background(), "sess-abc")
	sessionsAfter, _ := uc.ListSessions(context.Background())
	if len(sessionsAfter) != 0 {
		t.Fatalf("Expected 0 sessions after deletion, got %d", len(sessionsAfter))
	}
}

func TestProfileUseCase(t *testing.T) {
	repo := memory.NewProfileRepository()
	uc := usecase.NewProfileUseCase(repo)

	profile, err := uc.GetProfile(context.Background())
	if err != nil || profile == nil {
		t.Fatalf("Expected default profile, got err: %v", err)
	}

	// Test automatic learning
	uc.LearnFromMessage(context.Background(), "Na nossa empresa usamos AWS e preferimos PostgreSQL com Go e conformidade LGPD.")

	updated, _ := uc.GetProfile(context.Background())
	if updated.PreferredCloud != "AWS" {
		t.Errorf("Expected PreferredCloud 'AWS', got %s", updated.PreferredCloud)
	}

	// Test note addition
	_ = uc.AddNote(context.Background(), "Preferência por arquitetura orientada a microsserviços.")
	pWithNote, _ := uc.GetProfile(context.Background())
	if len(pWithNote.CustomNotes) == 0 {
		t.Errorf("Expected custom notes to be recorded")
	}

	// Test reset
	_ = uc.ResetProfile(context.Background())
	pReset, _ := uc.GetProfile(context.Background())
	if pReset.PreferredCloud != "Multi-Cloud" {
		t.Errorf("Expected PreferredCloud 'Multi-Cloud' after reset, got %s", pReset.PreferredCloud)
	}
	if len(pReset.PrimaryLanguages) != 2 {
		t.Errorf("Expected 2 default primary languages after reset, got %d", len(pReset.PrimaryLanguages))
	}
}

func TestHealthUseCase(t *testing.T) {
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	defer closeCache()

	uc := usecase.NewHealthUseCase(cacheRepo)
	health := uc.GetHealth(context.Background(), "1h 30m", 150, map[string]interface{}{"rps": 15})

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", health.Status)
	}
	if health.Uptime != "1h 30m" {
		t.Errorf("Expected uptime '1h 30m', got %s", health.Uptime)
	}
}

func TestChatUseCase_ExecuteChat_WithMockAI(t *testing.T) {
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	defer closeCache()
	toolService := tools.NewToolService()

	mockAI := &mockAIService{
		reply: "Para alta disponibilidade, utilize réplicas de leitura e Multi-AZ.",
	}

	chatUC := usecase.NewChatUseCase(mockAI, toolService, sessionRepo, cacheRepo, profileRepo)

	out, err := chatUC.ExecuteChat(context.Background(), usecase.ChatInput{
		SessionID: "sess-test-1",
		Message:   "Como garantir alta disponibilidade no PostgreSQL?",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if out.Reply != mockAI.reply {
		t.Errorf("Expected reply '%s', got '%s'", mockAI.reply, out.Reply)
	}

	_, errEmpty := chatUC.ExecuteChat(context.Background(), usecase.ChatInput{
		Message: "",
	})
	if errEmpty == nil {
		t.Fatalf("Expected validation error on empty message")
	}
}
