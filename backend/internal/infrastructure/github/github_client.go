package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
)

// GitHubClient handles pulling public repository code trees for architectural analysis.
type GitHubClient struct {
	httpClient *http.Client
}

// NewGitHubClient creates an initialized GitHub client.
func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

type repoMetadataResponse struct {
	DefaultBranch string `json:"default_branch"`
	Message       string `json:"message"`
}

type treeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
	Size int64  `json:"size"`
}

type gitTreeResponse struct {
	SHA       string      `json:"sha"`
	Tree      []treeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
	Message   string      `json:"message"`
}

// ParseRepoURL parses "https://github.com/owner/repo" into owner, repo and optional specified branch.
func ParseRepoURL(rawURL string) (string, string, string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("URL inválida: %w", err)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("URL deve estar no formato https://github.com/owner/repo")
	}

	owner := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")
	specifiedBranch := ""
	if len(parts) >= 4 && (parts[2] == "tree" || parts[2] == "blob") {
		specifiedBranch = parts[3]
	}

	return owner, repo, specifiedBranch, nil
}

// FetchPublicRepository fetches code files from a public GitHub repository with dynamic default branch detection.
func (c *GitHubClient) FetchPublicRepository(ctx context.Context, repoURL string) ([]entity.CodeFile, []string, string, error) {
	owner, repo, specifiedBranch, err := ParseRepoURL(repoURL)
	if err != nil {
		return nil, nil, "", err
	}

	// 1. Resolve Target Branch dynamically
	targetBranch := specifiedBranch
	if targetBranch == "" {
		// Detect default branch (e.g. 'dev', 'main', 'master', 'develop') from repository metadata
		detectedBranch, bErr := c.detectDefaultBranch(ctx, owner, repo)
		if bErr == nil && detectedBranch != "" {
			targetBranch = detectedBranch
		} else {
			targetBranch = "main" // Initial guess
		}
	}

	// 2. Fetch Git Tree with branch fallbacks
	candidateBranches := []string{targetBranch, "dev", "main", "master", "develop", "trunk"}
	var treeResp *gitTreeResponse
	var activeBranch string
	var lastStatus int

	for _, branchToTry := range uniqueStrings(candidateBranches) {
		apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repo, branchToTry)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "SystemCrafter-Architect/1.0")
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			continue
		}

		lastStatus = resp.StatusCode
		if resp.StatusCode == http.StatusOK {
			var tr gitTreeResponse
			if err := json.NewDecoder(resp.Body).Decode(&tr); err == nil && len(tr.Tree) > 0 {
				treeResp = &tr
				activeBranch = branchToTry
				resp.Body.Close()
				break
			}
		}
		resp.Body.Close()
	}

	if treeResp == nil {
		if lastStatus == http.StatusForbidden {
			return nil, nil, "", fmt.Errorf("limite de taxa da API pública do GitHub excedido (HTTP 403 Rate Limit). Tente novamente em alguns minutos ou faça upload da pasta de código local")
		}
		return nil, nil, "", fmt.Errorf("repositório não encontrado ou nenhuma branch acessível (HTTP %d). Verifique se a URL 'https://github.com/%s/%s' é pública", lastStatus, owner, repo)
	}

	var rawFiles []entity.CodeFile
	maxFilesToFetch := 35
	fetchedCount := 0

	// 3. Filter tree and fetch contents of key architecture files
	for _, entry := range treeResp.Tree {
		if entry.Type != "blob" {
			continue
		}

		// Security check: Ignore secrets, .env, keys, credentials
		if sensitive, _ := storage.IsSensitivePath(entry.Path); sensitive {
			continue
		}

		allowed, lang := storage.IsAllowedCodeFile(entry.Path)
		if !allowed {
			continue
		}

		if fetchedCount >= maxFilesToFetch {
			break
		}

		// Fetch raw content
		rawUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, activeBranch, entry.Path)
		fileReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, rawUrl, nil)
		fileReq.Header.Set("User-Agent", "SystemCrafter-Architect/1.0")
		fileResp, err := c.httpClient.Do(fileReq)
		if err != nil || fileResp.StatusCode != http.StatusOK {
			continue
		}

		contentBytes, _ := io.ReadAll(io.LimitReader(fileResp.Body, 350*1024))
		fileResp.Body.Close()

		rawFiles = append(rawFiles, entity.CodeFile{
			Path:     entry.Path,
			Content:  string(contentBytes),
			Size:     int64(len(contentBytes)),
			Language: lang,
		})
		fetchedCount++
	}

	// 4. Sanitize and redact
	sanitized, ignored := storage.SanitizeCodeFiles(rawFiles)
	repoFullName := fmt.Sprintf("%s/%s (branch: %s)", owner, repo, activeBranch)

	return sanitized, ignored, repoFullName, nil
}

// detectDefaultBranch calls GitHub repo metadata API to determine default branch (e.g. 'dev', 'main', 'master').
func (c *GitHubClient) detectDefaultBranch(ctx context.Context, owner, repo string) (string, error) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "SystemCrafter-Architect/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var meta repoMetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", err
	}

	return strings.TrimSpace(meta.DefaultBranch), nil
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
