package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/service"
)

// ChatInput represents user input parameters for a conversation turn.
type ChatInput struct {
	SessionID     string
	Message       string
	ImageBase64   string
	ImageMimeType string
	History       []entity.Message
}

// ChatOutput represents synchronous result of a chat turn.
type ChatOutput struct {
	SessionID  string             `json:"sessionId"`
	Reply      string             `json:"reply"`
	ToolEvents []entity.ToolEvent `json:"toolEvents,omitempty"`
	DurationMs int64              `json:"durationMs"`
	Success    bool               `json:"success"`
}

// StreamEmitter is a callback interface for sending streaming events to client.
type StreamEmitter interface {
	Emit(event entity.StreamEvent)
}

// ChatUseCase encapsulates the core conversation and agentic workflow business rules.
type ChatUseCase struct {
	aiService   service.AIService
	toolService service.ToolService
	sessionRepo repository.SessionRepository
	cacheRepo   repository.CacheRepository
	profileRepo repository.ProfileRepository
}

// NewChatUseCase creates a new ChatUseCase instance.
func NewChatUseCase(
	ai service.AIService,
	tools service.ToolService,
	sessions repository.SessionRepository,
	cache repository.CacheRepository,
	profile repository.ProfileRepository,
) *ChatUseCase {
	return &ChatUseCase{
		aiService:   ai,
		toolService: tools,
		sessionRepo: sessions,
		cacheRepo:   cache,
		profileRepo: profile,
	}
}

// ExecuteChat handles synchronous chat queries with caching.
func (uc *ChatUseCase) ExecuteChat(ctx context.Context, input ChatInput) (*ChatOutput, error) {
	if strings.TrimSpace(input.Message) == "" && strings.TrimSpace(input.ImageBase64) == "" {
		return nil, fmt.Errorf("message or image is required")
	}

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	// Learn user preferences from message
	uc.autoLearn(ctx, input.Message)

	// 1. Cache lookup for pure text queries
	cacheKey := ""
	if input.ImageBase64 == "" && len(input.History) == 0 && uc.cacheRepo != nil {
		hSum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(input.Message))))
		cacheKey = "chat:" + hex.EncodeToString(hSum[:])
		if cached, found := uc.cacheRepo.Get(ctx, cacheKey); found {
			if out, ok := cached.(ChatOutput); ok {
				out.SessionID = sessionID
				return &out, nil
			}
		}
	}

	start := time.Now()

	// 2. Persist User Message
	userMsg := entity.Message{
		ID:            fmt.Sprintf("usr-%d", time.Now().UnixNano()),
		Role:          entity.RoleUser,
		Content:       input.Message,
		ImageBase64:   input.ImageBase64,
		ImageMimeType: input.ImageMimeType,
		CreatedAt:     time.Now(),
	}
	if err := uc.sessionRepo.AppendMessage(ctx, sessionID, userMsg); err != nil {
		return nil, fmt.Errorf("failed to persist user message: %w", err)
	}

	history, err := uc.sessionRepo.GetSessionMessages(ctx, sessionID)
	if err != nil || len(history) == 0 {
		history = input.History
	}

	// Sliding window: keep only the most recent turns to maintain low latency and low token consumption
	const maxRecentMessages = 8
	if len(history) > maxRecentMessages {
		history = history[len(history)-maxRecentMessages:]
	}

	// Inject Long-Term Memory Context
	historyWithMemory := uc.injectMemoryContext(ctx, history)

	toolDecls := uc.toolService.GetToolDeclarations()
	reply, fnName, fnArgs, err := uc.aiService.Generate(ctx, historyWithMemory, toolDecls)
	if err != nil {
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	var executedTools []entity.ToolEvent

	// 3. Handle Tool Execution Loop if requested by AI
	if fnName != "" {
		tResult, durationMs, tErr := uc.toolService.ExecuteTool(ctx, fnName, fnArgs)
		evt := entity.ToolEvent{
			ToolName:   fnName,
			Args:       fnArgs,
			Result:     tResult,
			DurationMs: durationMs,
		}
		if tErr != nil {
			evt.Error = tErr.Error()
		}
		executedTools = append(executedTools, evt)

		// Synthesis turn with compact sliding context
		resJSON, _ := json.Marshal(tResult)
		var synthesisContext []entity.Message
		if len(historyWithMemory) > 4 {
			synthesisContext = historyWithMemory[len(historyWithMemory)-4:]
		} else {
			synthesisContext = historyWithMemory
		}

		historyWithTool := append(synthesisContext,
			entity.Message{
				Role:       entity.RoleUser,
				Content:    fmt.Sprintf("[Resultado da ferramenta '%s']:\n```json\n%s\n```\nSintetize a resposta para o usuário de forma concisa e com o diagrama Mermaid.", fnName, string(resJSON)),
				ToolEvents: []entity.ToolEvent{evt},
			},
		)
		synthReply, _, _, sErr := uc.aiService.Generate(ctx, historyWithTool, nil)
		if sErr == nil && synthReply != "" {
			reply = synthReply
		}
	}

	duration := time.Since(start).Milliseconds()

	// 4. Persist Assistant Message
	modelMsg := entity.Message{
		ID:         fmt.Sprintf("asst-%d", time.Now().UnixNano()),
		Role:       entity.RoleModel,
		Content:    reply,
		ToolEvents: executedTools,
		CreatedAt:  time.Now(),
	}
	_ = uc.sessionRepo.AppendMessage(ctx, sessionID, modelMsg)

	output := &ChatOutput{
		SessionID:  sessionID,
		Reply:      reply,
		ToolEvents: executedTools,
		DurationMs: duration,
		Success:    true,
	}

	if cacheKey != "" && uc.cacheRepo != nil {
		uc.cacheRepo.Set(ctx, cacheKey, *output, 30*time.Minute)
	}

	return output, nil
}

