package http

import (
	"net/http"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/handler"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

// RouterConfig holds all use cases and middleware configs required to build the HTTP router.
type RouterConfig struct {
	ChatUseCase     *usecase.ChatUseCase
	TemplateUseCase *usecase.TemplateUseCase
	SessionUseCase  *usecase.SessionUseCase
	ProfileUseCase  *usecase.ProfileUseCase
	AuditUseCase    *usecase.CodeAuditUseCase
	HealthUseCase   *usecase.HealthUseCase
	RateLimiter     *middleware.RateLimiter
	SecurityOpts    middleware.SecurityOptions
}

// NewRouter wires up all HTTP endpoints and applies the middleware pipeline.
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Handlers
	chatStreamHandler := handler.NewChatStreamHandler(cfg.ChatUseCase)
	chatHandler := handler.NewChatHandler(cfg.ChatUseCase)
	templateHandler := handler.NewTemplateHandler(cfg.TemplateUseCase)
	sessionHandler := handler.NewSessionHandler(cfg.SessionUseCase)
	profileHandler := handler.NewProfileHandler(cfg.ProfileUseCase)
	auditHandler := handler.NewAuditHandler(cfg.AuditUseCase)
	healthHandler := handler.NewHealthHandler(cfg.HealthUseCase, cfg.RateLimiter)

	// Routes
	mux.Handle("/api/chat/stream", chatStreamHandler)
	mux.Handle("/api/chat", chatHandler)
	mux.Handle("/api/templates", templateHandler)
	mux.Handle("/api/templates/", templateHandler)
	mux.Handle("/api/sessions", sessionHandler)
	mux.Handle("/api/sessions/", sessionHandler)
	mux.Handle("/api/profile", profileHandler)
	mux.Handle("/api/profile/", profileHandler)
	mux.Handle("/api/audit", auditHandler)
	mux.Handle("/api/audit/", auditHandler)
	mux.Handle("/health", healthHandler)

	// Middleware pipeline: Recovery -> Security/CORS -> Rate Limiting -> Router
	var h http.Handler = mux
	h = cfg.RateLimiter.Middleware(h)
	h = middleware.SecurityMiddleware(cfg.SecurityOpts)(h)
	h = middleware.RecoveryMiddleware(h)

	return h
}
