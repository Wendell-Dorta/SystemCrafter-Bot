package usecase

import (
	"context"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

// HealthUseCase generates real-time telemetry and health data.
type HealthUseCase struct {
	cacheRepo repository.CacheRepository
}

// NewHealthUseCase creates a new HealthUseCase.
func NewHealthUseCase(cacheRepo repository.CacheRepository) *HealthUseCase {
	return &HealthUseCase{cacheRepo: cacheRepo}
}

// GetHealth collects telemetry metrics.
func (uc *HealthUseCase) GetHealth(ctx context.Context, uptime string, totalRequests int64, rateLimitStats map[string]interface{}) entity.SystemHealth {
	var cacheStats map[string]interface{}
	if uc.cacheRepo != nil {
		cacheStats = uc.cacheRepo.Stats()
	}

	return entity.SystemHealth{
		Status:         "healthy",
		Version:        "2.0.0-clean-arch",
		Uptime:         uptime,
		TotalRequests:  totalRequests,
		CacheStats:     cacheStats,
		RateLimitStats: rateLimitStats,
	}
}
