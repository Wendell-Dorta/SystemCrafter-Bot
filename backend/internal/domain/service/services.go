package service

import (
	"context"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// AIService defines the port for Generative AI model interactions (e.g. Gemini).
type AIService interface {
	Generate(ctx context.Context, history []entity.Message, toolDecls interface{}) (reply string, toolCallName string, toolCallArgs map[string]interface{}, err error)
	StreamGenerate(ctx context.Context, history []entity.Message, toolDecls interface{}, chunkHandler func(token string, fnName string, fnArgs map[string]interface{}) error) error
}

// ToolService defines the port for registering and executing architectural tools.
type ToolService interface {
	GetToolDeclarations() interface{}
	ExecuteTool(ctx context.Context, name string, args map[string]interface{}) (result interface{}, durationMs int64, err error)
}
