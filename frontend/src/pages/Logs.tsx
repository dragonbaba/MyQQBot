import { useAppStore } from '../store/useAppStore';
import { Card } from '../components/common/Card';

const levelColors = {
  INFO: 'text-text-secondary',
  WARN: 'text-status-warning',
  ERROR: 'text-status-error',
};

export function Logs() {
  const logs = useAppStore((s) => s.logs);

  return (
    <div className="h-full flex flex-col p-6">
      <Card className="flex-1 flex flex-col overflow-hidden">
        <div className="flex items-center gap-2 mb-3">
          {(['ALL', 'INFO', 'WARN', 'ERROR'] as const).map((filter) => (
            <button
              key={filter}
              className="px-3 py-1 text-xs font-medium rounded-md bg-surface-elevated text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
            >
              {filter}
            </button>
          ))}
        </div>
        <div className="flex-1 overflow-auto font-mono text-xs space-y-1">
          {logs.length === 0 ? (
            <p className="text-text-secondary">暂无日志</p>
          ) : (
            logs.map((log, idx) => (
              <div key={idx} className="flex gap-3 py-1 border-b border-border-DEFAULT/50">
                <span className="text-text-muted shrink-0">{log.time}</span>
                <span className={`shrink-0 font-bold ${levelColors[log.level]}`}>{log.level}</span>
                <span className="text-status-info shrink-0">[{log.source}]</span>
                <span className="text-text-secondary">{log.message}</span>
              </div>
            ))
          )}
        </div>
      </Card>
    </div>
  );
}
