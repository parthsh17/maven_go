import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Order, OrderEvent } from '../types';
import { api } from '../services/api';
import { Navbar } from '../components/Navbar';
import { Badge } from '../components/Badge';
import { StateTimeline } from '../components/StateTimeline';

export function OrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [events, setEvents] = useState<OrderEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!id) return;
    const fetchData = async () => {
      try {
        const [o, e] = await Promise.all([api.getOrder(id), api.getOrderEvents(id)]);
        setOrder(o);
        setEvents(e ?? []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load order');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
    const interval = setInterval(fetchData, 3000);
    return () => clearInterval(interval);
  }, [id]);

  if (loading) return <><Navbar /><main className="page"><div className="skeleton-detail" /></main></>;
  if (error || !order) return <><Navbar /><main className="page"><div className="form-error">{error || 'Order not found'}</div></main></>;

  return (
    <>
      <Navbar />
      <main className="page">
        <div className="page-header">
          <div>
            <Link to="/orders" className="back-link">← Back to Orders</Link>
            <h1 className="page-title">Order Detail</h1>
            <p className="page-sub monospace">{order.id}</p>
          </div>
          <Badge state={order.state} />
        </div>

        <div className="detail-grid">
          {}
          <div className="glass-card">
            <h2 className="card-title">Order Metadata</h2>
            <dl className="meta-list">
              <div className="meta-row"><dt>Asset</dt><dd className="asset-big">{order.asset}</dd></div>
              <div className="meta-row"><dt>Quantity</dt><dd>{order.quantity.toLocaleString()}</dd></div>
              <div className="meta-row"><dt>Type</dt><dd><span className="type-chip">{order.order_type}</span></dd></div>
              <div className="meta-row"><dt>State</dt><dd><Badge state={order.state} /></dd></div>
              <div className="meta-row"><dt>Retry Count</dt><dd className={order.retry_count > 0 ? 'retry-count' : ''}>{order.retry_count} / 3</dd></div>
              <div className="meta-row"><dt>Created</dt><dd>{new Date(order.created_at).toLocaleString()}</dd></div>
              <div className="meta-row"><dt>Updated</dt><dd>{new Date(order.updated_at).toLocaleString()}</dd></div>
            </dl>
          </div>

          {}
          <div className="glass-card">
            <h2 className="card-title">State Timeline</h2>
            <StateTimeline events={events} />
          </div>
        </div>
      </main>
    </>
  );
}
