package storage

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

var (
	// Regex patterns for secret token redaction inside code contents
	secretRegexes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|secret|password|token|auth[_-]?token)\s*[:=]\s*["']([^"']{8,})["']`),
		regexp.MustCompile(`(?i)(AIzaSy[0-9A-Za-z-_]{33})`),                                    // Google API Key
		regexp.MustCompile(`(?i)(ghp_[0-9A-Za-z]{36}|github_pat_[0-9A-Za-z_]{82})`),           // GitHub Token
		regexp.MustCompile(`(?i)(xox[baprs]-[0-9A-Za-z-]{10,48})`),                             // Slack Token
		regexp.MustCompile(`(?i)(postgres|mysql|mongodb(?:\+srv)?):\/\/[^:\s]+:([^@\s]+)@`),   // DB Connection strings
	}

	// Supported code extensions for architectural review
	allowedExtensions = map[string]string{
		".go":         "Go",
		".ts":         "TypeScript",
		".tsx":        "TypeScript (React)",
		".js":         "JavaScript",
		".jsx":        "JavaScript (React)",
		".py":         "Python",
		".java":       "Java",
		".rs":         "Rust",
		".cs":         "C#",
		".cpp":        "C++",
		".c":          "C",
		".sql":        "SQL",
		".yaml":       "YAML",
		".yml":        "YAML",
		".json":       "JSON",
		".toml":       "TOML",
		".proto":      "Protocol Buffers",
		".dockerfile": "Dockerfile",
		".md":         "Markdown",
	}
)

// IsSensitivePath checks if a relative file path matches sensitive patterns that must not be ingested.
func IsSensitivePath(path string) (bool, string) {
	norm := strings.ToLower(filepath.ToSlash(path))
	base := strings.ToLower(filepath.Base(norm))

	// 1. Blacklisted Directories
	if strings.Contains(norm, "/node_modules/") || strings.HasPrefix(norm, "node_modules/") ||
		strings.Contains(norm, "/vendor/") || strings.HasPrefix(norm, "vendor/") ||
		strings.Contains(norm, "/.git/") || strings.HasPrefix(norm, ".git/") ||
		strings.Contains(norm, "/.next/") || strings.HasPrefix(norm, ".next/") ||
		strings.Contains(norm, "/dist/") || strings.HasPrefix(norm, "dist/") ||
		strings.Contains(norm, "/build/") || strings.HasPrefix(norm, "build/") {
		return true, "Diretório de dependências/build ignorado"
	}

	// 2. Secret & Environment Files
	if strings.HasPrefix(base, ".env") || strings.HasSuffix(base, ".env") {
		return true, "Arquivo de variáveis de ambiente (.env) bloqueado por segurança"
	}

	if strings.Contains(base, "credential") || strings.Contains(base, "secret") ||
		strings.Contains(base, "service-account") || strings.Contains(base, "id_rsa") ||
		strings.Contains(base, "id_ed25519") || strings.HasSuffix(base, ".pem") ||
		strings.HasSuffix(base, ".key") || strings.HasSuffix(base, ".pfx") ||
		strings.HasSuffix(base, ".p12") || strings.HasSuffix(base, ".htpasswd") {
		return true, "Chave criptográfica/Credencial bloqueada por segurança"
	}

	// 3. Binaries & Media Files
	ext := filepath.Ext(base)
	binaryExts := map[string]bool{
		".exe": true, ".dll": true, ".so": true, ".dylib": true, ".bin": true,
		".zip": true, ".tar": true, ".gz": true, ".7z": true, ".rar": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
		".ico": true, ".pdf": true, ".mp4": true, ".mp3": true, ".wasm": true,
		".class": true, ".o": true, ".obj": true,
	}
	if binaryExts[ext] {
		return true, "Arquivo binário/mídia não processável"
	}

	return false, ""
}

// IsAllowedCodeFile checks if a file is a valid source code / config file.
func IsAllowedCodeFile(path string) (bool, string) {
	base := strings.ToLower(filepath.Base(path))
	if base == "dockerfile" || base == "docker-compose.yml" || base == "docker-compose.yaml" ||
		base == "makefile" || base == "go.mod" || base == "package.json" || base == "cargo.toml" ||
		base == "pom.xml" || base == "requirements.txt" {
		return true, "Config/Build"
	}

	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := allowedExtensions[ext]; ok {
		return true, lang
	}

	return false, ""
}

// SanitizeContent removes/masks sensitive API tokens and credentials from code contents.
func SanitizeContent(content string) string {
	for _, re := range secretRegexes {
		content = re.ReplaceAllString(content, "$1: \"[REDACTED_SECRET]\"")
	}
	return content
}

// SanitizeCodeFiles filters out sensitive files and redacts secrets from acceptable code files.
func SanitizeCodeFiles(files []entity.CodeFile) ([]entity.CodeFile, []string) {
	var sanitized []entity.CodeFile
	var ignored []string

	for _, f := range files {
		// Check sensitive path
		if sensitive, reason := IsSensitivePath(f.Path); sensitive {
			ignored = append(ignored, f.Path+" ("+reason+")")
			continue
		}

		// Check allowed code extensions
		allowed, lang := IsAllowedCodeFile(f.Path)
		if !allowed {
			ignored = append(ignored, f.Path+" (Formato de arquivo não suportado)")
			continue
		}

		// Limit individual file size to 350KB
		if len(f.Content) > 350*1024 {
			ignored = append(ignored, f.Path+" (Arquivo excede o limite de 350KB)")
			continue
		}

		cleanContent := SanitizeContent(f.Content)
		sanitized = append(sanitized, entity.CodeFile{
			Path:     f.Path,
			Content:  cleanContent,
			Size:     int64(len(cleanContent)),
			Language: lang,
		})
	}

	return sanitized, ignored
}
