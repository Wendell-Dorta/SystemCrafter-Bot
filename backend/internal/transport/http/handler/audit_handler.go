package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// AuditHandler handles code and GitHub repository auditing endpoints.
type AuditHandler struct {
	auditUseCase *usecase.CodeAuditUseCase
}

// NewAuditHandler creates an instance of AuditHandler.
func NewAuditHandler(uc *usecase.CodeAuditUseCase) *AuditHandler {
	return &AuditHandler{auditUseCase: uc}
}

func (h *AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/audit")
	path = strings.TrimPrefix(path, "/")

	switch {
	case r.Method == http.MethodPost && path == "github":
		var req struct {
			GitHubURL string `json:"githubUrl"`
			SessionID string `json:"sessionId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GitHubURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "URL do repositório GitHub é obrigatória (ex: https://github.com/owner/repo)"})
			return
		}

		report, err := h.auditUseCase.AuditGitHubRepository(r.Context(), req.GitHubURL, req.SessionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(report)

	case r.Method == http.MethodPost && path == "upload":
		var req struct {
			ProjectName string            `json:"projectName"`
			SessionID   string            `json:"sessionId"`
			Files       []entity.CodeFile `json:"files"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Files) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Lista de arquivos de código é obrigatória"})
			return
		}

		report, err := h.auditUseCase.AuditUploadedFiles(r.Context(), req.Files, req.SessionID, req.ProjectName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(report)

	case r.Method == http.MethodGet && strings.HasPrefix(path, "workspaces/"):
		workspaceID := strings.TrimPrefix(path, "workspaces/")
		report, err := h.auditUseCase.GetReport(r.Context(), workspaceID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(report)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}