// ExecuteStreamChat executes an agentic streaming conversation turn.
func (uc *ChatUseCase) ExecuteStreamChat(ctx context.Context, input ChatInput, emitter StreamEmitter) error {
	if strings.TrimSpace(input.Message) == "" && strings.TrimSpace(input.ImageBase64) == "" {
		return fmt.Errorf("message or image is required")
	}

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	// Learn user preferences from message
	uc.autoLearn(ctx, input.Message)

	startTime := time.Now()

	// Persist User Message
	userMsg := entity.Message{
		ID:            fmt.Sprintf("usr-%d", time.Now().UnixNano()),
		Role:          entity.RoleUser,
		Content:       input.Message,
		ImageBase64:   input.ImageBase64,
		ImageMimeType: input.ImageMimeType,
		CreatedAt:     time.Now(),
	}
	_ = uc.sessionRepo.AppendMessage(ctx, sessionID, userMsg)

	history, err := uc.sessionRepo.GetSessionMessages(ctx, sessionID)
	if err != nil || len(history) == 0 {
		history = input.History
	}

	// Sliding window: keep only the most recent turns to maintain low latency and low token consumption
	const maxRecentMessages = 8
	if len(history) > maxRecentMessages {
		history = history[len(history)-maxRecentMessages:]
	}

	// Send Connected Ping
	emitter.Emit(entity.StreamEvent{
		Type:      entity.EventPing,
		Content:   "connected",
		Data:      map[string]string{"sessionId": sessionID},
		Timestamp: time.Now().UnixMilli(),
	})

	// Inject Long-Term Memory Context
	historyWithMemory := uc.injectMemoryContext(ctx, history)

	toolDecls := uc.toolService.GetToolDeclarations()
	var fullText strings.Builder
	var executedTools []entity.ToolEvent

	streamErr := uc.aiService.StreamGenerate(ctx, historyWithMemory, toolDecls, func(token string, fnName string, fnArgs map[string]interface{}) error {
		// Tool call execution triggered
		if fnName != "" {
			emitter.Emit(entity.StreamEvent{
				Type:    entity.EventToolStart,
				Content: fmt.Sprintf("Executando ferramenta de arquitetura: %s", fnName),
				Data: map[string]interface{}{
					"toolName": fnName,
					"args":     fnArgs,
				},
				Timestamp: time.Now().UnixMilli(),
			})

			tResult, durationMs, tErr := uc.toolService.ExecuteTool(ctx, fnName, fnArgs)
			evt := entity.ToolEvent{
				ToolName:   fnName,
				Args:       fnArgs,
				Result:     tResult,
				DurationMs: durationMs,
			}
			if tErr != nil {
				evt.Error = tErr.Error()
			}
			executedTools = append(executedTools, evt)

			emitter.Emit(entity.StreamEvent{
				Type:      entity.EventToolResult,
				Content:   fmt.Sprintf("Ferramenta %s concluída em %dms", fnName, durationMs),
				Data:      evt,
				Timestamp: time.Now().UnixMilli(),
			})

			// Sub-stream synthesis after tool execution: use compact sliding context
			resJSON, _ := json.Marshal(tResult)
			var synthesisContext []entity.Message
			if len(historyWithMemory) > 4 {
				synthesisContext = historyWithMemory[len(historyWithMemory)-4:]
			} else {
				synthesisContext = historyWithMemory
			}

			historyWithTool := append(synthesisContext, entity.Message{
				Role:       entity.RoleUser,
				Content:    fmt.Sprintf("[Resultado da ferramenta '%s']:\n```json\n%s\n```\nSintetize a resposta para o usuário de forma concisa e com o diagrama Mermaid.", fnName, string(resJSON)),
				ToolEvents: []entity.ToolEvent{evt},
			})
			return uc.aiService.StreamGenerate(ctx, historyWithTool, nil, func(subToken string, _ string, _ map[string]interface{}) error {
				if subToken != "" {
					fullText.WriteString(subToken)
					emitter.Emit(entity.StreamEvent{
						Type:      entity.EventToken,
						Content:   subToken,
						Timestamp: time.Now().UnixMilli(),
					})
				}
				return nil
			})
		}

		if token != "" {
			fullText.WriteString(token)
			emitter.Emit(entity.StreamEvent{
				Type:      entity.EventToken,
				Content:   token,
				Timestamp: time.Now().UnixMilli(),
			})
		}
		return nil
	})

	if streamErr != nil {
		emitter.Emit(entity.StreamEvent{
			Type:      entity.EventError,
			Content:   "Erro no processamento da IA: " + streamErr.Error(),
			Timestamp: time.Now().UnixMilli(),
		})
		return streamErr
	}

	// Persist Model response
	duration := time.Since(startTime).Milliseconds()
	modelMsg := entity.Message{
		ID:         fmt.Sprintf("asst-%d", time.Now().UnixNano()),
		Role:       entity.RoleModel,
		Content:    fullText.String(),
		ToolEvents: executedTools,
		CreatedAt:  time.Now(),
	}
	_ = uc.sessionRepo.AppendMessage(ctx, sessionID, modelMsg)

	// Send Done
	emitter.Emit(entity.StreamEvent{
		Type:    entity.EventDone,
		Content: "completed",
		Data: map[string]interface{}{
			"sessionId":  sessionID,
			"durationMs": duration,
			"toolsCount": len(executedTools),
		},
		Timestamp: time.Now().UnixMilli(),
	})

	return nil
}

