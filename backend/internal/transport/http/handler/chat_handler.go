package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// ChatHandler handles standard synchronous chat requests.
type ChatHandler struct {
	chatUseCase *usecase.ChatUseCase
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(chatUseCase *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{chatUseCase: chatUseCase}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID     string           `json:"sessionId"`
		Message       string           `json:"message"`
		ImageBase64   string           `json:"imageBase64,omitempty"`
		ImageMimeType string           `json:"imageMimeType,omitempty"`
		History       []entity.Message `json:"history,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if strings.TrimSpace(req.Message) == "" && strings.TrimSpace(req.ImageBase64) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Message or image is required"})
		return
	}

	input := usecase.ChatInput{
		SessionID:     req.SessionID,
		Message:       req.Message,
		ImageBase64:   req.ImageBase64,
		ImageMimeType: req.ImageMimeType,
		History:       req.History,
	}

	output, err := h.chatUseCase.ExecuteChat(r.Context(), input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
