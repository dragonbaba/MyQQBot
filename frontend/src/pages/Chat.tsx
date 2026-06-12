import { useState, useEffect, useRef } from 'react';
import { Messages } from '../components/chat/Messages';
import { ChatInput } from '../components/chat/ChatInput';
import { SessionList } from '../components/chat/SessionList';
import { Card } from '../components/common/Card';
import { useAppStore } from '../store/useAppStore';
import * as ChatService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/chatservice.js';
import type { ChatMessage } from '../types';

type ChatMode = 'direct' | 'qq';

export function Chat() {
  const [mode, setMode] = useState<ChatMode>('direct');
  const [loading, setLoading] = useState(false);
  const [pendingImages, setPendingImages] = useState<string[]>([]);
  const chatMessages = useAppStore((s) => s.chatMessages);
  const addChatMessage = useAppStore((s) => s.addChatMessage);
  const sessions = useAppStore((s) => s.sessions);
  const activeSessionId = useAppStore((s) => s.activeSessionId);
  const setActiveSessionId = useAppStore((s) => s.setActiveSessionId);
  const addToast = useAppStore((s) => s.addToast);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    ChatService.GetHistory().then((history) => {
      useAppStore.getState().setChatMessages(
        history.map((m) => ({
          role: m.role,
          content: m.content,
          timestamp: m.timestamp,
        }))
      );
    });
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [chatMessages]);

  const handleDirectSend = async (text: string) => {
    const images = pendingImages.length > 0 ? pendingImages : undefined;
    if (images) {
      addChatMessage({ role: 'user', content: text, timestamp: Date.now(), imageUrls: images });
    } else {
      addChatMessage({ role: 'user', content: text, timestamp: Date.now() });
    }
    setLoading(true);
    try {
      let reply: string;
      if (images && images.length > 0) {
        reply = await ChatService.SendVisionMessage(text, images);
        setPendingImages([]);
      } else {
        reply = await ChatService.SendMessage(text);
      }
      addChatMessage({ role: 'assistant', content: reply, timestamp: Date.now() });
    } catch (err) {
      addToast({ type: 'error', message: `请求失败: ${err}` });
    } finally {
      setLoading(false);
    }
  };

  const activeSession = sessions.find((s) => s.id === activeSessionId);

  return (
    <div className="h-full flex">
      {mode === 'direct' && (
        <aside className="w-72 border-r border-border-DEFAULT bg-surface-card/50 p-3 flex flex-col">
          <h3 className="text-xs font-mono uppercase tracking-wider text-text-muted mb-3">图片附件</h3>
          {pendingImages.length === 0 ? (
            <p className="text-sm text-text-secondary">暂无图片</p>
          ) : (
            <div className="space-y-2">
              {pendingImages.map((url, idx) => (
                <div key={idx} className="relative group">
                  <img src={url} alt="" className="w-full h-24 object-cover rounded-lg border border-border-subtle" />
                  <button
                    onClick={() => setPendingImages((prev) => prev.filter((_, i) => i !== idx))}
                    className="absolute top-1 right-1 p-1 bg-status-error text-white rounded-md opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
          <input
            type="text"
            className="mt-2 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-xs text-text-primary focus:outline-none focus:border-brand-cta/50"
            placeholder="粘贴图片 URL 按回车"
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                const value = (e.target as HTMLInputElement).value.trim();
                if (value) {
                  setPendingImages((prev) => [...prev, value]);
                  (e.target as HTMLInputElement).value = '';
                }
              }
            }}
          />
        </aside>
      )}

      {mode === 'qq' && (
        <aside className="w-72 border-r border-border-DEFAULT bg-surface-card/50 p-3 flex flex-col">
          <SessionList sessions={sessions} activeId={activeSessionId} onSelect={setActiveSessionId} />
        </aside>
      )}

      <div className="flex-1 flex flex-col p-4 min-w-0">
        <div className="flex items-center gap-2 mb-3">
          <button
            onClick={() => setMode('direct')}
            className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors cursor-pointer ${
              mode === 'direct'
                ? 'bg-brand-cta/15 text-brand-cta'
                : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
            }`}
          >
            直接对话
          </button>
          <button
            onClick={() => setMode('qq')}
            className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors cursor-pointer ${
              mode === 'qq'
                ? 'bg-brand-cta/15 text-brand-cta'
                : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
            }`}
          >
            QQ 会话
          </button>
        </div>

        {mode === 'qq' && activeSession && (
          <div className="mb-3 px-3 py-2 bg-surface-card rounded-lg border border-border-subtle">
            <span className="text-sm font-medium text-text-primary">{activeSession.name}</span>
          </div>
        )}

        <Card className="flex-1 flex flex-col overflow-hidden p-4 mb-3">
          <Messages messages={chatMessages as ChatMessage[]} loading={loading} />
          <div ref={messagesEndRef} />
        </Card>
        <ChatInput
          placeholder={mode === 'direct' ? '直接对话 LLM...' : '回复选中会话...'}
          loading={loading}
          onSend={handleDirectSend}
        />
      </div>
    </div>
  );
}
