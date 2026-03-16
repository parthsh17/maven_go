interface StatCardProps {
  label: string;
  value: number | string;
  color?: string;
  icon?: string;
}

export function StatCard({ label, value, color = '#6C63FF', icon }: StatCardProps) {
  return (
    <div className="stat-card" style={{ '--accent': color } as React.CSSProperties}>
      <div className="stat-icon">{icon}</div>
      <div className="stat-value">{value}</div>
      <div className="stat-label">{label}</div>
      <div className="stat-glow" />
    </div>
  );
}
