package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// SessionHandler handles session listing, retrieval, and deletion.
type SessionHandler struct {
	sessionUseCase *usecase.SessionUseCase
}

// NewSessionHandler creates a new SessionHandler.
func NewSessionHandler(sessionUseCase *usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{sessionUseCase: sessionUseCase}
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions")
	path = strings.TrimPrefix(path, "/")

	if r.Method == http.MethodDelete {
		if path == "" {
			http.Error(w, "Session ID is required for deletion", http.StatusBadRequest)
			return
		}
		if err := h.sessionUseCase.DeleteSession(r.Context(), path); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"deletedId": path,
		})
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if path != "" {
		messages, err := h.sessionUseCase.GetSessionHistory(r.Context(), path)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Session not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sessionId": path,
			"messages":  messages,
		})
		return
	}

	sessions, err := h.sessionUseCase.ListSessions(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    len(sessions),
		"sessions": sessions,
	})
}
