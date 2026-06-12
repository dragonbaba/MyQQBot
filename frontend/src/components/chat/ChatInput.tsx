import { Send } from 'lucide-react';
import { useState, KeyboardEvent } from 'react';

interface ChatInputProps {
  placeholder?: string;
  loading?: boolean;
  onSend: (text: string) => void;
}

export function ChatInput({ placeholder = '输入消息...', loading = false, onSend }: ChatInputProps) {
  const [text, setText] = useState('');

  const handleSend = () => {
    const trimmed = text.trim();
    if (!trimmed || loading) return;
    onSend(trimmed);
    setText('');
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex gap-2">
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={loading}
        placeholder={placeholder}
        rows={2}
        className="flex-1 bg-surface-card border border-border-subtle rounded-xl px-4 py-2 text-sm text-text-primary placeholder:text-text-muted resize-none focus:outline-none focus:border-brand-cta/50 disabled:opacity-60"
      />
      <button
        onClick={handleSend}
        disabled={loading}
        title="Enter 发送，Shift+Enter 换行"
        className="px-4 bg-brand-cta hover:bg-brand-cta/80 disabled:bg-brand-cta/40 text-surface-base rounded-xl transition-colors cursor-pointer flex items-center justify-center"
      >
        <Send className="w-5 h-5" />
      </button>
    </div>
  );
}
