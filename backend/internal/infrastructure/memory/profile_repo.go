package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

type memoryProfileRepo struct {
	mu      sync.RWMutex
	profile entity.ArchitectProfile
}

// NewProfileRepository creates a thread-safe in-memory profile memory repository.
func NewProfileRepository() repository.ProfileRepository {
	return &memoryProfileRepo{
		profile: entity.ArchitectProfile{
			ID:                 "default-architect",
			PrimaryLanguages:   []string{"Go", "TypeScript"},
			PreferredCloud:     "AWS",
			PreferredDatabases: []string{"PostgreSQL", "Redis"},
			PreferredPatterns:  []string{"Event-Driven", "CQRS", "Transactional Outbox"},
			ComplianceRules:    []string{"LGPD", "OWASP Top 10"},
			CustomNotes: []string{
				"Prioriza simplicidade operacional, alta coesão e baixo acoplamento.",
				"Prefere soluções cloud native gerenciadas e contêineres Docker.",
			},
			UpdatedAt: time.Now(),
		},
	}
}

func (r *memoryProfileRepo) GetProfile(ctx context.Context) (*entity.ArchitectProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pCopy := r.profile
	return &pCopy, nil
}

func (r *memoryProfileRepo) UpdateProfile(ctx context.Context, profile entity.ArchitectProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	profile.ID = "default-architect"
	profile.UpdatedAt = time.Now()
	r.profile = profile
	return nil
}

func (r *memoryProfileRepo) AddNote(ctx context.Context, note string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if note == "" {
		return nil
	}
	r.profile.CustomNotes = append(r.profile.CustomNotes, note)
	r.profile.UpdatedAt = time.Now()
	return nil
}

func (r *memoryProfileRepo) ResetProfile(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.profile = entity.ArchitectProfile{
		ID:                 "default-architect",
		PrimaryLanguages:   []string{"Go", "TypeScript"},
		PreferredCloud:     "Multi-Cloud",
		PreferredDatabases: []string{"PostgreSQL", "Redis"},
		PreferredPatterns:  []string{"Event-Driven", "CQRS", "Transactional Outbox"},
		ComplianceRules:    []string{"LGPD", "OWASP Top 10"},
		CustomNotes: []string{
			"Prioriza simplicidade operacional, alta coesão e baixo acoplamento.",
			"Prefere soluções cloud native gerenciadas e contêineres Docker.",
		},
		UpdatedAt: time.Now(),
	}
	return nil
}
