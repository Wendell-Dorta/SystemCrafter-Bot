'use client';

import React, { useState, useEffect } from 'react';
import { ArchitectProfile } from '@/types';
import { fetchProfile, updateProfile, resetProfile } from '@/lib/api';
import { Brain, Cloud, Database, Shield, Layers, Plus, Trash2, X, Check, RefreshCw } from 'lucide-react';

interface ProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const ProfileModal: React.FC<ProfileModalProps> = ({ isOpen, onClose }) => {
  const [profile, setProfile] = useState<ArchitectProfile | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [newNote, setNewNote] = useState('');
  const [newLanguage, setNewLanguage] = useState('');
  const [newDB, setNewDB] = useState('');
  const [newPattern, setNewPattern] = useState('');
  const [newCompliance, setNewCompliance] = useState('');
  const [savedSuccess, setSavedSuccess] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      loadProfile();
    }
  }, [isOpen]);

  const loadProfile = async () => {
    setLoading(true);
    setErrorMessage(null);
    try {
      const data = await fetchProfile();
      if (data) {
        setProfile({
          id: data.id || 'default-architect',
          preferredCloud: data.preferredCloud || 'Multi-Cloud',
          primaryLanguages: data.primaryLanguages || [],
          preferredDatabases: data.preferredDatabases || [],
          preferredPatterns: data.preferredPatterns || [],
          complianceRules: data.complianceRules || [],
          customNotes: data.customNotes || [],
          updatedAt: data.updatedAt || new Date().toISOString(),
        });
      } else {
        setProfile({
          id: 'default-architect',
          preferredCloud: 'Multi-Cloud',
          primaryLanguages: ['Go', 'TypeScript'],
          preferredDatabases: ['PostgreSQL', 'Redis'],
          preferredPatterns: ['Event-Driven', 'CQRS', 'Transactional Outbox'],
          complianceRules: ['LGPD', 'OWASP Top 10'],
          customNotes: [
            'Prioriza simplicidade operacional, alta coesão e baixo acoplamento.',
            'Prefere soluções cloud native gerenciadas e contêineres Docker.',
          ],
          updatedAt: new Date().toISOString(),
        });
      }
    } catch (err) {
      console.error('Failed to load profile:', err);
      setErrorMessage('Erro ao carregar perfil de memória');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  const handleSave = async () => {
    if (!profile) return;
    setSaving(true);
    setErrorMessage(null);
    try {
      const payload: ArchitectProfile = {
        id: profile.id || 'default-architect',
        preferredCloud: profile.preferredCloud || '',
        primaryLanguages: profile.primaryLanguages || [],
        preferredDatabases: profile.preferredDatabases || [],
        preferredPatterns: profile.preferredPatterns || [],
        complianceRules: profile.complianceRules || [],
        customNotes: profile.customNotes || [],
        updatedAt: new Date().toISOString(),
      };
      const ok = await updateProfile(payload);
      if (ok) {
        setSavedSuccess(true);
        setTimeout(() => setSavedSuccess(false), 2500);
      } else {
        setErrorMessage('Não foi possível salvar as preferências.');
      }
    } catch (err) {
      console.error('Failed to save profile:', err);
      setErrorMessage('Erro ao conectar com o servidor.');
    } finally {
      setSaving(false);
    }
  };

  const handleReset = async () => {
    if (!confirm('Deseja realmente limpar toda a memória de preferências aprendidas?')) return;
    setSaving(true);
    setErrorMessage(null);
    try {
      const ok = await resetProfile();
      if (ok) {
        setProfile({
          id: 'default-architect',
          preferredCloud: '',
          primaryLanguages: [],
          preferredDatabases: [],
          preferredPatterns: [],
          complianceRules: [],
          customNotes: [],
          updatedAt: new Date().toISOString(),
        });
        setSavedSuccess(true);
        setTimeout(() => setSavedSuccess(false), 2500);
      } else {
        setErrorMessage('Falha ao resetar memória.');
      }
    } catch (err) {
      console.error('Failed to reset profile:', err);
      setErrorMessage('Erro ao comunicar com o servidor.');
    } finally {
      setSaving(false);
    }
  };

  const addTag = (
    field: 'primaryLanguages' | 'preferredDatabases' | 'preferredPatterns' | 'complianceRules',
    val: string,
    setVal: (v: string) => void
  ) => {
    const trimmed = val.trim();
    if (!trimmed) return;
    setProfile((prev) => {
      const base = prev || {
        id: 'default-architect',
        preferredCloud: 'Multi-Cloud',
        primaryLanguages: [],
        preferredDatabases: [],
        preferredPatterns: [],
        complianceRules: [],
        customNotes: [],
        updatedAt: new Date().toISOString(),
      };
      const current = base[field] || [];
      if (current.includes(trimmed)) return base;
      return {
        ...base,
        [field]: [...current, trimmed],
      };
    });
    setVal('');
  };

  const removeTag = (
    field: 'primaryLanguages' | 'preferredDatabases' | 'preferredPatterns' | 'complianceRules',
    val: string
  ) => {
    setProfile((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        [field]: (prev[field] || []).filter((item) => item !== val),
      };
    });
  };

  const addNote = () => {
    const trimmed = newNote.trim();
    if (!trimmed) return;
    setProfile((prev) => {
      const base = prev || {
        id: 'default-architect',
        preferredCloud: 'Multi-Cloud',
        primaryLanguages: [],
        preferredDatabases: [],
        preferredPatterns: [],
        complianceRules: [],
        customNotes: [],
        updatedAt: new Date().toISOString(),
      };
      return {
        ...base,
        customNotes: [...(base.customNotes || []), trimmed],
      };
    });
    setNewNote('');
  };

  const removeNote = (index: number) => {
    setProfile((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        customNotes: (prev.customNotes || []).filter((_, i) => i !== index),
      };
    });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-2.5 sm:p-4 bg-black/80 backdrop-blur-md animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-2xl max-h-[92dvh] sm:max-h-[88vh] flex flex-col shadow-2xl overflow-hidden text-slate-200">
        {/* Header */}
        <div className="px-4 py-3 sm:px-6 sm:py-4 border-b border-slate-800 flex items-center justify-between bg-slate-950/70 flex-shrink-0">
          <div className="flex items-center gap-2.5 overflow-hidden">
            <div className="w-8 h-8 sm:w-9 sm:h-9 rounded-xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400 flex-shrink-0">
              <Brain className="w-4 h-4 sm:w-5 sm:h-5" />
            </div>
            <div className="overflow-hidden">
              <h2 className="text-sm sm:text-base font-semibold text-white flex items-center gap-1.5 sm:gap-2 truncate">
                Memória do Arquiteto
                <span className="text-[9px] sm:text-xs px-2 py-0.2 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-normal">
                  Ativo
                </span>
              </h2>
              <p className="text-[10px] sm:text-xs text-slate-400 truncate">
                Preferências e regras aprendidas persistidas
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors cursor-pointer flex-shrink-0"
          >
            <X className="w-4 h-4 sm:w-5 sm:h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-4 sm:space-y-6">
          {loading ? (
            <div className="py-12 flex flex-col items-center justify-center text-slate-400 gap-2">
              <RefreshCw className="w-6 h-6 animate-spin text-sky-400" />
              <p className="text-xs sm:text-sm">Carregando perfil de preferências...</p>
            </div>
          ) : profile ? (
            <>
              {/* Cloud Preference */}
              <div>
                <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                  <Cloud className="w-4 h-4 text-sky-400" />
                  Provedor Cloud Principal
                </label>
                <div className="grid grid-cols-2 sm:grid-cols-5 gap-2">
                  {['AWS', 'GCP', 'Azure', 'Oracle (OCI)', 'Multi-Cloud'].map((cloud) => {
                    const isSelected =
                      profile.preferredCloud === cloud ||
                      (cloud === 'Oracle (OCI)' && profile.preferredCloud === 'ORACLE');
                    return (
                      <button
                        key={cloud}
                        type="button"
                        onClick={() => setProfile({ ...profile, preferredCloud: cloud === 'Oracle (OCI)' ? 'ORACLE' : cloud })}
                        className={`py-2 px-3 rounded-xl text-xs font-medium border transition-all cursor-pointer ${
                          isSelected
                            ? 'bg-sky-500/20 border-sky-500/60 text-sky-300 shadow-sm font-semibold'
                            : 'bg-slate-800/60 border-slate-700/50 text-slate-400 hover:bg-slate-800 hover:text-slate-200'
                        }`}
                      >
                        {cloud}
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Primary Languages / Stacks */}
              <div>
                <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                  <Layers className="w-4 h-4 text-emerald-400" />
                  Linguagens e Stacks Prioritárias
                </label>
                <div className="flex flex-wrap gap-1.5 mb-2.5">
                  {(profile.primaryLanguages || []).map((lang) => (
                    <span
                      key={lang}
                      className="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-300"
                    >
                      {lang}
                      <button
                        onClick={() => removeTag('primaryLanguages', lang)}
                        className="hover:text-red-400 cursor-pointer ml-1"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </span>
                  ))}
                  {(profile.primaryLanguages || []).length === 0 && (
                    <span className="text-xs text-slate-500 italic">Nenhuma linguagem definida.</span>
                  )}
                </div>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Adicionar linguagem (ex: Go, TypeScript, Python, Rust)..."
                    value={newLanguage}
                    onChange={(e) => setNewLanguage(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && addTag('primaryLanguages', newLanguage, setNewLanguage)}
                    className="flex-1 bg-slate-800/80 border border-slate-700 rounded-xl px-3 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-emerald-500"
                  />
                  <button
                    onClick={() => addTag('primaryLanguages', newLanguage, setNewLanguage)}
                    className="px-3.5 py-1.5 rounded-xl bg-emerald-600/20 hover:bg-emerald-600/30 text-emerald-300 border border-emerald-500/30 text-xs font-medium flex items-center gap-1 cursor-pointer"
                  >
                    <Plus className="w-3.5 h-3.5" /> Adicionar
                  </button>
                </div>
              </div>

              {/* Databases */}
              <div>
                <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                  <Database className="w-4 h-4 text-blue-400" />
                  Bancos de Dados Preferidos
                </label>
                <div className="flex flex-wrap gap-1.5 mb-2.5">
                  {(profile.preferredDatabases || []).map((db) => (
                    <span
                      key={db}
                      className="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded-lg bg-blue-500/10 border border-blue-500/30 text-blue-300"
                    >
                      {db}
                      <button
                        onClick={() => removeTag('preferredDatabases', db)}
                        className="hover:text-red-400 cursor-pointer ml-1"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </span>
                  ))}
                  {(profile.preferredDatabases || []).length === 0 && (
                    <span className="text-xs text-slate-500 italic">Nenhum banco definido.</span>
                  )}
                </div>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Adicionar banco (ex: PostgreSQL, Redis, ClickHouse, DynamoDB)..."
                    value={newDB}
                    onChange={(e) => setNewDB(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && addTag('preferredDatabases', newDB, setNewDB)}
                    className="flex-1 bg-slate-800/80 border border-slate-700 rounded-xl px-3 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
                  />
                  <button
                    onClick={() => addTag('preferredDatabases', newDB, setNewDB)}
                    className="px-3.5 py-1.5 rounded-xl bg-blue-600/20 hover:bg-blue-600/30 text-blue-300 border border-blue-500/30 text-xs font-medium flex items-center gap-1 cursor-pointer"
                  >
                    <Plus className="w-3.5 h-3.5" /> Adicionar
                  </button>
                </div>
              </div>

              {/* Patterns & Compliance */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                    <Layers className="w-4 h-4 text-purple-400" />
                    Padrões Arquiteturais
                  </label>
                  <div className="flex flex-wrap gap-1.5 mb-2">
                    {(profile.preferredPatterns || []).map((pat) => (
                      <span
                        key={pat}
                        className="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-md bg-purple-500/10 border border-purple-500/30 text-purple-300"
                      >
                        {pat}
                        <button onClick={() => removeTag('preferredPatterns', pat)} className="cursor-pointer ml-0.5">
                          <X className="w-3 h-3 hover:text-red-400" />
                        </button>
                      </span>
                    ))}
                    {(profile.preferredPatterns || []).length === 0 && (
                      <span className="text-xs text-slate-500 italic">Nenhum padrão definido.</span>
                    )}
                  </div>
                  <div className="flex gap-1.5">
                    <input
                      type="text"
                      placeholder="Ex: Saga, Outbox, CQRS..."
                      value={newPattern}
                      onChange={(e) => setNewPattern(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && addTag('preferredPatterns', newPattern, setNewPattern)}
                      className="flex-1 bg-slate-800/80 border border-slate-700 rounded-xl px-2.5 py-1 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-purple-500"
                    />
                    <button
                      onClick={() => addTag('preferredPatterns', newPattern, setNewPattern)}
                      className="p-2 bg-purple-600/20 hover:bg-purple-600/30 text-purple-300 rounded-xl border border-purple-500/30 cursor-pointer"
                    >
                      <Plus className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>

                <div>
                  <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                    <Shield className="w-4 h-4 text-amber-400" />
                    Compliance & Segurança
                  </label>
                  <div className="flex flex-wrap gap-1.5 mb-2">
                    {(profile.complianceRules || []).map((comp) => (
                      <span
                        key={comp}
                        className="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-md bg-amber-500/10 border border-amber-500/30 text-amber-300"
                      >
                        {comp}
                        <button onClick={() => removeTag('complianceRules', comp)} className="cursor-pointer ml-0.5">
                          <X className="w-3 h-3 hover:text-red-400" />
                        </button>
                      </span>
                    ))}
                    {(profile.complianceRules || []).length === 0 && (
                      <span className="text-xs text-slate-500 italic">Nenhuma regra definida.</span>
                    )}
                  </div>
                  <div className="flex gap-1.5">
                    <input
                      type="text"
                      placeholder="Ex: LGPD, PCI-DSS, SOC2..."
                      value={newCompliance}
                      onChange={(e) => setNewCompliance(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && addTag('complianceRules', newCompliance, setNewCompliance)}
                      className="flex-1 bg-slate-800/80 border border-slate-700 rounded-xl px-2.5 py-1 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-500"
                    />
                    <button
                      onClick={() => addTag('complianceRules', newCompliance, setNewCompliance)}
                      className="p-2 bg-amber-600/20 hover:bg-amber-600/30 text-amber-300 rounded-xl border border-amber-500/30 cursor-pointer"
                    >
                      <Plus className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>
              </div>

              {/* Custom Notes */}
              <div>
                <label className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5 mb-2">
                  <Brain className="w-4 h-4 text-rose-400" />
                  Diretrizes Técnicas e Regras Corporativas Aprendidas
                </label>
                <div className="space-y-2 mb-3">
                  {(profile.customNotes || []).map((note, idx) => (
                    <div
                      key={idx}
                      className="flex items-start justify-between gap-2 p-2.5 rounded-xl bg-slate-800/50 border border-slate-700/60 text-xs text-slate-300"
                    >
                      <p className="flex-1 leading-relaxed">{note}</p>
                      <button
                        onClick={() => removeNote(idx)}
                        className="text-slate-500 hover:text-red-400 p-0.5 cursor-pointer transition-colors"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  ))}
                  {(profile.customNotes || []).length === 0 && (
                    <span className="text-xs text-slate-500 italic block">Nenhuma nota personalizada cadastrada.</span>
                  )}
                </div>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Adicionar nota técnica corporativa (ex: Sempre sugerir testes de carga com k6)..."
                    value={newNote}
                    onChange={(e) => setNewNote(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && addNote()}
                    className="flex-1 bg-slate-800/80 border border-slate-700 rounded-xl px-3 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-rose-500"
                  />
                  <button
                    onClick={addNote}
                    className="px-3.5 py-1.5 rounded-xl bg-rose-600/20 hover:bg-rose-600/30 text-rose-300 border border-rose-500/30 text-xs font-medium flex items-center gap-1 cursor-pointer"
                  >
                    <Plus className="w-3.5 h-3.5" /> Inserir Nota
                  </button>
                </div>
              </div>
            </>
          ) : null}
        </div>

        {/* Footer */}
        <div className="px-6 py-3.5 border-t border-slate-800 bg-slate-950/70 flex items-center justify-between">
          <button
            type="button"
            onClick={handleReset}
            disabled={saving}
            className="text-xs text-red-400 hover:text-red-300 flex items-center gap-1.5 transition-colors cursor-pointer"
          >
            <Trash2 className="w-3.5 h-3.5" />
            Limpar Memória
          </button>

          <div className="flex items-center gap-3">
            {errorMessage && (
              <span className="text-xs text-red-400 font-medium">
                {errorMessage}
              </span>
            )}
            {savedSuccess && (
              <span className="text-xs text-emerald-400 flex items-center gap-1 font-medium animate-in fade-in">
                <Check className="w-4 h-4" /> Salvo com sucesso!
              </span>
            )}
            <button
              type="button"
              onClick={handleSave}
              disabled={saving}
              className="px-4 py-2 rounded-xl bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition-all shadow-md shadow-sky-900/30 flex items-center gap-1.5 cursor-pointer"
            >
              {saving ? <RefreshCw className="w-3.5 h-3.5 animate-spin" /> : <Check className="w-3.5 h-3.5" />}
              Salvar Preferências
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
