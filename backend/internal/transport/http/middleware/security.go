package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

// SecurityOptions holds config for security headers and CORS.
type SecurityOptions struct {
	AllowedOrigins []string
	MaxBodyBytes   int64
}

// SecurityMiddleware applies CORS, standard security headers, and body size limits.
func SecurityMiddleware(opts SecurityOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if isOriginAllowed(origin, opts.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if len(opts.AllowedOrigins) == 1 && opts.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Cache-Control")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length, X-RateLimit-Limit, X-RateLimit-Remaining")

			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			if opts.MaxBodyBytes > 0 && r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, opts.MaxBodyBytes)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware catches panics and returns a clean 500 JSON without crashing.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED] %v\nStack: %s", err, string(debug.Stack()))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "An unexpected internal error occurred.",
					"status":  http.StatusInternalServerError,
					"success": false,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowedList []string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range allowedList {
		if allowed == "*" || strings.EqualFold(origin, allowed) {
			return true
		}
		if strings.HasPrefix(allowed, "*.") && strings.HasSuffix(origin, allowed[1:]) {
			return true
		}
	}
	if strings.HasSuffix(origin, ".trycloudflare.com") || strings.HasSuffix(origin, ".cloudflare.com") {
		return true
	}
	return false
}
