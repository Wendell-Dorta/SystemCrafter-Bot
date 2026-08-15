package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/service"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/github"
	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/infrastructure/storage"
)

// CodeAuditUseCase coordinates real source code extraction, sanitization, storage, and AI-grounded architectural auditing.
type CodeAuditUseCase struct {
	aiService     service.AIService
	workspaceRepo repository.WorkspaceRepository
	sessionRepo   repository.SessionRepository
	profileRepo   repository.ProfileRepository
	githubClient  *github.GitHubClient
}

// NewCodeAuditUseCase creates an initialized CodeAuditUseCase instance.
func NewCodeAuditUseCase(
	ai service.AIService,
	workspaceRepo repository.WorkspaceRepository,
	sessionRepo repository.SessionRepository,
	profileRepo repository.ProfileRepository,
) *CodeAuditUseCase {
	return &CodeAuditUseCase{
		aiService:     ai,
		workspaceRepo: workspaceRepo,
		sessionRepo:   sessionRepo,
		profileRepo:   profileRepo,
		githubClient:  github.NewGitHubClient(),
	}
}

// AuditGitHubRepository pulls, sanitizes, saves, and performs grounded architectural audit on a public GitHub repo.
func (uc *CodeAuditUseCase) AuditGitHubRepository(ctx context.Context, repoURL, sessionID string) (*entity.CodeAuditReport, error) {
	if strings.TrimSpace(repoURL) == "" {
		return nil, fmt.Errorf("URL do repositório GitHub é obrigatória")
	}

	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	// 1. Pull & Sanitize repository files from GitHub API
	files, ignored, repoFullName, err := uc.githubClient.FetchPublicRepository(ctx, repoURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao analisar repositório GitHub: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("nenhum arquivo de código suportado foi encontrado no repositório")
	}

	// 2. Persist sanitized files to local workspace disk storage
	workspaceID := fmt.Sprintf("ws-%s", sessionID)
	if _, err := uc.workspaceRepo.SaveWorkspace(ctx, workspaceID, files); err != nil {
		return nil, fmt.Errorf("falha ao salvar arquivos do workspace no disco: %w", err)
	}

	// 3. Perform AI-grounded architectural & security review
	report, err := uc.generateGroundedAudit(ctx, files, ignored, repoFullName, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("falha na análise com a IA: %w", err)
	}

	// 4. Save report in workspace storage
	_ = uc.workspaceRepo.SaveAuditReport(ctx, workspaceID, *report)

	// 5. Inject Audit Result into the Chat Session History so user can ask follow-up questions
	uc.persistAuditToSession(ctx, sessionID, report)

	return report, nil
}

// AuditUploadedFiles sanitizes, saves, and audits user-uploaded source files.
func (uc *CodeAuditUseCase) AuditUploadedFiles(ctx context.Context, rawFiles []entity.CodeFile, sessionID, projectName string) (*entity.CodeAuditReport, error) {
	if len(rawFiles) == 0 {
		return nil, fmt.Errorf("nenhum arquivo de código foi enviado para análise")
	}

	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	if projectName == "" {
		projectName = "Projeto Local Enviado"
	}

	// 1. Sanitize & Filter sensitive files (.env, credentials, secrets)
	files, ignored := storage.SanitizeCodeFiles(rawFiles)
	if len(files) == 0 {
		return nil, fmt.Errorf("todos os arquivos enviados foram filtrados por segurança (.env/credenciais) ou formato não suportado")
	}

	// 2. Persist to local disk storage
	workspaceID := fmt.Sprintf("ws-%s", sessionID)
	if _, err := uc.workspaceRepo.SaveWorkspace(ctx, workspaceID, files); err != nil {
		return nil, fmt.Errorf("falha ao persistir arquivos no workspace: %w", err)
	}

	// 3. Perform AI-grounded review
	report, err := uc.generateGroundedAudit(ctx, files, ignored, projectName, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("falha na análise com a IA: %w", err)
	}

	// 4. Save report
	_ = uc.workspaceRepo.SaveAuditReport(ctx, workspaceID, *report)

	// 5. Inject into session history
	uc.persistAuditToSession(ctx, sessionID, report)

	return report, nil
}

func (uc *CodeAuditUseCase) GetReport(ctx context.Context, workspaceID string) (*entity.CodeAuditReport, error) {
	return uc.workspaceRepo.GetAuditReport(ctx, workspaceID)
}

