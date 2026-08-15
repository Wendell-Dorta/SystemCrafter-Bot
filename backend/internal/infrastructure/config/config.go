package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config stores all application configuration parameters.
type Config struct {
	Port                string
	Environment         string
	AllowedOrigins      []string
	GeminiAPIKey        string
	GeminiModel         string
	RateLimitRPS        float64
	RateLimitBurst      int
	MaxBodyBytes        int64
	CacheTTLMinutes     int
	CacheCleanupMinutes int
}

// LoadConfig loads configuration from environment variables and an optional .env file.
func LoadConfig() *Config {
	loadDotEnv(".env")
	loadDotEnv("../.env")

	port := getEnv("PORT", "8080")
	env := getEnv("ENVIRONMENT", "development")
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000")
	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	geminiKey := getEnv("GEMINI_API_KEY", "")
	geminiModel := getEnv("GEMINI_MODEL", "gemini-3.5-flash-lite")

	rps, err := strconv.ParseFloat(getEnv("RATE_LIMIT_RPS", "15"), 64)
	if err != nil || rps <= 0 {
		rps = 15
	}

	burst, err := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "30"))
	if err != nil || burst <= 0 {
		burst = 30
	}

	maxBody, err := strconv.ParseInt(getEnv("MAX_BODY_BYTES", "15728640"), 10, 64) // 15MB default
	if err != nil || maxBody <= 0 {
		maxBody = 15728640
	}

	cacheTTL, err := strconv.Atoi(getEnv("CACHE_TTL_MINUTES", "60"))
	if err != nil || cacheTTL <= 0 {
		cacheTTL = 60
	}

	cacheCleanup, err := strconv.Atoi(getEnv("CACHE_CLEANUP_INTERVAL_MINUTES", "10"))
	if err != nil || cacheCleanup <= 0 {
		cacheCleanup = 10
	}

	return &Config{
		Port:                port,
		Environment:         env,
		AllowedOrigins:      origins,
		GeminiAPIKey:        geminiKey,
		GeminiModel:         geminiModel,
		RateLimitRPS:        rps,
		RateLimitBurst:      burst,
		MaxBodyBytes:        maxBody,
		CacheTTLMinutes:     cacheTTL,
		CacheCleanupMinutes: cacheCleanup,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, value)
			}
		}
	}
}
