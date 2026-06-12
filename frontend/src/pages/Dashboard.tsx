import { useEffect, useState } from 'react';
import { Card } from '../components/common/Card';
import { StatusBadge } from '../components/common/StatusBadge';
import { Skeleton } from '../components/common/Skeleton';
import { useAppStore } from '../store/useAppStore';
import { Bot, MessageSquare, Users, Zap } from 'lucide-react';
import * as BotService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/botservice.js';

export function Dashboard() {
  const botStatus = useAppStore((s) => s.botStatus);
  const setBotStatus = useAppStore((s) => s.setBotStatus);
  const sessions = useAppStore((s) => s.sessions);
  const logs = useAppStore((s) => s.logs);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    BotService.GetBotStatus()
      .then((status) => {
        setBotStatus(status);
      })
      .finally(() => setLoading(false));
  }, [setBotStatus]);

  const recentMessages = sessions.slice(0, 20);

  return (
    <div className="h-full overflow-auto p-6 space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <Card className="flex items-center gap-4">
          <div className="p-3 rounded-xl bg-brand-cta/10 text-brand-cta">
            <Bot className="w-6 h-6" />
          </div>
          <div>
            <p className="text-text-muted text-xs font-mono uppercase tracking-wider">Bot 状态</p>
            <div className="mt-1">
              {loading ? (
                <Skeleton width={80} height={16} />
              ) : (
                <StatusBadge
                  color={botStatus.running ? 'online' : 'error'}
                  label={botStatus.running ? '运行中' : '已停止'}
                  pulse={botStatus.running}
                />
              )}
            </div>
          </div>
        </Card>
        <Card className="flex items-center gap-4">
          <div className="p-3 rounded-xl bg-status-info/10 text-status-info">
            <MessageSquare className="w-6 h-6" />
          </div>
          <div>
            <p className="text-text-muted text-xs font-mono uppercase tracking-wider">今日消息</p>
            <p className="text-2xl font-mono text-text-primary mt-1">
              {logs.filter((l) => l.source === 'QQ').length}
            </p>
          </div>
        </Card>
        <Card className="flex items-center gap-4">
          <div className="p-3 rounded-xl bg-status-warning/10 text-status-warning">
            <Users className="w-6 h-6" />
          </div>
          <div>
            <p className="text-text-muted text-xs font-mono uppercase tracking-wider">活跃会话</p>
            <p className="text-2xl font-mono text-text-primary mt-1">{sessions.length}</p>
          </div>
        </Card>
        <Card className="flex items-center gap-4">
          <div className="p-3 rounded-xl bg-brand-cta/10 text-brand-cta">
            <Zap className="w-6 h-6" />
          </div>
          <div>
            <p className="text-text-muted text-xs font-mono uppercase tracking-wider">LLM 延迟</p>
            <p className="text-2xl font-mono text-text-primary mt-1">-- ms</p>
          </div>
        </Card>
      </div>

      <Card>
        <h3 className="text-sm font-semibold font-mono text-text-primary mb-4">最近消息</h3>
        <div className="space-y-2">
          {recentMessages.length === 0 ? (
            <p className="text-text-secondary text-sm">暂无消息</p>
          ) : (
            recentMessages.map((session) => (
              <div
                key={session.id}
                className="flex items-center gap-3 p-3 rounded-lg bg-surface-base/50 hover:bg-surface-elevated transition-colors cursor-pointer"
              >
                <div className="w-8 h-8 rounded-full bg-surface-elevated flex items-center justify-center text-xs font-bold text-brand-cta">
                  {session.avatarText}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-text-primary truncate">{session.name}</p>
                  <p className="text-xs text-text-secondary truncate">{session.lastMessage}</p>
                </div>
                <span className="text-xs text-text-muted font-mono">{session.lastTime}</span>
              </div>
            ))
          )}
        </div>
      </Card>
    </div>
  );
}
