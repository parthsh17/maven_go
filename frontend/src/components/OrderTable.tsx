import { useNavigate } from 'react-router-dom';
import type { Order } from '../types';
import { Badge } from './Badge';

interface OrderTableProps {
  orders: Order[];
  filterState?: string;
}

export function OrderTable({ orders, filterState }: OrderTableProps) {
  const navigate = useNavigate();

  const filtered = filterState && filterState !== 'ALL'
    ? orders.filter((o) => o.state === filterState)
    : orders;

  if (filtered.length === 0) {
    return (
      <div className="empty-state">
        <span>◌</span>
        <p>No orders found</p>
      </div>
    );
  }

  return (
    <div className="table-wrapper">
      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Asset</th>
            <th>Quantity</th>
            <th>Type</th>
            <th>State</th>
            <th>Retries</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((o) => (
            <tr
              key={o.id}
              className="table-row clickable"
              onClick={() => navigate(`/orders/${o.id}`)}
            >
              <td className="id-cell" title={o.id}>{o.id.slice(0, 8)}…</td>
              <td className="asset-cell">{o.asset}</td>
              <td>{o.quantity.toLocaleString()}</td>
              <td><span className="type-chip">{o.order_type}</span></td>
              <td><Badge state={o.state} /></td>
              <td>{o.retry_count}</td>
              <td>{new Date(o.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
