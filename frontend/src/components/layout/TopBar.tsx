import { useAppStore } from '../../store/useAppStore';

const pageTitles: Record<string, string> = {
  dashboard: 'Dashboard',
  chat: 'Chat',
  settings: 'Settings',
  logs: 'Logs',
};

export function TopBar() {
  const currentPage = useAppStore((s) => s.currentPage);

  return (
    <header className="h-12 px-5 flex items-center justify-between bg-surface-card/60 backdrop-blur-md border-b border-border-DEFAULT">
      <h2 className="text-sm font-semibold font-mono text-text-primary tracking-wide">
        {pageTitles[currentPage] || 'MyQQBot'}
      </h2>
      <div className="flex items-center gap-2">
        {/* Window controls are handled by the OS; placeholder for future custom buttons */}
      </div>
    </header>
  );
}
