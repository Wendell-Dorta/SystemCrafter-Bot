package unit_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/transport/http/middleware"
)

func TestRateLimiter_AllowAndBlock(t *testing.T) {
	rl := middleware.NewRateLimiter(5.0, 3)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := rl.Middleware(dummyHandler)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.50:12345"
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Request %d should be allowed, got status %d", i+1, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.50:12345"
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected status 429 Too Many Requests, got %d", rec.Code)
	}

	time.Sleep(250 * time.Millisecond)

	recRefill := httptest.NewRecorder()
	wrapped.ServeHTTP(recRefill, req)

	if recRefill.Code != http.StatusOK {
		t.Fatalf("Refilled token request should be allowed, got status %d", recRefill.Code)
	}
}

func TestSecurityMiddleware_HeadersAndCORS(t *testing.T) {
	opts := middleware.SecurityOptions{
		AllowedOrigins: []string{"http://localhost:3000"},
		MaxBodyBytes:   1024,
	}

	handler := middleware.SecurityMiddleware(opts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected CORS origin http://localhost:3000, got '%s'", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options: nosniff")
	}
}

func TestSecurityMiddleware_MaxBodyBytes(t *testing.T) {
	opts := middleware.SecurityOptions{
		AllowedOrigins: []string{"*"},
		MaxBodyBytes:   10,
	}

	handler := middleware.SecurityMiddleware(opts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Payload too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	largeBody := bytes.NewReader([]byte(strings.Repeat("A", 50)))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", largeBody)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("Expected 413 Payload Too Large, got %d", rec.Code)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Simulated critical error")
	})

	wrapped := middleware.RecoveryMiddleware(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500 for recovered panic, got %d", rec.Code)
	}
}
