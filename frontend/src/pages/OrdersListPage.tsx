import { useState } from 'react';
import { Navbar } from '../components/Navbar';
import { OrderTable } from '../components/OrderTable';
import { useOrders } from '../hooks/useOrders';

const ALL_STATES = ['ALL', 'CREATED', 'VALIDATED', 'QUEUED', 'EXECUTING', 'COMPLETED', 'FAILED', 'RETRYING'];

export function OrdersListPage() {
  const { orders, loading, error, refetch } = useOrders();
  const [filterState, setFilterState] = useState('ALL');

  return (
    <>
      <Navbar />
      <main className="page">
        <div className="page-header">
          <div>
            <h1 className="page-title">Orders</h1>
            <p className="page-sub">All orders across the lifecycle — auto-refreshes every 3s</p>
          </div>
          <button className="btn btn-ghost" onClick={refetch}>↻ Refresh</button>
        </div>

        <div className="filter-bar">
          {ALL_STATES.map((state) => (
            <button
              key={state}
              className={`filter-chip ${filterState === state ? 'active' : ''}`}
              onClick={() => setFilterState(state)}
            >
              {state}
            </button>
          ))}
        </div>

        {error && <div className="form-error">{error}</div>}

        {loading ? (
          <div className="skeleton-table" />
        ) : (
          <OrderTable orders={orders} filterState={filterState} />
        )}
      </main>
    </>
  );
}
