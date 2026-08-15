'use client';

import { useState, useRef, useCallback, useEffect } from 'react';
import { Message, ToolEvent, StreamEvent } from '@/types';
import { getStreamEndpoint, fetchSessionMessages } from '@/lib/api';

export function useChatStream() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState<boolean>(false);
  const [currentTool, setCurrentTool] = useState<{ toolName: string; status: string } | null>(null);
  const [sessionID, setSessionID] = useState<string>('');
  const [error, setError] = useState<string | null>(null);

  const abortControllerRef = useRef<AbortController | null>(null);

  // Initialize or restore session ID
  useEffect(() => {
    const saved = localStorage.getItem('archmind_session_id');
    if (saved) {
      setSessionID(saved);
      fetchSessionMessages(saved).then((msgs) => {
        if (msgs && msgs.length > 0) {
          setMessages(msgs);
        }
      });
    } else {
      const newID = `sess-${Date.now()}`;
      setSessionID(newID);
      localStorage.setItem('archmind_session_id', newID);
    }
  }, []);

  const switchSession = useCallback(async (newSessionID: string) => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    setIsStreaming(false);
    setCurrentTool(null);
    setError(null);
    setSessionID(newSessionID);
    localStorage.setItem('archmind_session_id', newSessionID);

    const msgs = await fetchSessionMessages(newSessionID);
    setMessages(msgs || []);
  }, []);

  const clearMessages = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    setMessages([]);
    setError(null);
    setCurrentTool(null);
    setIsStreaming(false);

    const newID = `sess-${Date.now()}`;
    setSessionID(newID);
    localStorage.setItem('archmind_session_id', newID);
  }, []);

  const stopStream = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      setIsStreaming(false);
      setCurrentTool(null);
    }
  }, []);

  const sendMessage = useCallback(
    async (content: string, imageBase64?: string, imageMimeType?: string) => {
      if (!content.trim() && !imageBase64) return;

      setError(null);
      setIsStreaming(true);
      setCurrentTool(null);

      const userMsg: Message = {
        id: `usr-${Date.now()}`,
        role: 'user',
        content,
        imageBase64,
        imageMimeType,
        createdAt: new Date().toISOString(),
      };

      const assistantMsgId = `asst-${Date.now()}`;
      const assistantMsg: Message = {
        id: assistantMsgId,
        role: 'model',
        content: '',
        toolEvents: [],
        createdAt: new Date().toISOString(),
      };

      setMessages((prev) => [...prev, userMsg, assistantMsg]);

      abortControllerRef.current = new AbortController();

      try {
        const streamUrl = getStreamEndpoint();
        const response = await fetch(streamUrl, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            sessionId: sessionID,
            message: content,
            imageBase64,
            imageMimeType,
          }),
          signal: abortControllerRef.current.signal,
        });

        if (!response.ok) {
          throw new Error(`HTTP Error: ${response.status} ${response.statusText}`);
        }

        const reader = response.body?.getReader();
        if (!reader) throw new Error('Response body is not readable');

        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n');
          buffer = lines.pop() || '';

          for (const line of lines) {
            const trimmed = line.trim();
            if (!trimmed || trimmed.startsWith(':') || trimmed.startsWith('event:')) continue;

            if (trimmed.startsWith('data:')) {
              const jsonStr = trimmed.replace('data:', '').trim();
              if (!jsonStr || jsonStr === '[DONE]') continue;

              try {
                const event: StreamEvent = JSON.parse(jsonStr);

                if (event.type === 'ping' && event.data?.sessionId) {
                  setSessionID(event.data.sessionId);
                  localStorage.setItem('archmind_session_id', event.data.sessionId);
                }

                if (event.type === 'tool_start') {
                  setCurrentTool({
                    toolName: event.data?.toolName || 'tool',
                    status: event.content || 'Executando ferramenta de arquitetura...',
                  });
                }

                if (event.type === 'tool_result') {
                  const toolEvt: ToolEvent = event.data;
                  setCurrentTool(null);
                  setMessages((prev) =>
                    prev.map((msg) => {
                      if (msg.id === assistantMsgId) {
                        const existing = msg.toolEvents || [];
                        return { ...msg, toolEvents: [...existing, toolEvt] };
                      }
                      return msg;
                    })
                  );
                }

                if (event.type === 'token' && event.content) {
                  setMessages((prev) =>
                    prev.map((msg) => {
                      if (msg.id === assistantMsgId) {
                        return { ...msg, content: msg.content + event.content };
                      }
                      return msg;
                    })
                  );
                }

                if (event.type === 'error') {
                  setError(event.content || 'Erro no processamento da IA.');
                }

                if (event.type === 'done') {
                  setIsStreaming(false);
                  setCurrentTool(null);
                }
              } catch (e) {
                // Ignore parse errors
              }
            }
          }
        }
      } catch (err: any) {
        if (err.name !== 'AbortError') {
          setError(err.message || 'Falha ao conectar com o backend.');
        }
      } finally {
        setIsStreaming(false);
        setCurrentTool(null);
      }
    },
    [sessionID]
  );

  return {
    messages,
    isStreaming,
    currentTool,
    sessionID,
    error,
    sendMessage,
    stopStream,
    clearMessages,
    switchSession,
    setMessages,
  };
}
