import { Navbar } from '../components/Navbar';
import { StatCard } from '../components/StatCard';
import { useMetrics } from '../hooks/useMetrics';
import { useOrders } from '../hooks/useOrders';

export function MetricsPage() {
  const { metrics, loading, error } = useMetrics();
  const { orders } = useOrders();

  
  const stateDist = orders.reduce<Record<string, number>>((acc, o) => {
    acc[o.state] = (acc[o.state] ?? 0) + 1;
    return acc;
  }, {});

  const stateColors: Record<string, string> = {
    CREATED: '#94a3b8', VALIDATED: '#60a5fa', QUEUED: '#a78bfa',
    EXECUTING: '#fbbf24', COMPLETED: '#34d399', FAILED: '#f87171', RETRYING: '#fb923c',
  };

  return (
    <>
      <Navbar />
      <main className="page">
        <div className="page-header">
          <div>
            <h1 className="page-title">System Metrics</h1>
            <p className="page-sub">Live counters — refreshes every 2 seconds</p>
          </div>
          <div className="live-dot">
            <span className="pulse" />
            LIVE
          </div>
        </div>

        {error && <div className="form-error">{error}</div>}

        {loading ? (
          <div className="skeleton-grid">{[1,2,3,4,5].map(i => <div key={i} className="skeleton-card" />)}</div>
        ) : (
          <div className="stats-grid">
            <StatCard label="Total Orders" value={metrics?.total_orders ?? 0} icon="📋" color="#6C63FF" />
            <StatCard label="Processing" value={metrics?.processing_orders ?? 0} icon="⚡" color="#fbbf24" />
            <StatCard label="Completed" value={metrics?.completed_orders ?? 0} icon="✅" color="#34d399" />
            <StatCard label="Failed" value={metrics?.failed_orders ?? 0} icon="❌" color="#f87171" />
            <StatCard label="Workers" value={metrics?.worker_count ?? 0} icon="⚙️" color="#60a5fa" />
            <StatCard 
              label="Success Rate (MA)" 
              value={`${((metrics?.success_rate_ma as number ?? 0) * 100).toFixed(1)}%`} 
              icon="📈" 
              color="#8b5cf6" 
            />
            <StatCard 
              label="Avg. Slippage" 
              value={`${((metrics?.avg_slippage as number ?? 0) * 100).toFixed(4)}%`} 
              icon="📉" 
              color="#f43f5e" 
            />
          </div>
        )}

        <div className="glass-card" style={{ marginTop: '2rem' }}>
          <h2 className="card-title">State Distribution</h2>
          {Object.keys(stateColors).map((state) => {
            const count = stateDist[state] ?? 0;
            const total = orders.length || 1;
            const pct = Math.round((count / total) * 100);
            return (
              <div key={state} className="dist-row">
                <span className="dist-label" style={{ color: stateColors[state] }}>{state}</span>
                <div className="dist-bar-bg">
                  <div
                    className="dist-bar-fill"
                    style={{ width: `${pct}%`, background: stateColors[state] }}
                  />
                </div>
                <span className="dist-count">{count}</span>
              </div>
            );
          })}
          {orders.length === 0 && (
            <p style={{ color: '#64748b', textAlign: 'center', padding: '1rem' }}>No orders yet</p>
          )}
        </div>
      </main>
    </>
  );
}
