import { describe, it, beforeEach } from 'node:test';
import assert from 'node:assert/strict';
import {
  fetchHealth,
  fetchTemplates,
  fetchTemplateByID,
  fetchSessions,
  fetchSessionMessages,
  deleteSession,
  fetchProfile,
  updateProfile,
  addProfileNote,
  resetProfile,
  auditGitHubRepo,
} from './api';

describe('Frontend API Client Unit Tests', () => {
  beforeEach(() => {
    // Reset global fetch mock
    // @ts-ignore
    globalThis.fetch = async (url: string, options?: any) => {
      const urlStr = String(url);

      if (urlStr.includes('/health')) {
        return {
          ok: true,
          json: async () => ({
            status: { status: 'healthy', version: '1.0.0', uptime: '10m', totalRequests: 50, cacheStats: {}, rateLimitStats: {} },
            system: 'SystemCrafter',
            uptime: '10m',
          }),
        };
      }

      if (urlStr.includes('/api/templates/microservices-1')) {
        return {
          ok: true,
          json: async () => ({ id: 'microservices-1', title: 'Microservices Blueprint' }),
        };
      }

      if (urlStr.includes('/api/templates')) {
        return {
          ok: true,
          json: async () => ({ templates: [{ id: 't1', title: 'Template 1' }] }),
        };
      }

      if (urlStr.includes('/api/sessions') && options?.method === 'DELETE') {
        return { ok: true };
      }

      if (urlStr.includes('/api/sessions/session-123')) {
        return {
          ok: true,
          json: async () => ({ messages: [{ id: 'm1', role: 'user', content: 'Hello' }] }),
        };
      }

      if (urlStr.includes('/api/sessions')) {
        return {
          ok: true,
          json: async () => ({ sessions: [{ id: 's1', title: 'Chat 1' }] }),
        };
      }

      if (urlStr.includes('/api/profile/note')) {
        return { ok: true, json: async () => ({ success: true }) };
      }

      if (urlStr.includes('/api/profile/reset')) {
        return { ok: true, json: async () => ({ success: true }) };
      }

      if (urlStr.includes('/api/profile') && options?.method === 'PUT') {
        return { ok: true, json: async () => ({ success: true }) };
      }

      if (urlStr.includes('/api/profile')) {
        return {
          ok: true,
          json: async () => ({ preferredCloud: 'AWS', primaryLanguages: ['Go'] }),
        };
      }

      if (urlStr.includes('/api/audit/github')) {
        return {
          ok: true,
          json: async () => ({
            repositoryName: 'zen-browser/desktop',
            riskScore: 15,
            riskLevel: 'LOW',
            detectedStack: ['TypeScript', 'Gecko'],
          }),
        };
      }

      return { ok: false, status: 404, json: async () => ({ error: 'Not Found' }) };
    };
  });

  it('fetchHealth() returns system health data', async () => {
    const res = await fetchHealth();
    assert.ok(res !== null);
    assert.equal(res?.system, 'SystemCrafter');
  });

  it('fetchTemplates() returns template array', async () => {
    const templates = await fetchTemplates();
    assert.equal(templates.length, 1);
    assert.equal(templates[0].id, 't1');
  });

  it('fetchTemplateByID() returns specific template', async () => {
    const tmpl = await fetchTemplateByID('microservices-1');
    assert.ok(tmpl !== null);
    assert.equal(tmpl?.id, 'microservices-1');
  });

  it('fetchSessions() returns session list', async () => {
    const sessions = await fetchSessions();
    assert.equal(sessions.length, 1);
    assert.equal(sessions[0].id, 's1');
  });

  it('fetchSessionMessages() returns message list', async () => {
    const messages = await fetchSessionMessages('session-123');
    assert.equal(messages.length, 1);
    assert.equal(messages[0].content, 'Hello');
  });

  it('deleteSession() returns true on success', async () => {
    const ok = await deleteSession('session-123');
    assert.equal(ok, true);
  });

  it('fetchProfile() and updateProfile() work correctly', async () => {
    const prof = await fetchProfile();
    assert.equal(prof?.preferredCloud, 'AWS');

    const updateOk = await updateProfile({
      id: 'default',
      primaryLanguages: ['Go', 'TypeScript'],
      preferredCloud: 'GCP',
      preferredDatabases: ['PostgreSQL'],
      preferredPatterns: ['CQRS'],
      complianceRules: ['LGPD'],
      customNotes: [],
      updatedAt: new Date().toISOString(),
    });
    assert.equal(updateOk, true);
  });

  it('addProfileNote() and resetProfile() work correctly', async () => {
    const noteOk = await addProfileNote('Prefere Kubernetes');
    assert.equal(noteOk, true);

    const resetOk = await resetProfile();
    assert.equal(resetOk, true);
  });

  it('auditGitHubRepo() returns code audit findings', async () => {
    const report = await auditGitHubRepo('https://github.com/zen-browser/desktop', 'session-123');
    assert.equal(report.repositoryName, 'zen-browser/desktop');
    assert.equal(report.riskScore, 15);
    assert.equal(report.riskLevel, 'LOW');
  });
});
