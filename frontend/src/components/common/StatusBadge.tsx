interface StatusBadgeProps {
  color: 'online' | 'warning' | 'error' | 'info';
  label: string;
  pulse?: boolean;
}

const colorMap = {
  online: 'bg-status-online',
  warning: 'bg-status-warning',
  error: 'bg-status-error',
  info: 'bg-status-info',
};

export function StatusBadge({ color, label, pulse = false }: StatusBadgeProps) {
  return (
    <span className="inline-flex items-center gap-1.5 text-sm text-text-secondary">
      <span
        className={`w-2 h-2 rounded-full ${colorMap[color]} ${pulse ? 'animate-pulse' : ''}`}
      />
      {label}
    </span>
  );
}
