package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// ChatStreamHandler handles HTTP SSE streaming for agentic conversations.
type ChatStreamHandler struct {
	chatUseCase *usecase.ChatUseCase
}

// NewChatStreamHandler creates a new ChatStreamHandler.
func NewChatStreamHandler(chatUseCase *usecase.ChatUseCase) *ChatStreamHandler {
	return &ChatStreamHandler{chatUseCase: chatUseCase}
}

type sseEmitter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (e *sseEmitter) Emit(event entity.StreamEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	fmt.Fprintf(e.w, "event: message\ndata: %s\n\n", string(data))
	e.flusher.Flush()
}

func (h *ChatStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by server", http.StatusInternalServerError)
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.Message) == "" && strings.TrimSpace(req.ImageBase64) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Message or image is required"})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	emitter := &sseEmitter{w: w, flusher: flusher}

	input := usecase.ChatInput{
		SessionID:     req.SessionID,
		Message:       req.Message,
		ImageBase64:   req.ImageBase64,
		ImageMimeType: req.ImageMimeType,
		History:       req.History,
	}

	_ = h.chatUseCase.ExecuteStreamChat(r.Context(), input, emitter)
}
