import { Link } from 'react-router-dom';
import { StatCard } from '../components/StatCard';
import { OrderTable } from '../components/OrderTable';
import { useOrders } from '../hooks/useOrders';
import { useMetrics } from '../hooks/useMetrics';
import { Navbar } from '../components/Navbar';

export function DashboardPage() {
  const { orders, loading: ordersLoading } = useOrders();
  const { metrics, loading: metricsLoading } = useMetrics();

  const recent = [...orders]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5);

  return (
    <>
      <Navbar />
      <main className="page">
        <div className="page-header">
          <div>
            <h1 className="page-title">Dashboard</h1>
            <p className="page-sub">System overview and recent activity</p>
          </div>
          <Link to="/submit" className="btn btn-primary">+ New Order</Link>
        </div>

        {metricsLoading ? (
          <div className="skeleton-grid">
            {[1, 2, 3, 4].map((i) => <div key={i} className="skeleton-card" />)}
          </div>
        ) : (
          <div className="stats-grid">
            <StatCard label="Total Orders" value={metrics?.total_orders ?? 0} icon="📋" color="#6C63FF" />
            <StatCard label="Processing" value={metrics?.processing_orders ?? 0} icon="⚡" color="#fbbf24" />
            <StatCard label="Completed" value={metrics?.completed_orders ?? 0} icon="✅" color="#34d399" />
            <StatCard label="Failed" value={metrics?.failed_orders ?? 0} icon="❌" color="#f87171" />
          </div>
        )}

        <section className="section">
          <div className="section-header">
            <h2 className="section-title">Recent Orders</h2>
            <Link to="/orders" className="btn btn-ghost">View All →</Link>
          </div>
          {ordersLoading ? (
            <div className="skeleton-table" />
          ) : (
            <OrderTable orders={recent} />
          )}
        </section>
      </main>
    </>
  );
}
