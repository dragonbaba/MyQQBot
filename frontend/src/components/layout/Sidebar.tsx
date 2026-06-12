import { LayoutGrid, MessageSquare, Settings, List, ChevronLeft, ChevronRight } from 'lucide-react';
import { useAppStore } from '../../store/useAppStore';
import { StatusBadge } from '../common/StatusBadge';
import type { NavPage } from '../../types';

interface NavItem {
  id: NavPage;
  label: string;
  icon: typeof LayoutGrid;
}

const navItems: NavItem[] = [
  { id: 'dashboard', label: 'Dashboard', icon: LayoutGrid },
  { id: 'chat', label: 'Chat', icon: MessageSquare },
  { id: 'settings', label: 'Settings', icon: Settings },
  { id: 'logs', label: 'Logs', icon: List },
];

interface SidebarProps {
  collapsed: boolean;
  onToggle: () => void;
}

export function Sidebar({ collapsed, onToggle }: SidebarProps) {
  const currentPage = useAppStore((s) => s.currentPage);
  const setCurrentPage = useAppStore((s) => s.setCurrentPage);
  const botStatus = useAppStore((s) => s.botStatus);

  return (
    <aside
      className={`h-full bg-[#0A0F1A] border-r border-border-DEFAULT flex flex-col transition-all duration-200 ${
        collapsed ? 'w-16' : 'w-60'
      }`}
    >
      <div className="h-12 flex items-center px-4 border-b border-border-DEFAULT">
        {!collapsed && (
          <span className="text-lg font-bold font-mono text-brand-cta">MyQQBot</span>
        )}
      </div>

      <nav className="flex-1 py-3 space-y-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = currentPage === item.id;
          return (
            <button
              key={item.id}
              onClick={() => setCurrentPage(item.id)}
              className={`w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium transition-all duration-200 cursor-pointer ${
                active
                  ? 'bg-brand-cta/15 text-brand-cta border-l-2 border-brand-cta'
                  : 'text-text-secondary hover:text-text-primary hover:bg-surface-card border-l-2 border-transparent'
              }`}
              title={collapsed ? item.label : undefined}
            >
              <Icon className="w-5 h-5 shrink-0" />
              {!collapsed && <span>{item.label}</span>}
            </button>
          );
        })}
      </nav>

      <div className="p-3 border-t border-border-DEFAULT space-y-3">
        <div className={`flex items-center ${collapsed ? 'justify-center' : 'gap-2 px-2'}`}>
          <StatusBadge
            color={botStatus.running ? 'online' : 'error'}
            label={collapsed ? '' : botStatus.running ? '已连接' : '未连接'}
            pulse={botStatus.running}
          />
        </div>
        <button
          onClick={onToggle}
          className="w-full flex items-center justify-center p-2 text-text-secondary hover:text-text-primary hover:bg-surface-card rounded-lg transition-colors cursor-pointer"
          title={collapsed ? '展开' : '收起'}
        >
          {collapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
        </button>
      </div>
    </aside>
  );
}
