import type { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
}

export function Card({ children, className = '' }: CardProps) {
  return (
    <div className={`bg-surface-card border border-border-subtle rounded-xl p-5 ${className}`}>
      {children}
    </div>
  );
}
