package repository

import (
	"context"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// SessionRepository defines the port for managing conversation sessions and message history.
type SessionRepository interface {
	GetOrCreateSession(ctx context.Context, sessionID string) (*entity.Session, error)
	AppendMessage(ctx context.Context, sessionID string, msg entity.Message) error
	GetSessionMessages(ctx context.Context, sessionID string) ([]entity.Message, error)
	GetAllSessions(ctx context.Context) ([]*entity.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

// ProfileRepository defines the port for storing and retrieving learned engineering profile memory.
type ProfileRepository interface {
	GetProfile(ctx context.Context) (*entity.ArchitectProfile, error)
	UpdateProfile(ctx context.Context, profile entity.ArchitectProfile) error
	AddNote(ctx context.Context, note string) error
	ResetProfile(ctx context.Context) error
}

// TemplateRepository defines the port for querying architecture blueprints and reference designs.
type TemplateRepository interface {
	GetTemplates(ctx context.Context) ([]entity.ArchitectureTemplate, error)
	GetTemplateByID(ctx context.Context, id string) (*entity.ArchitectureTemplate, error)
}

// CacheRepository defines the port for in-memory caching and performance acceleration.
type CacheRepository interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration)
	Delete(ctx context.Context, key string)
	Stats() map[string]interface{}
}

// WorkspaceRepository defines the port for storing and retrieving uploaded code workspaces on disk.
type WorkspaceRepository interface {
	SaveWorkspace(ctx context.Context, id string, files []entity.CodeFile) (string, error)
	LoadWorkspace(ctx context.Context, id string) ([]entity.CodeFile, error)
	SaveAuditReport(ctx context.Context, id string, report entity.CodeAuditReport) error
	GetAuditReport(ctx context.Context, id string) (*entity.CodeAuditReport, error)
}
