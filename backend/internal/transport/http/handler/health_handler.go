package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// HealthHandler serves system health, rate limiter stats and cache metrics.
type HealthHandler struct {
	startTime     time.Time
	healthUseCase *usecase.HealthUseCase
	rateLimiter   *middleware.RateLimiter
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(healthUseCase *usecase.HealthUseCase, rateLimiter *middleware.RateLimiter) *HealthHandler {
	return &HealthHandler{
		startTime:     time.Now(),
		healthUseCase: healthUseCase,
		rateLimiter:   rateLimiter,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Round(time.Second)
	rateStats := h.rateLimiter.Stats()

	totalReqs := int64(0)
	if val, ok := rateStats["totalRequests"].(int64); ok {
		totalReqs = val
	}

	health := h.healthUseCase.GetHealth(r.Context(), uptime.String(), totalReqs, rateStats)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"system":    "ArchMind AI Core Engine (Clean Architecture)",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    fmt.Sprintf("%s (%s)", uptime, h.startTime.Format("2006-01-02 15:04:05 MST")),
		"status":    health,
	})
}