func (uc *CodeAuditUseCase) generateGroundedAudit(
	ctx context.Context,
	files []entity.CodeFile,
	ignored []string,
	repoName, workspaceID string,
) (*entity.CodeAuditReport, error) {
	// Build concise code summary for prompt
	var treeSummary strings.Builder
	treeSummary.WriteString("Árvore de Arquivos Analisados:\n")
	for _, f := range files {
		treeSummary.WriteString(fmt.Sprintf("- %s (%s, %d bytes)\n", f.Path, f.Language, f.Size))
	}

	// Select top critical code snippets (e.g. main.go, routes, docker-compose, config, handlers)
	var codeSnippets strings.Builder
	codeSnippets.WriteString("\nTrechos Relevantes de Código:\n")
	snippetCount := 0
	for _, f := range files {
		base := strings.ToLower(filepath.Base(f.Path))
		isCritical := strings.Contains(base, "main") || strings.Contains(base, "route") ||
			strings.Contains(base, "handler") || strings.Contains(base, "docker") ||
			strings.Contains(base, "config") || strings.Contains(base, "server") ||
			strings.Contains(base, "app") || strings.Contains(base, "schema")

		if isCritical && snippetCount < 8 {
			contentPreview := f.Content
			if len(contentPreview) > 2500 {
				contentPreview = contentPreview[:2500] + "\n... [truncado para concisão]"
			}
			codeSnippets.WriteString(fmt.Sprintf("\n--- Arquivo: %s ---\n```%s\n%s\n```\n", f.Path, strings.ToLower(f.Language), contentPreview))
			snippetCount++
		}
	}

	prompt := fmt.Sprintf(`Você é o ArchMind AI. Realize uma Auditoria de Arquitetura e Segurança BASEADA NO CÓDIGO REAL DO PROJETO '%s'.

%s
%s

Responda OBRIGATORIAMENTE em JSON estrito com o seguinte esquema:
{
  "detectedStack": ["Go 1.24", "Next.js", "Docker", "PostgreSQL"],
  "riskScore": 25,
  "riskLevel": "LOW",
  "architectureSummary": "Visão geral da arquitetura real identificada no código...",
  "vulnerabilities": [
    {
      "id": "SEC-01",
      "severity": "HIGH",
      "category": "Auth/Network/SPOF",
      "description": "Explicação baseada no arquivo analisado...",
      "remediation": "Como corrigir no código..."
    }
  ],
  "mermaidDiagram": "graph TD\n  Client --> Gateway\n  Gateway --> API",
  "recommendations": ["Recomendação prática 1", "Recomendação prática 2"]
}`, repoName, treeSummary.String(), codeSnippets.String())

	messages := []entity.Message{
		{
			Role:    entity.RoleUser,
			Content: prompt,
		},
	}

	reply, _, _, err := uc.aiService.Generate(ctx, messages, nil)
	if err != nil {
		return nil, err
	}

	// Clean Markdown code fences from JSON response
	jsonClean := strings.TrimSpace(reply)
	jsonClean = strings.TrimPrefix(jsonClean, "```json")
	jsonClean = strings.TrimPrefix(jsonClean, "```")
	jsonClean = strings.TrimSuffix(jsonClean, "```")
	jsonClean = strings.TrimSpace(jsonClean)

	var parsed struct {
		DetectedStack       []string                 `json:"detectedStack"`
		RiskScore           int                      `json:"riskScore"`
		RiskLevel           string                   `json:"riskLevel"`
		ArchitectureSummary string                   `json:"architectureSummary"`
		Vulnerabilities     []entity.SecurityFinding `json:"vulnerabilities"`
		MermaidDiagram      string                   `json:"mermaidDiagram"`
		Recommendations     []string                 `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(jsonClean), &parsed); err != nil {
		// Fallback structuring if model returned raw text
		parsed.ArchitectureSummary = reply
		parsed.RiskScore = 30
		parsed.RiskLevel = "MEDIUM"
		parsed.DetectedStack = []string{"Detectado via Código-Fonte"}
		parsed.MermaidDiagram = "graph TD\n    Client --> Server[Servidor de Aplicação]\n    Server --> DB[(Banco de Dados)]"
	}

	report := &entity.CodeAuditReport{
		WorkspaceID:           workspaceID,
		RepositoryName:        repoName,
		AnalyzedFilesCount:    len(files),
		FilteredFilesCount:    len(ignored),
		IgnoredSensitiveFiles: ignored,
		DetectedStack:         parsed.DetectedStack,
		RiskScore:             parsed.RiskScore,
		RiskLevel:             parsed.RiskLevel,
		Vulnerabilities:       parsed.Vulnerabilities,
		ArchitectureSummary:   parsed.ArchitectureSummary,
		MermaidDiagram:        parsed.MermaidDiagram,
		Recommendations:       parsed.Recommendations,
		CreatedAt:             time.Now(),
	}

	// Update architect learned memory with detected stacks
	if uc.profileRepo != nil {
		_ = uc.profileRepo.AddNote(ctx, fmt.Sprintf("Auditado repositório '%s' com stack: %s", repoName, strings.Join(parsed.DetectedStack, ", ")))
	}

	return report, nil
}

func (uc *CodeAuditUseCase) persistAuditToSession(ctx context.Context, sessionID string, report *entity.CodeAuditReport) {
	if uc.sessionRepo == nil || report == nil {
		return
	}

	var vulnText strings.Builder
	for _, v := range report.Vulnerabilities {
		vulnText.WriteString(fmt.Sprintf("- **[%s] %s**: %s *(Correção: %s)*\n", v.Severity, v.ID, v.Description, v.Remediation))
	}

	content := fmt.Sprintf(`### 🛡️ Auditoria Arquitetural de Código: **%s**

- **Arquivos Analisados**: %d arquivos de código
- **Arquivos Sensíveis Bloqueados (.env/segredos)**: %d arquivos
- **Stack Identificada**: %s
- **Nível de Risco**: %s (Score: %d/100)

#### 📋 Visão Geral da Arquitetura:
%s

#### 🔍 Vulnerabilidades e Riscos Detectados:
%s

#### 🗺️ Diagrama da Topologia Real:
%s%s
%s
%s

#### 💡 Recomendações:
%s`,
		report.RepositoryName,
		report.AnalyzedFilesCount,
		report.FilteredFilesCount,
		strings.Join(report.DetectedStack, ", "),
		report.RiskLevel, report.RiskScore,
		report.ArchitectureSummary,
		vulnText.String(),
		"```", "mermaid",
		report.MermaidDiagram,
		"```",
		"- "+strings.Join(report.Recommendations, "\n- "),
	)

	_ = uc.sessionRepo.AppendMessage(ctx, sessionID, entity.Message{
		ID:        fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		Role:      entity.RoleModel,
		Content:   content,
		CreatedAt: time.Now(),
	})
}
