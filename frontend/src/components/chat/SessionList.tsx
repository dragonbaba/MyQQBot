import { Users, User } from 'lucide-react';
import type { QQSession } from '../../types';

interface SessionListProps {
  sessions: QQSession[];
  activeId: string | null;
  onSelect: (id: string) => void;
}

export function SessionList({ sessions, activeId, onSelect }: SessionListProps) {
  const groups = sessions.filter((s) => s.type === 'group');
  const privates = sessions.filter((s) => s.type === 'private');

  const renderItem = (session: QQSession) => (
    <button
      key={session.id}
      onClick={() => onSelect(session.id)}
      className={`w-full flex items-center gap-3 p-3 rounded-lg transition-colors text-left cursor-pointer ${
        activeId === session.id
          ? 'bg-surface-elevated border border-brand-cta/30'
          : 'hover:bg-surface-elevated/50 border border-transparent'
      }`}
    >
      <div className="w-9 h-9 rounded-full bg-brand-primary flex items-center justify-center text-sm font-bold text-brand-cta shrink-0">
        {session.avatarText}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium text-text-primary truncate">{session.name}</span>
          <span className="text-xs text-text-muted font-mono">{session.lastTime}</span>
        </div>
        <p className="text-xs text-text-secondary truncate">{session.lastMessage}</p>
      </div>
      {session.unread > 0 && (
        <span className="min-w-[18px] h-[18px] px-1.5 rounded-full bg-status-error text-white text-[10px] font-bold flex items-center justify-center">
          {session.unread}
        </span>
      )}
    </button>
  );

  return (
    <div className="h-full overflow-auto space-y-4 pr-1">
      {groups.length > 0 && (
        <div>
          <div className="flex items-center gap-2 text-xs font-mono text-text-muted uppercase tracking-wider mb-2">
            <Users className="w-3.5 h-3.5" />
            群聊
          </div>
          <div className="space-y-1">{groups.map(renderItem)}</div>
        </div>
      )}
      {privates.length > 0 && (
        <div>
          <div className="flex items-center gap-2 text-xs font-mono text-text-muted uppercase tracking-wider mb-2">
            <User className="w-3.5 h-3.5" />
            私聊
          </div>
          <div className="space-y-1">{privates.map(renderItem)}</div>
        </div>
      )}
      {sessions.length === 0 && (
        <p className="text-sm text-text-secondary text-center py-8">暂无 QQ 会话</p>
      )}
    </div>
  );
}
