package unit_test

import (
	"context"
	"testing"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/ai"
)

func TestGeminiService_InitializationAndCancellation(t *testing.T) {
	service := ai.NewGeminiAIService("test-key", "gemini-2.5-flash")
	if service == nil {
		t.Fatal("Expected non-nil GeminiAIService")
	}

	// Test with immediate context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	history := []entity.Message{
		{Role: "user", Content: "Hello architect"},
	}

	_, _, _, err := service.Generate(ctx, history, nil)
	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}

	errStream := service.StreamGenerate(ctx, history, nil, func(token string, fnName string, fnArgs map[string]interface{}) error {
		return nil
	})
	if errStream == nil {
		t.Error("Expected error with cancelled context in stream, got nil")
	}
}
