import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { ChatMessage } from '../../types';

interface MessagesProps {
  messages: ChatMessage[];
  loading?: boolean;
}

export function Messages({ messages, loading = false }: MessagesProps) {
  return (
    <div className="flex-1 overflow-auto space-y-3 pr-2">
      {messages.map((msg, idx) => {
        const isUser = msg.role === 'user';
        return (
          <div
            key={idx}
            className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[80%] px-4 py-3 rounded-2xl text-sm leading-relaxed ${
                isUser
                  ? 'bg-brand-cta/20 rounded-tr-sm text-text-primary'
                  : 'bg-surface-elevated rounded-tl-sm text-text-primary'
              }`}
            >
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                  code({ children, className }) {
                    const isInline = !className;
                    return (
                      <code
                        className={`${
                          isInline
                            ? 'bg-surface-base px-1 py-0.5 rounded text-brand-cta'
                            : 'block bg-surface-base p-3 rounded-lg overflow-x-auto text-xs font-mono text-text-primary my-2'
                        }`}
                      >
                        {children}
                      </code>
                    );
                  },
                  a: ({ href, children }) => (
                    <a
                      href={href}
                      className="text-brand-cta hover:underline"
                      target="_blank"
                      rel="noreferrer"
                    >
                      {children}
                    </a>
                  ),
                  ul: ({ children }) => <ul className="list-disc pl-5 my-2">{children}</ul>,
                  ol: ({ children }) => <ol className="list-decimal pl-5 my-2">{children}</ol>,
                  p: ({ children }) => <p className="mb-1 last:mb-0">{children}</p>,
                }}
              >
                {msg.content}
              </ReactMarkdown>
            </div>
          </div>
        );
      })}
      {loading && (
        <div className="flex justify-start">
          <div className="bg-surface-elevated rounded-2xl rounded-tl-sm px-4 py-3 space-y-2">
            <div className="animate-pulse bg-surface-base rounded h-3 w-48" />
            <div className="animate-pulse bg-surface-base rounded h-3 w-32" />
          </div>
        </div>
      )}
    </div>
  );
}
