import { create } from 'zustand';
import type { Config, BotStatus, ChatMessage, ToastMessage, NavPage, LogEntry, QQSession, QQMessage } from '../types';

interface AppState {
  currentPage: NavPage;
  setCurrentPage: (page: NavPage) => void;

  config: Config | null;
  setConfig: (config: Config) => void;

  botStatus: BotStatus;
  setBotStatus: (status: BotStatus) => void;

  chatMessages: ChatMessage[];
  setChatMessages: (messages: ChatMessage[]) => void;
  addChatMessage: (message: ChatMessage) => void;
  clearChatMessages: () => void;

  sessions: QQSession[];
  setSessions: (sessions: QQSession[]) => void;
  activeSessionId: string | null;
  setActiveSessionId: (id: string | null) => void;

  qqMessages: Record<string, QQMessage[]>;
  addQQMessage: (sessionId: string, message: QQMessage) => void;

  logs: LogEntry[];
  addLog: (log: LogEntry) => void;

  toasts: ToastMessage[];
  addToast: (toast: Omit<ToastMessage, 'id'>) => void;
  removeToast: (id: string) => void;
}

const generateId = () => Math.random().toString(36).slice(2, 9);

export const useAppStore = create<AppState>((set) => ({
  currentPage: 'dashboard',
  setCurrentPage: (page) => set({ currentPage: page }),

  config: null,
  setConfig: (config) => set({ config }),

  botStatus: { running: false, connectedAt: '', oneBotURL: '' },
  setBotStatus: (status) => set({ botStatus: status }),

  chatMessages: [],
  setChatMessages: (messages) => set({ chatMessages: messages }),
  addChatMessage: (message) =>
    set((state) => ({ chatMessages: [...state.chatMessages, message] })),
  clearChatMessages: () => set({ chatMessages: [] }),

  sessions: [],
  setSessions: (sessions) => set({ sessions }),
  activeSessionId: null,
  setActiveSessionId: (id) => set({ activeSessionId: id }),

  qqMessages: {},
  addQQMessage: (sessionId, message) =>
    set((state) => ({
      qqMessages: {
        ...state.qqMessages,
        [sessionId]: [...(state.qqMessages[sessionId] || []), message],
      },
    })),

  logs: [],
  addLog: (log) =>
    set((state) => ({ logs: [...state.logs.slice(-499), log] })),

  toasts: [],
  addToast: (toast) => {
    const id = generateId();
    set((state) => ({ toasts: [...state.toasts, { ...toast, id }] }));
    setTimeout(() => {
      set((state) => ({
        toasts: state.toasts.filter((t) => t.id !== id),
      }));
    }, 3000);
  },
  removeToast: (id) =>
    set((state) => ({ toasts: state.toasts.filter((t) => t.id !== id) })),
}));
