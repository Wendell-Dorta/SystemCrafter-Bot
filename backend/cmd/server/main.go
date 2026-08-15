package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/ai"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/config"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/memory"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/tools"
	transporthttp "github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("==================================================")
	log.Printf("   🚀 ArchMind AI - Strict Clean Architecture     ")
	log.Printf("   Environment: %s | Port: %s                     ", cfg.Environment, cfg.Port)
	log.Printf("   Gemini Model: %s                              ", cfg.GeminiModel)
	log.Printf("   Rate Limit: %.1f req/s (Burst: %d)            ", cfg.RateLimitRPS, cfg.RateLimitBurst)
	log.Printf("==================================================")

	if cfg.GeminiAPIKey == "" {
		log.Printf("⚠️  WARNING: GEMINI_API_KEY is not set in environment or .env file.")
		log.Printf("👉 Please provide GEMINI_API_KEY in backend/.env to enable AI generation.")
	}

	// -------------------------------------------------------------
	// 1. INFRASTRUCTURE LAYER (Adapters, Repositories, AI Service)
	// -------------------------------------------------------------
	cacheRepo, closeCache := memory.NewCacheRepository(
		time.Duration(cfg.CacheTTLMinutes)*time.Minute,
		time.Duration(cfg.CacheCleanupMinutes)*time.Minute,
	)
	defer closeCache()

	sessionRepo := memory.NewSessionRepository()
	profileRepo := memory.NewProfileRepository()
	templateRepo := memory.NewTemplateRepository()
	workspaceRepo, err := storage.NewWorkspaceRepository("storage/workspaces")
	if err != nil {
		log.Fatalf("Failed to initialize workspace storage: %v", err)
	}

	toolService := tools.NewToolService()
	aiService := ai.NewGeminiAIService(cfg.GeminiAPIKey, cfg.GeminiModel)

	// -------------------------------------------------------------
	// 2. APPLICATION USE CASES LAYER (Business Rules / SOLID)
	// -------------------------------------------------------------
	chatUseCase := usecase.NewChatUseCase(aiService, toolService, sessionRepo, cacheRepo, profileRepo)
	templateUseCase := usecase.NewTemplateUseCase(templateRepo)
	sessionUseCase := usecase.NewSessionUseCase(sessionRepo)
	profileUseCase := usecase.NewProfileUseCase(profileRepo)
	auditUseCase := usecase.NewCodeAuditUseCase(aiService, workspaceRepo, sessionRepo, profileRepo)
	healthUseCase := usecase.NewHealthUseCase(cacheRepo)

	// -------------------------------------------------------------
	// 3. TRANSPORT & MIDDLEWARE LAYER (HTTP Router & Pipeline)
	// -------------------------------------------------------------
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	securityOpts := middleware.SecurityOptions{
		AllowedOrigins: cfg.AllowedOrigins,
		MaxBodyBytes:   cfg.MaxBodyBytes,
	}

	httpHandler := transporthttp.NewRouter(transporthttp.RouterConfig{
		ChatUseCase:     chatUseCase,
		TemplateUseCase: templateUseCase,
		SessionUseCase:  sessionUseCase,
		ProfileUseCase:  profileUseCase,
		AuditUseCase:    auditUseCase,
		HealthUseCase:   healthUseCase,
		RateLimiter:     rateLimiter,
		SecurityOpts:    securityOpts,
	})

	// -------------------------------------------------------------
	// 4. HTTP SERVER & GRACEFUL SHUTDOWN
	// -------------------------------------------------------------
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpHandler,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       300 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       180 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server stopped successfully.")
}
