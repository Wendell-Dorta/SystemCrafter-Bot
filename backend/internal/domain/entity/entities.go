package entity

import "time"

// Role constants
const (
	RoleUser   = "user"
	RoleModel  = "model"
	RoleSystem = "system"
)

// Message represents a message entity in a conversation.
type Message struct {
	ID            string      `json:"id"`
	Role          string      `json:"role"` // "user", "model", "system"
	Content       string      `json:"content"`
	ImageBase64   string      `json:"imageBase64,omitempty"`
	ImageMimeType string      `json:"imageMimeType,omitempty"`
	ToolEvents    []ToolEvent `json:"toolEvents,omitempty"`
	CreatedAt     time.Time   `json:"createdAt"`
}

// ToolEvent records a tool call execution and its result.
type ToolEvent struct {
	ToolName   string      `json:"toolName"`
	Args       interface{} `json:"args"`
	Result     interface{} `json:"result"`
	DurationMs int64       `json:"durationMs"`
	Error      string      `json:"error,omitempty"`
}

// Session represents a stored chat session entity.
type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Messages  []Message `json:"messages"`
}

// ArchitectureTemplate represents a reference system blueprint entity.
type ArchitectureTemplate struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Tags             []string `json:"tags"`
	Complexity       string   `json:"complexity"`
	RecommendedStack []string `json:"recommendedStack"`
	MermaidDiagram   string   `json:"mermaidDiagram"`
	EstimatedCost    string   `json:"estimatedCost"`
	KeyTradeoffs     []string `json:"keyTradeoffs"`
	PromptStarter    string   `json:"promptStarter"`
}

// CloudCostEstimate represents calculated cloud infrastructure pricing entity.
type CloudCostEstimate struct {
	Provider           string             `json:"provider"`
	WorkloadTier       string             `json:"workloadTier"`
	MonthlyTotalUSD    float64            `json:"monthlyTotalUSD"`
	BreakdownUSD       map[string]float64 `json:"breakdownUSD"`
	Calculations       map[string]string  `json:"calculations"`
	ComparativePricing map[string]float64 `json:"comparativePricing,omitempty"`
	Optimizations      []string           `json:"optimizations"`
}

// SecurityAuditResult represents security findings and compliance status entity.
type SecurityAuditResult struct {
	RiskScore        int               `json:"riskScore"`
	RiskLevel        string            `json:"riskLevel"`
	Vulnerabilities  []SecurityFinding `json:"vulnerabilities"`
	ComplianceStatus map[string]string `json:"complianceStatus"`
	Recommendations  []string          `json:"recommendations"`
}

// SecurityFinding describes an architectural vulnerability entity.
type SecurityFinding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
}

// TechStackMatrix represents architectural tradeoffs entity.
type TechStackMatrix struct {
	WorkloadType  string            `json:"workloadType"`
	Recommended   map[string]string `json:"recommended"`
	TradeoffNotes map[string]string `json:"tradeoffNotes"`
	Alternatives  map[string]string `json:"alternatives"`
}

// ArchitecturePattern represents reference architectural pattern entity.
type ArchitecturePattern struct {
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	Summary        string   `json:"summary"`
	WhenToUse      []string `json:"whenToUse"`
	Pitfalls       []string `json:"pitfalls"`
	MermaidDiagram string   `json:"mermaidDiagram"`
}

// StreamEventType defines SSE stream event types.
type StreamEventType string

const (
	EventToken      StreamEventType = "token"
	EventToolStart  StreamEventType = "tool_start"
	EventToolResult StreamEventType = "tool_result"
	EventDone       StreamEventType = "done"
	EventError      StreamEventType = "error"
	EventPing       StreamEventType = "ping"
)

// StreamEvent represents an outgoing SSE packet entity.
type StreamEvent struct {
	Type      StreamEventType `json:"type"`
	Content   string          `json:"content,omitempty"`
	Data      interface{}     `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

// SystemHealth represents health metrics entity.
type SystemHealth struct {
	Status         string                 `json:"status"`
	Version        string                 `json:"version"`
	Uptime         string                 `json:"uptime"`
	TotalRequests  int64                  `json:"totalRequests"`
	CacheStats     map[string]interface{} `json:"cacheStats"`
	RateLimitStats map[string]interface{} `json:"rateLimitStats"`
}

// ArchitectProfile represents learned cross-session engineering preferences and decisions.
type ArchitectProfile struct {
	ID                 string    `json:"id"`
	PrimaryLanguages   []string  `json:"primaryLanguages"`
	PreferredCloud     string    `json:"preferredCloud"`
	PreferredDatabases []string  `json:"preferredDatabases"`
	PreferredPatterns  []string  `json:"preferredPatterns"`
	ComplianceRules    []string  `json:"complianceRules"`
	CustomNotes        []string  `json:"customNotes"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// CodeFile represents an ingested source code file for architectural review.
type CodeFile struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
	Language string `json:"language"`
}

// CodeAuditRequest represents an incoming code audit trigger.
type CodeAuditRequest struct {
	SessionID string     `json:"sessionId"`
	GitHubURL string     `json:"githubUrl,omitempty"`
	Files     []CodeFile `json:"files,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

// CodeAuditReport represents the grounded architectural and security audit findings.
type CodeAuditReport struct {
	WorkspaceID           string            `json:"workspaceId"`
	RepositoryName        string            `json:"repositoryName"`
	AnalyzedFilesCount    int               `json:"analyzedFilesCount"`
	FilteredFilesCount    int               `json:"filteredFilesCount"`
	IgnoredSensitiveFiles []string          `json:"ignoredSensitiveFiles"`
	DetectedStack         []string          `json:"detectedStack"`
	RiskScore             int               `json:"riskScore"`
	RiskLevel             string            `json:"riskLevel"`
	Vulnerabilities       []SecurityFinding `json:"vulnerabilities"`
	ArchitectureSummary   string            `json:"architectureSummary"`
	MermaidDiagram        string            `json:"mermaidDiagram"`
	Recommendations       []string          `json:"recommendations"`
	CreatedAt             time.Time         `json:"createdAt"`
}
