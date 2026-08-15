package usecase

import (
	"context"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

// SessionUseCase handles conversation sessions and history retrieval.
type SessionUseCase struct {
	repo repository.SessionRepository
}

// NewSessionUseCase creates a new SessionUseCase.
func NewSessionUseCase(repo repository.SessionRepository) *SessionUseCase {
	return &SessionUseCase{repo: repo}
}

// ListSessions retrieves all stored sessions.
func (uc *SessionUseCase) ListSessions(ctx context.Context) ([]*entity.Session, error) {
	return uc.repo.GetAllSessions(ctx)
}

// GetSessionHistory retrieves messages for a specific session.
func (uc *SessionUseCase) GetSessionHistory(ctx context.Context, sessionID string) ([]entity.Message, error) {
	return uc.repo.GetSessionMessages(ctx, sessionID)
}

// DeleteSession removes a session from history.
func (uc *SessionUseCase) DeleteSession(ctx context.Context, sessionID string) error {
	return uc.repo.DeleteSession(ctx, sessionID)
}