func (uc *ChatUseCase) injectMemoryContext(ctx context.Context, history []entity.Message) []entity.Message {
	if uc.profileRepo == nil {
		return history
	}

	profile, err := uc.profileRepo.GetProfile(ctx)
	if err != nil || profile == nil {
		return history
	}

	var parts []string
	if profile.PreferredCloud != "" {
		parts = append(parts, "Cloud: "+profile.PreferredCloud)
	}
	if len(profile.PrimaryLanguages) > 0 {
		parts = append(parts, "Stacks: "+strings.Join(profile.PrimaryLanguages, ", "))
	}
	if len(profile.PreferredDatabases) > 0 {
		parts = append(parts, "Bancos: "+strings.Join(profile.PreferredDatabases, ", "))
	}
	if len(profile.PreferredPatterns) > 0 {
		parts = append(parts, "Padrões: "+strings.Join(profile.PreferredPatterns, ", "))
	}
	if len(profile.ComplianceRules) > 0 {
		parts = append(parts, "Compliance: "+strings.Join(profile.ComplianceRules, ", "))
	}
	if len(profile.CustomNotes) > 0 {
		parts = append(parts, "Diretrizes: "+strings.Join(profile.CustomNotes, " | "))
	}

	if len(parts) == 0 {
		return history
	}

	memoryPrompt := fmt.Sprintf("[MEMÓRIA ARQUITETURAL APRENDIDA DO USUÁRIO]: %s. Alinhe suas respostas a essas preferências técnicas.", strings.Join(parts, "; "))

	// Prepend as initial context
	return append([]entity.Message{
		{
			Role:    entity.RoleUser,
			Content: memoryPrompt,
		},
		{
			Role:    entity.RoleModel,
			Content: "Entendido! Manterei as recomendações e diagramas alinhados a essas preferências técnicas.",
		},
	}, history...)
}

