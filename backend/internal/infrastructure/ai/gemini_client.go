package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	GeminiBaseURL = "https://generativelanguage.googleapis.com/v1beta"
)

// GeminiClient interacts with Google Gemini API with fallback resilience and rate-limit backoff.
type GeminiClient struct {
	apiKey         string
	model          string
	fallbackModels []string
	httpClient     *http.Client
}

// NewGeminiClient creates an initialized Gemini client.
func NewGeminiClient(apiKey, model string) *GeminiClient {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
		fallbackModels: []string{
			"gemini-2.5-flash",
			"gemini-2.0-flash",
			"gemini-1.5-flash",
			"gemini-2.0-flash-lite",
		},
		httpClient: &http.Client{
			Timeout: 300 * time.Second, // 5 minutes timeout for complex multi-language code generation
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       90 * time.Second,
				ResponseHeaderTimeout: 60 * time.Second,
			},
		},
	}
}

// GetSystemInstruction returns a token-efficient, highly direct persona prompt.
func GetSystemInstruction() Content {
	return Content{
		Role: "system",
		Parts: []Part{
			{
				Text: `Você é o **ArchMind AI**, Arquiteto de Software Principal e Especialista em Computação em Nuvem (AWS, GCP, Azure, Oracle Cloud OCI), Sistemas Distribuídos e FinOps.

### REGRAS OPERACIONAIS:
1. **Respostas Diretas & Alto Valor**: Seja conciso, técnico e objetivo. Priorize tabelas de trade-offs, diagramas e recomendações práticas. Evite saudações ou explicações genéricas.
2. **FinOps & Custos Multi-Cloud**:
   - Suporte completo: **AWS, GCP, Azure e Oracle Cloud (OCI)**.
   - Se o usuário pedir custos sem especificar volume, faça uma estimativa preliminar e pergunte os 5 pontos-chave para calibrar (Volume reqs/mês, Compute Serverless vs K8s vs VMs, Banco/Storage GB, Provedor e Cache/Egress).
   - Com os dados fornecidos, invoque 'estimate_cloud_costs'. Apresente decomposição por serviço, comparativo multi-cloud e estratégias de economia (ARM Graviton/Ampere, Savings Plans, Lifecycle).
3. **Ferramentas Especializadas**:
   - Riscos, SPOF e conformidade (LGPD, PCI-DSS, OWASP): 'audit_security_compliance'.
   - Matriz de tecnologias/linguagens: 'generate_tech_stack_matrix'.
   - Padrões arquiteturais (Outbox, CQRS, Saga): 'lookup_architecture_patterns'.
4. **Diagramas Mermaid Válidos**:
   - Sempre que propor ou analisar arquiteturas, gere um diagrama Mermaid limpo ('graph TD', 'sequenceDiagram').
   - REGRA DE SINTAXE: Sempre use aspas duplas em rótulos com caracteres especiais: NodeA["API Gateway (Go)"] --> NodeB["PostgreSQL (HA)"].
5. **Idioma**: Português do Brasil com terminologia técnica precisa.`,
			},
		},
	}
}

// GenerateContent sends a non-streaming request with automatic model fallback and backoff.
func (c *GeminiClient) GenerateContent(ctx context.Context, contents []Content, toolDecls []FunctionDeclaration) (*GenerateContentResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	modelsToTry := append([]string{c.model}, c.fallbackModels...)
	var lastErr error

	for idx, modelName := range uniqueStrings(modelsToTry) {
		if idx > 0 {
			// Small backoff before trying fallback model to prevent immediate rate limit cascade
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(idx*800) * time.Millisecond):
			}
		}

		resp, err := c.doGenerateContent(ctx, modelName, contents, toolDecls)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		// Retry on 404, 429 quota exhaustion or 503 service unavailable
		if !isRetryableError(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("all models failed, last error: %w", lastErr)
}

func (c *GeminiClient) doGenerateContent(ctx context.Context, model string, contents []Content, toolDecls []FunctionDeclaration) (*GenerateContentResponse, error) {
	sysInst := GetSystemInstruction()
	reqBody := GenerateContentRequest{
		Contents:          contents,
		SystemInstruction: &sysInst,
		GenerationConfig: &GenerationConfig{
			Temperature:     0.3,
			TopP:            0.9,
			MaxOutputTokens: 8192,
		},
	}

	if len(toolDecls) > 0 {
		reqBody.Tools = []ToolContainer{
			{FunctionDeclarations: toolDecls},
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", GeminiBaseURL, model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API error (HTTP %d, model %s): %s", resp.StatusCode, model, string(bodyBytes))
	}

	var genResp GenerateContentResponse
	if err := json.Unmarshal(bodyBytes, &genResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &genResp, nil
}

// StreamGenerateContent sends a streaming SSE request with fallback and rate-limit recovery.
func (c *GeminiClient) StreamGenerateContent(
	ctx context.Context,
	contents []Content,
	toolDecls []FunctionDeclaration,
	chunkHandler func(chunk *GenerateContentResponse) error,
) error {
	if c.apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	modelsToTry := append([]string{c.model}, c.fallbackModels...)
	var lastErr error

	for idx, modelName := range uniqueStrings(modelsToTry) {
		if idx > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(idx*1000) * time.Millisecond):
			}
		}

		err := c.doStreamGenerateContent(ctx, modelName, contents, toolDecls, chunkHandler)
		if err == nil {
			return nil
		}
		lastErr = err
		if !isRetryableError(err) {
			return err
		}
	}

	return fmt.Errorf("all streaming models failed, last error: %w", lastErr)
}

func (c *GeminiClient) doStreamGenerateContent(
	ctx context.Context,
	model string,
	contents []Content,
	toolDecls []FunctionDeclaration,
	chunkHandler func(chunk *GenerateContentResponse) error,
) error {
	sysInst := GetSystemInstruction()
	reqBody := GenerateContentRequest{
		Contents:          contents,
		SystemInstruction: &sysInst,
		GenerationConfig: &GenerationConfig{
			Temperature:     0.3,
			TopP:            0.9,
			MaxOutputTokens: 8192,
		},
	}

	if len(toolDecls) > 0 {
		reqBody.Tools = []ToolContainer{
			{FunctionDeclarations: toolDecls},
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal streaming request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?alt=sse&key=%s", GeminiBaseURL, model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create streaming request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gemini stream request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gemini streaming API error (HTTP %d, model %s): %s", resp.StatusCode, model, string(errBytes))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("stream read error: %w", err)
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" || data == "[DONE]" {
				continue
			}

			var chunk GenerateContentResponse
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if err := chunkHandler(&chunk); err != nil {
				return err
			}
		}
	}

	return nil
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "500") ||
		strings.Contains(msg, "RESOURCE_EXHAUSTED") ||
		strings.Contains(msg, "quota") ||
		strings.Contains(msg, "rate limit")
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
