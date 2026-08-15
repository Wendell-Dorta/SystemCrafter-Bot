package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

type diskWorkspaceRepo struct {
	mu       sync.RWMutex
	baseDir  string
	reports  map[string]entity.CodeAuditReport
}

// NewWorkspaceRepository creates a WorkspaceRepository backed by local disk storage and in-memory reports cache.
func NewWorkspaceRepository(baseDir string) (repository.WorkspaceRepository, error) {
	if baseDir == "" {
		baseDir = "storage/workspaces"
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		// Graceful fallback to OS temp dir
		fallback := filepath.Join(os.TempDir(), "systemcrafter_workspaces")
		_ = os.MkdirAll(fallback, 0755)
		baseDir = fallback
	}

	return &diskWorkspaceRepo{
		baseDir: baseDir,
		reports: make(map[string]entity.CodeAuditReport),
	}, nil
}

func (r *diskWorkspaceRepo) SaveWorkspace(ctx context.Context, id string, files []entity.CodeFile) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	workspacePath := filepath.Join(r.baseDir, id)
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspace folder: %w", err)
	}

	for _, f := range files {
		// Prevent path traversal attacks
		cleanRel := filepath.Clean(f.Path)
		if filepath.IsAbs(cleanRel) || stringsHasPrefixDotDot(cleanRel) {
			cleanRel = filepath.Base(cleanRel)
		}

		fullPath := filepath.Join(workspacePath, cleanRel)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			continue
		}

		_ = os.WriteFile(fullPath, []byte(f.Content), 0644)
	}

	return workspacePath, nil
}

func (r *diskWorkspaceRepo) LoadWorkspace(ctx context.Context, id string) ([]entity.CodeFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workspacePath := filepath.Join(r.baseDir, id)
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace '%s' not found", id)
	}

	var files []entity.CodeFile
	_ = filepath.Walk(workspacePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(workspacePath, path)
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		_, lang := IsAllowedCodeFile(relPath)
		files = append(files, entity.CodeFile{
			Path:     filepath.ToSlash(relPath),
			Content:  string(contentBytes),
			Size:     info.Size(),
			Language: lang,
		})
		return nil
	})

	return files, nil
}

func (r *diskWorkspaceRepo) SaveAuditReport(ctx context.Context, id string, report entity.CodeAuditReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.reports[id] = report

	// Also persist to disk as audit_report.json
	workspacePath := filepath.Join(r.baseDir, id)
	if err := os.MkdirAll(workspacePath, 0755); err == nil {
		reportFile := filepath.Join(workspacePath, "audit_report.json")
		data, _ := json.MarshalIndent(report, "", "  ")
		_ = os.WriteFile(reportFile, data, 0644)
	}

	return nil
}

func (r *diskWorkspaceRepo) GetAuditReport(ctx context.Context, id string) (*entity.CodeAuditReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if rep, ok := r.reports[id]; ok {
		return &rep, nil
	}

	// Try reading from disk
	reportFile := filepath.Join(r.baseDir, id, "audit_report.json")
	data, err := os.ReadFile(reportFile)
	if err != nil {
		return nil, fmt.Errorf("audit report for '%s' not found", id)
	}

	var report entity.CodeAuditReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to parse audit report JSON: %w", err)
	}

	return &report, nil
}

func stringsHasPrefixDotDot(p string) bool {
	return len(p) >= 2 && p[0] == '.' && p[1] == '.'
}
