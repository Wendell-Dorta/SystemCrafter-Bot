package unit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/handler"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

func TestHealthHandler(t *testing.T) {
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	defer closeCache()

	healthUC := usecase.NewHealthUseCase(cacheRepo)
	rl := middleware.NewRateLimiter(60, 100)
	h := handler.NewHealthHandler(healthUC, rl)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestTemplateHandler(t *testing.T) {
	templateRepo := memory.NewTemplateRepository()
	templateUC := usecase.NewTemplateUseCase(templateRepo)
	h := handler.NewTemplateHandler(templateUC)

	// Test GET all
	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Test GET by ID
	reqID := httptest.NewRequest(http.MethodGet, "/api/templates/ecommerce-event-driven", nil)
	rrID := httptest.NewRecorder()
	h.ServeHTTP(rrID, reqID)

	if rrID.Code != http.StatusOK {
		t.Errorf("Expected status 200 for template ID, got %d", rrID.Code)
	}
}

func TestSessionHandler(t *testing.T) {
	sessionRepo := memory.NewSessionRepository()
	sessionUC := usecase.NewSessionUseCase(sessionRepo)
	h := handler.NewSessionHandler(sessionUC)

	// 1. Create a session
	_, _ = sessionRepo.GetOrCreateSession(context.Background(), "session-test-123")

	// 2. GET all sessions
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// 3. GET session messages
	reqMsg := httptest.NewRequest(http.MethodGet, "/api/sessions/session-test-123/messages", nil)
	rrMsg := httptest.NewRecorder()
	h.ServeHTTP(rrMsg, reqMsg)

	if rrMsg.Code != http.StatusOK {
		t.Errorf("Expected status 200 for messages, got %d", rrMsg.Code)
	}

	// 4. DELETE session
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/sessions/session-test-123", nil)
	rrDel := httptest.NewRecorder()
	h.ServeHTTP(rrDel, reqDel)

	if rrDel.Code != http.StatusOK {
		t.Errorf("Expected status 200 for delete, got %d", rrDel.Code)
	}
}

func TestProfileHandler(t *testing.T) {
	profileRepo := memory.NewProfileRepository()
	profileUC := usecase.NewProfileUseCase(profileRepo)
	h := handler.NewProfileHandler(profileUC)

	// 1. GET Profile
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// 2. PUT Profile
	updatedProfile := entity.ArchitectProfile{
		PreferredCloud:     "Azure",
		PrimaryLanguages:   []string{"C#", "TypeScript"},
		PreferredDatabases: []string{"CosmosDB", "Redis"},
	}
	body, _ := json.Marshal(updatedProfile)
	reqPut := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	rrPut := httptest.NewRecorder()
	h.ServeHTTP(rrPut, reqPut)

	if rrPut.Code != http.StatusOK {
		t.Errorf("Expected status 200 for PUT, got %d", rrPut.Code)
	}

	// 3. POST Note (Route: /api/profile/note)
	noteBody, _ := json.Marshal(map[string]string{"note": "Preferência por containers Linux"})
	reqNote := httptest.NewRequest(http.MethodPost, "/api/profile/note", bytes.NewReader(noteBody))
	rrNote := httptest.NewRecorder()
	h.ServeHTTP(rrNote, reqNote)

	if rrNote.Code != http.StatusOK {
		t.Errorf("Expected status 200 for note, got %d", rrNote.Code)
	}

	// 4. POST Reset (Route: /api/profile/reset)
	reqReset := httptest.NewRequest(http.MethodPost, "/api/profile/reset", nil)
	rrReset := httptest.NewRecorder()
	h.ServeHTTP(rrReset, reqReset)

	if rrReset.Code != http.StatusOK {
		t.Errorf("Expected status 200 for reset, got %d", rrReset.Code)
	}
}

func TestChatHandler_And_StreamHandler(t *testing.T) {
	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	cacheRepo, closeCache := memory.NewCacheRepository(10*time.Minute, 1*time.Minute)
	defer closeCache()

	toolService := tools.NewToolService()
	mockAI := &mockAIService{
		reply:       "Arquitetura bancária orientada a eventos e isolamento de falhas",
		streamToken: "Stream token",
	}
	chatUC := usecase.NewChatUseCase(mockAI, toolService, sessionRepo, cacheRepo, profileRepo)

	// 1. Test Sync ChatHandler
	chatH := handler.NewChatHandler(chatUC)
	syncPayload, _ := json.Marshal(map[string]string{
		"message":   "Como desenhar um sistema bancário resiliente?",
		"sessionId": "test-session-sync",
	})
	reqSync := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(syncPayload))
	rrSync := httptest.NewRecorder()
	chatH.ServeHTTP(rrSync, reqSync)

	if rrSync.Code != http.StatusOK {
		t.Errorf("Expected status 200 for sync chat, got %d", rrSync.Code)
	}

	// 2. Test Stream ChatStreamHandler
	streamH := handler.NewChatStreamHandler(chatUC)
	streamPayload, _ := json.Marshal(map[string]string{
		"message":   "Gerar diagrama de microsserviços",
		"sessionId": "test-session-stream",
	})
	reqStream := httptest.NewRequest(http.MethodPost, "/api/chat/stream", bytes.NewReader(streamPayload))
	rrStream := httptest.NewRecorder()
	streamH.ServeHTTP(rrStream, reqStream)

	if rrStream.Code != http.StatusOK {
		t.Errorf("Expected status 200 for stream chat, got %d", rrStream.Code)
	}
	if !bytes.Contains(rrStream.Body.Bytes(), []byte("event: message")) {
		t.Errorf("Expected SSE event in stream body, got %s", rrStream.Body.String())
	}
}
