package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

type memorySessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*entity.Session
}

// NewSessionRepository creates an in-memory thread-safe SessionRepository.
func NewSessionRepository() repository.SessionRepository {
	return &memorySessionRepo{
		sessions: make(map[string]*entity.Session),
	}
}

func (r *memorySessionRepo) GetOrCreateSession(ctx context.Context, sessionID string) (*entity.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	session, exists := r.sessions[sessionID]
	if !exists {
		session = &entity.Session{
			ID:        sessionID,
			Title:     "Nova Sessão de Arquitetura",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Messages:  make([]entity.Message, 0),
		}
		r.sessions[sessionID] = session
	}

	return session, nil
}

func (r *memorySessionRepo) AppendMessage(ctx context.Context, sessionID string, msg entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		session = &entity.Session{
			ID:        sessionID,
			Title:     "Nova Sessão de Arquitetura",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Messages:  make([]entity.Message, 0),
		}
		r.sessions[sessionID] = session
	}

	if msg.ID == "" {
		msg.ID = fmt.Sprintf("msg-%d", time.Now().UnixNano())
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	// Auto generate title on first user message
	if len(session.Messages) == 0 && msg.Role == entity.RoleUser && len(msg.Content) > 0 {
		title := msg.Content
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		session.Title = title
	}

	session.Messages = append(session.Messages, msg)
	session.UpdatedAt = time.Now()
	return nil
}

func (r *memorySessionRepo) GetSessionMessages(ctx context.Context, sessionID string) ([]entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return []entity.Message{}, nil
	}

	msgs := make([]entity.Message, len(session.Messages))
	copy(msgs, session.Messages)
	return msgs, nil
}

func (r *memorySessionRepo) GetAllSessions(ctx context.Context) ([]*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*entity.Session, 0, len(r.sessions))
	for _, sess := range r.sessions {
		list = append(list, &entity.Session{
			ID:        sess.ID,
			Title:     sess.Title,
			CreatedAt: sess.CreatedAt,
			UpdatedAt: sess.UpdatedAt,
			Messages:  sess.Messages,
		})
	}
	return list, nil
}

func (r *memorySessionRepo) DeleteSession(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, sessionID)
	return nil
}
