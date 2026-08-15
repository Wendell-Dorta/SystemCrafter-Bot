package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/service"
)

type geminiAIService struct {
	client *GeminiClient
}

// NewGeminiAIService creates an AIService adapter for Google Gemini.
func NewGeminiAIService(apiKey, model string) service.AIService {
	return &geminiAIService{
		client: NewGeminiClient(apiKey, model),
	}
}

func (s *geminiAIService) Generate(ctx context.Context, history []entity.Message, toolDecls interface{}) (string, string, map[string]interface{}, error) {
	contents := mapDomainMessagesToGemini(history)
	var decls []FunctionDeclaration
	if td, ok := toolDecls.([]FunctionDeclaration); ok {
		decls = td
	}

	resp, err := s.client.GenerateContent(ctx, contents, decls)
	if err != nil {
		return "", "", nil, err
	}

	if len(resp.Candidates) == 0 {
		return "", "", nil, nil
	}

	var replyText strings.Builder
	var fnName string
	var fnArgs map[string]interface{}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			fnName = part.FunctionCall.Name
			fnArgs = part.FunctionCall.Args
		}
		if part.Text != "" {
			replyText.WriteString(part.Text)
		}
	}

	return replyText.String(), fnName, fnArgs, nil
}

func (s *geminiAIService) StreamGenerate(
	ctx context.Context,
	history []entity.Message,
	toolDecls interface{},
	chunkHandler func(token string, fnName string, fnArgs map[string]interface{}) error,
) error {
	contents := mapDomainMessagesToGemini(history)
	var decls []FunctionDeclaration
	if td, ok := toolDecls.([]FunctionDeclaration); ok {
		decls = td
	}

	return s.client.StreamGenerateContent(ctx, contents, decls, func(chunk *GenerateContentResponse) error {
		if len(chunk.Candidates) == 0 {
			return nil
		}
		for _, part := range chunk.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				if err := chunkHandler("", part.FunctionCall.Name, part.FunctionCall.Args); err != nil {
					return err
				}
			}
			if part.Text != "" {
				if err := chunkHandler(part.Text, "", nil); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// mapDomainMessagesToGemini applies sliding window, optimizes token usage and prevents re-uploading giant images.
func mapDomainMessagesToGemini(messages []entity.Message) []Content {
	// Sliding window: keep only the most recent 6 messages to prevent token explosion
	maxHistory := 6
	startIdx := 0
	if len(messages) > maxHistory {
		startIdx = len(messages) - maxHistory
	}
	recentMessages := messages[startIdx:]

	var contents []Content
	lastIdx := len(recentMessages) - 1

	for i, m := range recentMessages {
		var parts []Part

		// Text content
		if m.Content != "" {
			parts = append(parts, Part{Text: m.Content})
		}

		// Tool execution synthesis context
		if len(m.ToolEvents) > 0 {
			for _, te := range m.ToolEvents {
				resJSON, _ := json.Marshal(te.Result)
				toolSummary := fmt.Sprintf("\n[Resultado da Ferramenta %s]:\n```json\n%s\n```\n", te.ToolName, string(resJSON))
				parts = append(parts, Part{Text: toolSummary})
			}
		}

		// Image payload optimization: ONLY send raw base64 on the current/latest user message!
		// Prevents re-sending megabytes of images on every subsequent turn.
		if m.ImageBase64 != "" {
			if i == lastIdx {
				mime := m.ImageMimeType
				if mime == "" {
					mime = "image/png"
				}
				parts = append(parts, Part{
					InlineData: &Blob{
						MimeType: mime,
						Data:     m.ImageBase64,
					},
				})
			} else {
				// Historical placeholder
				parts = append(parts, Part{Text: "[Imagem de arquitetura enviada pelo usuário no início da sessão]"})
			}
		}

		role := "user"
		if m.Role == entity.RoleModel {
			role = "model"
		}
		if len(parts) > 0 {
			contents = append(contents, Content{
				Role:  role,
				Parts: parts,
			})
		}
	}
	return contents
}
