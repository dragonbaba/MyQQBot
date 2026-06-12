import { useState } from 'react';
import { Sidebar } from './Sidebar';
import { TopBar } from './TopBar';
import { Dashboard } from '../../pages/Dashboard';
import { Chat } from '../../pages/Chat';
import { Settings } from '../../pages/Settings';
import { Logs } from '../../pages/Logs';
import { ToastContainer } from '../common/Toast';
import { useAppStore } from '../../store/useAppStore';

const pages = {
  dashboard: Dashboard,
  chat: Chat,
  settings: Settings,
  logs: Logs,
};

export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const currentPage = useAppStore((s) => s.currentPage);
  const PageComponent = pages[currentPage] || Dashboard;

  return (
    <div className="h-screen w-screen grid grid-cols-[auto_1fr] grid-rows-[48px_1fr] bg-surface-base overflow-hidden">
      <div className="row-span-2">
        <Sidebar collapsed={collapsed} onToggle={() => setCollapsed((c) => !c)} />
      </div>
      <TopBar />
      <main className="overflow-hidden">
        <PageComponent />
      </main>
      <ToastContainer />
    </div>
  );
}
