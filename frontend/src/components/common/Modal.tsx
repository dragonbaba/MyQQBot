import type { ReactNode } from 'react';
import { X } from 'lucide-react';

interface ModalProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
}

export function Modal({ open, onClose, title, children }: ModalProps) {
  if (!open) return null;
  return (
    <div
      className="fixed inset-0 z-50 bg-surface-overlay backdrop-blur-sm flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div
        className="bg-surface-card rounded-2xl border border-border-subtle max-w-lg w-full p-6 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        {title && (
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold font-mono text-text-primary">{title}</h3>
            <button
              onClick={onClose}
              className="p-1 hover:bg-surface-elevated rounded-lg transition-colors cursor-pointer"
            >
              <X className="w-5 h-5 text-text-secondary" />
            </button>
          </div>
        )}
        {children}
      </div>
    </div>
  );
}
