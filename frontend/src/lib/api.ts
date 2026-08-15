import {
  ArchitectureTemplate,
  ArchitectProfile,
  CodeAuditReport,
  CodeFile,
  Message,
  Session,
  SystemHealth,
} from '@/types';

function getApiBaseUrl(): string {
  if (typeof window !== 'undefined') {
    // If accessing directly on localhost:3000 standalone dev mode
    if (
      window.location.port === '3000' &&
      (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ) {
      return 'http://localhost:8080';
    }
    // Universal relative path for Cloudflare Tunnels, Nginx Gateway, Custom Domains, or LAN IP
    return '';
  }
  return process.env.INTERNAL_BACKEND_URL || 'http://backend:8080';
}

export async function fetchHealth(): Promise<{ status: SystemHealth; system: string; uptime: string } | null> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/health`, { cache: 'no-store' });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    console.error('Failed to fetch health:', error);
    return null;
  }
}

export async function fetchTemplates(): Promise<ArchitectureTemplate[]> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/templates`, { cache: 'no-store' });
    if (!res.ok) return [];
    const data = await res.json();
    return data.templates || [];
  } catch (error) {
    console.error('Failed to fetch templates:', error);
    return [];
  }
}

export async function fetchTemplateByID(id: string): Promise<ArchitectureTemplate | null> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/templates/${id}`, { cache: 'no-store' });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    console.error('Failed to fetch template by id:', error);
    return null;
  }
}

export async function fetchSessions(): Promise<Session[]> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/sessions`, { cache: 'no-store' });
    if (!res.ok) return [];
    const data = await res.json();
    return data.sessions || [];
  } catch (error) {
    console.error('Failed to fetch sessions:', error);
    return [];
  }
}

export async function fetchSessionMessages(sessionID: string): Promise<Message[]> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/sessions/${sessionID}`, { cache: 'no-store' });
    if (!res.ok) return [];
    const data = await res.json();
    return data.messages || [];
  } catch (error) {
    console.error('Failed to fetch session messages:', error);
    return [];
  }
}

export async function deleteSession(sessionID: string): Promise<boolean> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/sessions/${sessionID}`, {
      method: 'DELETE',
    });
    return res.ok;
  } catch (error) {
    console.error('Failed to delete session:', error);
    return false;
  }
}

export async function fetchProfile(): Promise<ArchitectProfile | null> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/profile`, { cache: 'no-store' });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    console.error('Failed to fetch profile:', error);
    return null;
  }
}

export async function updateProfile(profile: Partial<ArchitectProfile>): Promise<boolean> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/profile`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(profile),
    });
    return res.ok;
  } catch (error) {
    console.error('Failed to update profile:', error);
    return false;
  }
}

export async function addProfileNote(note: string): Promise<boolean> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/profile/note`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ note }),
    });
    return res.ok;
  } catch (error) {
    console.error('Failed to add profile note:', error);
    return false;
  }
}

export async function resetProfile(): Promise<boolean> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/profile/reset`, {
      method: 'POST',
    });
    return res.ok;
  } catch (error) {
    console.error('Failed to reset profile:', error);
    return false;
  }
}

export async function auditGitHubRepo(githubUrl: string, sessionId?: string): Promise<CodeAuditReport> {
  const baseUrl = getApiBaseUrl();
  const res = await fetch(`${baseUrl}/api/audit/github`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ githubUrl, sessionId }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Falha ao auditar repositório GitHub' }));
    throw new Error(err.error || 'Falha ao auditar repositório GitHub');
  }
  return await res.json();
}

export async function auditUploadedFiles(
  files: CodeFile[],
  projectName: string,
  sessionId?: string
): Promise<CodeAuditReport> {
  const baseUrl = getApiBaseUrl();
  const res = await fetch(`${baseUrl}/api/audit/upload`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ files, projectName, sessionId }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Falha ao auditar arquivos enviados' }));
    throw new Error(err.error || 'Falha ao auditar arquivos enviados');
  }
  return await res.json();
}

export async function fetchAuditReport(workspaceId: string): Promise<CodeAuditReport | null> {
  try {
    const baseUrl = getApiBaseUrl();
    const res = await fetch(`${baseUrl}/api/audit/workspaces/${workspaceId}`, { cache: 'no-store' });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    console.error('Failed to fetch audit report:', error);
    return null;
  }
}

export function getStreamEndpoint(): string {
  const baseUrl = getApiBaseUrl();
  return `${baseUrl}/api/chat/stream`;
}
