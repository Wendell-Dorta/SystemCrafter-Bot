package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// TemplateHandler handles requests for architectural blueprints.
type TemplateHandler struct {
	templateUseCase *usecase.TemplateUseCase
}

// NewTemplateHandler creates a new TemplateHandler.
func NewTemplateHandler(templateUseCase *usecase.TemplateUseCase) *TemplateHandler {
	return &TemplateHandler{templateUseCase: templateUseCase}
}

func (h *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/templates")
	path = strings.TrimPrefix(path, "/")

	if path != "" {
		tmpl, err := h.templateUseCase.GetTemplateByID(r.Context(), path)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tmpl)
		return
	}

	templates, err := h.templateUseCase.ListTemplates(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":     len(templates),
		"templates": templates,
	})
}
