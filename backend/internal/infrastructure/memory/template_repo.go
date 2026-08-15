package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

type memoryTemplateRepo struct {
	mu        sync.RWMutex
	templates map[string]entity.ArchitectureTemplate
}

// NewTemplateRepository creates an initialized TemplateRepository with preloaded seeds.
func NewTemplateRepository() repository.TemplateRepository {
	repo := &memoryTemplateRepo{
		templates: make(map[string]entity.ArchitectureTemplate),
	}

	seeds := GetPreloadedTemplates()
	for _, s := range seeds {
		repo.templates[s.ID] = s
	}

	return repo
}

func (r *memoryTemplateRepo) GetTemplates(ctx context.Context) ([]entity.ArchitectureTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.ArchitectureTemplate, 0, len(r.templates))
	for _, t := range r.templates {
		result = append(result, t)
	}
	return result, nil
}

func (r *memoryTemplateRepo) GetTemplateByID(ctx context.Context, id string) (*entity.ArchitectureTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.templates[id]
	if !exists {
		return nil, fmt.Errorf("template with id '%s' not found", id)
	}
	return &t, nil
}
