export type MessageRole = 'user' | 'model' | 'system';

export interface ToolEvent {
  toolName: string;
  args: Record<string, any>;
  result: any;
  durationMs: number;
  error?: string;
}

export interface Message {
  id: string;
  role: MessageRole;
  content: string;
  imageBase64?: string;
  imageMimeType?: string;
  toolEvents?: ToolEvent[];
  createdAt: string;
}

export interface Session {
  id: string;
  title: string;
  createdAt: string;
  updatedAt: string;
  messages: Message[];
}

export interface ArchitectureTemplate {
  id: string;
  title: string;
  category: string;
  description: string;
  tags: string[];
  complexity: 'Intermediate' | 'High' | 'Very High';
  recommendedStack: string[];
  mermaidDiagram: string;
  estimatedCost: string;
  keyTradeoffs: string[];
  promptStarter: string;
}

export interface CloudCostEstimate {
  provider: 'AWS' | 'GCP' | 'AZURE' | 'ORACLE' | string;
  workloadTier: string;
  monthlyTotalUSD: number;
  breakdownUSD: Record<string, number>;
  calculations: Record<string, string>;
  comparativePricing?: Record<string, number>;
  optimizations: string[];
}

export interface SecurityFinding {
  id: string;
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  category: string;
  description: string;
  remediation: string;
}

export interface SecurityAuditResult {
  riskScore: number;
  riskLevel: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  vulnerabilities: SecurityFinding[];
  complianceStatus: Record<string, string>;
  recommendations: string[];
}

export interface SystemHealth {
  status: string;
  version: string;
  uptime: string;
  totalRequests: number;
  cacheStats?: {
    hits: number;
    misses: number;
    hitRatio: number;
    itemCount: number;
    evictions: number;
  };
  rateLimitStats?: {
    activeTrackedIPs: number;
    blockedRequests: number;
    rateLimitBurst: number;
    rateLimitRPS: number;
    totalRequests: number;
  };
}

export type StreamEventType = 'token' | 'tool_start' | 'tool_result' | 'done' | 'error' | 'ping';

export interface StreamEvent {
  type: StreamEventType;
  content?: string;
  data?: any;
  timestamp: number;
}

export interface ArchitectProfile {
  id: string;
  primaryLanguages: string[];
  preferredCloud: string;
  preferredDatabases: string[];
  preferredPatterns: string[];
  complianceRules: string[];
  customNotes: string[];
  updatedAt: string;
}

export interface CodeFile {
  path: string;
  content: string;
  size?: number;
  language?: string;
}

export interface CodeAuditReport {
  workspaceId: string;
  repositoryName: string;
  analyzedFilesCount: number;
  filteredFilesCount: number;
  ignoredSensitiveFiles: string[];
  detectedStack: string[];
  riskScore: number;
  riskLevel: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  vulnerabilities: SecurityFinding[];
  architectureSummary: string;
  mermaidDiagram: string;
  recommendations: string[];
  createdAt: string;
}
