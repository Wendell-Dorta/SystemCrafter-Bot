package usecase

import (
	"context"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

// TemplateUseCase handles architectural blueprints business logic.
type TemplateUseCase struct {
	repo repository.TemplateRepository
}

// NewTemplateUseCase creates a new TemplateUseCase.
func NewTemplateUseCase(repo repository.TemplateRepository) *TemplateUseCase {
	return &TemplateUseCase{repo: repo}
}

// ListTemplates retrieves all available preloaded architecture blueprints.
func (uc *TemplateUseCase) ListTemplates(ctx context.Context) ([]entity.ArchitectureTemplate, error) {
	return uc.repo.GetTemplates(ctx)
}

// GetTemplateByID retrieves a specific blueprint by its identifier.
func (uc *TemplateUseCase) GetTemplateByID(ctx context.Context, id string) (*entity.ArchitectureTemplate, error) {
	return uc.repo.GetTemplateByID(ctx, id)
}
