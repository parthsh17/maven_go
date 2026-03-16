import type { OrderState } from '../types';

const stateColors: Record<OrderState, string> = {
  CREATED:   '#94a3b8',
  VALIDATED: '#60a5fa',
  QUEUED:    '#a78bfa',
  EXECUTING: '#fbbf24',
  COMPLETED: '#34d399',
  FAILED:    '#f87171',
  RETRYING:  '#fb923c',
};

interface BadgeProps {
  state: string;
}

export function Badge({ state }: BadgeProps) {
  const color = stateColors[state as OrderState] ?? '#94a3b8';
  return (
    <span
      className="badge"
      style={{ '--badge-color': color } as React.CSSProperties}
    >
      {state}
    </span>
  );
}