func (uc *ChatUseCase) autoLearn(ctx context.Context, text string) {
	if uc.profileRepo == nil || strings.TrimSpace(text) == "" {
		return
	}

	profile, err := uc.profileRepo.GetProfile(ctx)
	if err != nil || profile == nil {
		return
	}

	lower := strings.ToLower(text)
	updated := false

	// 1. Cloud preference detection
	if strings.Contains(lower, "usamos oracle") || strings.Contains(lower, "preferimos oci") || strings.Contains(lower, "oracle cloud") {
		if profile.PreferredCloud != "ORACLE" {
			profile.PreferredCloud = "ORACLE"
			updated = true
		}
	} else if strings.Contains(lower, "usamos aws") || strings.Contains(lower, "nossa nuvem é aws") || strings.Contains(lower, "preferimos aws") {
		if profile.PreferredCloud != "AWS" {
			profile.PreferredCloud = "AWS"
			updated = true
		}
	} else if strings.Contains(lower, "usamos gcp") || strings.Contains(lower, "google cloud") || strings.Contains(lower, "preferimos gcp") {
		if profile.PreferredCloud != "GCP" {
			profile.PreferredCloud = "GCP"
			updated = true
		}
	} else if strings.Contains(lower, "usamos azure") || strings.Contains(lower, "preferimos azure") {
		if profile.PreferredCloud != "Azure" {
			profile.PreferredCloud = "Azure"
			updated = true
		}
	}

	// 2. Stacks learning
	if strings.Contains(lower, "trabalhamos com go") || strings.Contains(lower, "nossa stack é go") || strings.Contains(lower, "usamos golang") {
		if !containsStr(profile.PrimaryLanguages, "Go") {
			profile.PrimaryLanguages = append(profile.PrimaryLanguages, "Go")
			updated = true
		}
	}
	if strings.Contains(lower, "trabalhamos com typescript") || strings.Contains(lower, "usamos typescript") || strings.Contains(lower, "stack é ts") {
		if !containsStr(profile.PrimaryLanguages, "TypeScript") {
			profile.PrimaryLanguages = append(profile.PrimaryLanguages, "TypeScript")
			updated = true
		}
	}
	if strings.Contains(lower, "trabalhamos com python") || strings.Contains(lower, "usamos python") {
		if !containsStr(profile.PrimaryLanguages, "Python") {
			profile.PrimaryLanguages = append(profile.PrimaryLanguages, "Python")
			updated = true
		}
	}
	if strings.Contains(lower, "trabalhamos com rust") || strings.Contains(lower, "usamos rust") {
		if !containsStr(profile.PrimaryLanguages, "Rust") {
			profile.PrimaryLanguages = append(profile.PrimaryLanguages, "Rust")
			updated = true
		}
	}

	// 3. Database learning
	if strings.Contains(lower, "usamos postgres") || strings.Contains(lower, "banco é postgresql") {
		if !containsStr(profile.PreferredDatabases, "PostgreSQL") {
			profile.PreferredDatabases = append(profile.PreferredDatabases, "PostgreSQL")
			updated = true
		}
	}
	if strings.Contains(lower, "usamos redis") || strings.Contains(lower, "cache em redis") {
		if !containsStr(profile.PreferredDatabases, "Redis") {
			profile.PreferredDatabases = append(profile.PreferredDatabases, "Redis")
			updated = true
		}
	}
	if strings.Contains(lower, "usamos clickhouse") || strings.Contains(lower, "analytics em clickhouse") {
		if !containsStr(profile.PreferredDatabases, "ClickHouse") {
			profile.PreferredDatabases = append(profile.PreferredDatabases, "ClickHouse")
			updated = true
		}
	}

	// 4. Custom corporate rule learning
	if strings.HasPrefix(lower, "lembre-se:") || strings.HasPrefix(lower, "regra:") || strings.HasPrefix(lower, "diretriz:") {
		rule := strings.TrimSpace(text)
		if !containsStr(profile.CustomNotes, rule) {
			profile.CustomNotes = append(profile.CustomNotes, rule)
			updated = true
		}
	}

	if updated {
		_ = uc.profileRepo.UpdateProfile(ctx, *profile)
	}
}

func containsStr(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
