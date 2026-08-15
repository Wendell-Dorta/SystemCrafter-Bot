package unit_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
)

func TestMemoryCache_Operations(t *testing.T) {
	c := memory.NewMemoryCache(10*time.Minute, 1*time.Minute)
	defer c.Close()

	c.Set("user:123", "Alice", 0)

	val, found := c.Get("user:123")
	if !found || val != "Alice" {
		t.Fatalf("Expected Alice, got %v", val)
	}

	c.Delete("user:123")
	_, found = c.Get("user:123")
	if found {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := memory.NewMemoryCache(50*time.Millisecond, 20*time.Millisecond)
	defer c.Close()

	c.Set("temp", "data", 30*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	_, found := c.Get("temp")
	if found {
		t.Fatalf("Expected key to be expired")
	}
}

func TestSessionRepository_Operations(t *testing.T) {
	repo := memory.NewSessionRepository()
	ctx := context.Background()

	sess, err := repo.GetOrCreateSession(ctx, "s-1")
	if err != nil || sess.ID != "s-1" {
		t.Fatalf("Expected session s-1, got %v", sess)
	}

	_ = repo.AppendMessage(ctx, "s-1", entity.Message{
		Role:    entity.RoleUser,
		Content: "Como desenhar microsserviços?",
	})

	msgs, err := repo.GetSessionMessages(ctx, "s-1")
	if err != nil || len(msgs) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(msgs))
	}
}

func TestSessionRepository_ConcurrentAccess(t *testing.T) {
	repo := memory.NewSessionRepository()
	ctx := context.Background()
	var wg sync.WaitGroup
	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sessID := "concurrent-sess"
			_, _ = repo.GetOrCreateSession(ctx, sessID)
			_ = repo.AppendMessage(ctx, sessID, entity.Message{
				Role:    entity.RoleUser,
				Content: "Ping",
			})
			_, _ = repo.GetSessionMessages(ctx, sessID)
		}(i)
	}

	wg.Wait()

	msgs, _ := repo.GetSessionMessages(ctx, "concurrent-sess")
	if len(msgs) != workers {
		t.Fatalf("Expected %d messages, got %d", workers, len(msgs))
	}
}

func TestTemplateRepository_Operations(t *testing.T) {
	repo := memory.NewTemplateRepository()
	ctx := context.Background()

	templates, err := repo.GetTemplates(ctx)
	if err != nil || len(templates) == 0 {
		t.Fatalf("Expected templates, got %d", len(templates))
	}

	tmpl, err := repo.GetTemplateByID(ctx, "ecommerce-event-driven")
	if err != nil || tmpl.ID != "ecommerce-event-driven" {
		t.Fatalf("Expected template ecommerce-event-driven, got err: %v", err)
	}
}
