export interface BotConfig {
  onebot_ws_url: string;
  admin_qq: string;
}

export interface LLMConfig {
  base_url: string;
  api_key: string;
  model: string;
  system_prompt: string;
}

export interface SearchConfig {
  provider: string;
  tavily_api_key: string;
}

export interface Config {
  llm: LLMConfig;
  bot: BotConfig;
  search: SearchConfig;
}

export interface BotStatus {
  running: boolean;
  connectedAt: string;
  oneBotURL: string;
}

export interface ChatMessage {
  role: string;
  content: string;
  timestamp: number;
}

export type NavPage = 'dashboard' | 'chat' | 'settings' | 'logs';

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'warning';
  message: string;
}

export interface QQSession {
  id: string;
  type: 'group' | 'private';
  name: string;
  avatarText: string;
  lastMessage: string;
  lastTime: string;
  unread: number;
}

export interface QQMessage {
  id: string;
  role: 'self' | 'other';
  senderName: string;
  content: string;
  timestamp: number;
}

export interface LogEntry {
  time: string;
  level: 'INFO' | 'WARN' | 'ERROR';
  source: string;
  message: string;
}
