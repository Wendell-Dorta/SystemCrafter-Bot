'use client';

import React, { useState } from 'react';
import { useChatStream } from '@/hooks/useChatStream';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { QuickPrompts } from './QuickPrompts';
import { HistoryModal } from './HistoryModal';
import { TemplateModal } from './TemplateModal';
import { ProfileModal } from './ProfileModal';
import { CodeAuditModal } from './CodeAuditModal';
import { ArchitectureTemplate, CodeAuditReport } from '@/types';

export const ChatInterface: React.FC = () => {
  const {
    messages,
    isStreaming,
    currentTool,
    sessionID,
    error,
    sendMessage,
    stopStream,
    clearMessages,
    switchSession,
  } = useChatStream();

  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false);
  const [isHistoryModalOpen, setIsHistoryModalOpen] = useState<boolean>(false);
  const [isTemplateModalOpen, setIsTemplateModalOpen] = useState<boolean>(false);
  const [isProfileModalOpen, setIsProfileModalOpen] = useState<boolean>(false);
  const [isCodeAuditModalOpen, setIsCodeAuditModalOpen] = useState<boolean>(false);

  const handleSelectTemplate = (template: ArchitectureTemplate) => {
    sendMessage(template.promptStarter);
  };

  const handleAuditCompleted = (report: CodeAuditReport) => {
    // Reload active session messages to display the newly injected grounded audit
    switchSession(sessionID);
  };

  return (
    <div className="relative flex h-[100dvh] max-h-[100dvh] w-full flex-col overflow-hidden bg-[#0b0f19] text-slate-100 antialiased">
      {/* Top Navbar (Fixed & Mobile-Optimized) */}
      <Header
        onNewSession={clearMessages}
        onOpenHistory={() => setIsHistoryModalOpen(true)}
        onOpenTemplates={() => setIsTemplateModalOpen(true)}
        onOpenProfile={() => setIsProfileModalOpen(true)}
        onOpenCodeAudit={() => setIsCodeAuditModalOpen(true)}
        onToggleSidebar={() => setIsSidebarOpen((prev) => !prev)}
        isSidebarOpen={isSidebarOpen}
      />

      {/* Sidebar Navigation Drawer */}
      <Sidebar
        isOpen={isSidebarOpen}
        onClose={() => setIsSidebarOpen(false)}
        currentSessionId={sessionID}
        onSelectSession={switchSession}
        onNewSession={clearMessages}
        onSelectTemplate={handleSelectTemplate}
        onOpenHistoryModal={() => setIsHistoryModalOpen(true)}
        onOpenTemplatesModal={() => setIsTemplateModalOpen(true)}
        onOpenProfileModal={() => setIsProfileModalOpen(true)}
        onOpenCodeAuditModal={() => setIsCodeAuditModalOpen(true)}
      />

      {/* Dedicated Chat History Modal with Instant Search & Backdrop Blur */}
      <HistoryModal
        isOpen={isHistoryModalOpen}
        onClose={() => setIsHistoryModalOpen(false)}
        currentSessionId={sessionID}
        onSelectSession={switchSession}
        onNewSession={clearMessages}
      />

      {/* Blueprints Catalog Modal */}
      <TemplateModal
        isOpen={isTemplateModalOpen}
        onClose={() => setIsTemplateModalOpen(false)}
        onSelectTemplate={handleSelectTemplate}
      />

      {/* Long-Term Memory & Learned Profile Modal */}
      <ProfileModal
        isOpen={isProfileModalOpen}
        onClose={() => setIsProfileModalOpen(false)}
      />

      {/* Grounded Code & GitHub Repository Security Audit Modal */}
      <CodeAuditModal
        isOpen={isCodeAuditModalOpen}
        onClose={() => setIsCodeAuditModalOpen(false)}
        currentSessionId={sessionID}
        onAuditCompleted={handleAuditCompleted}
      />

      {/* Main Scrollable Viewport */}
      <main className="flex flex-1 flex-col overflow-hidden w-full relative min-h-0">
        {messages.length === 0 ? (
          <div className="flex-1 overflow-y-auto overscroll-contain w-full min-h-0">
            <QuickPrompts
              onSelectPrompt={(prompt) => sendMessage(prompt)}
              onOpenTemplates={() => setIsTemplateModalOpen(true)}
            />
          </div>
        ) : (
          <MessageList
            messages={messages}
            isStreaming={isStreaming}
            currentTool={currentTool}
            error={error}
          />
        )}

        {/* Bottom Chat Input Bar (Mobile-First Safe Area Aware) */}
        <ChatInput
          onSendMessage={sendMessage}
          onStopStream={stopStream}
          isStreaming={isStreaming}
        />
      </main>
    </div>
  );
};
