import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import type { OrderType } from '../types';
import { api } from '../services/api';
import { Navbar } from '../components/Navbar';

export function SubmitOrderPage() {
  const navigate = useNavigate();
  const [asset, setAsset] = useState('');
  const [quantity, setQuantity] = useState('');
  const [orderType, setOrderType] = useState<OrderType>('MARKET');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    const qty = parseInt(quantity, 10);
    if (!asset.trim()) { setError('Asset symbol is required.'); return; }
    if (isNaN(qty) || qty <= 0) { setError('Quantity must be a positive integer.'); return; }

    setSubmitting(true);
    try {
      const order = await api.createOrder({ asset: asset.toUpperCase(), quantity: qty, order_type: orderType });
      setSuccess(`Order ${order.id.slice(0, 8)}… submitted successfully!`);
      setTimeout(() => navigate(`/orders/${order.id}`), 1500);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to submit order');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Navbar />
      <main className="page">
        <div className="page-header">
          <div>
            <h1 className="page-title">Submit Order</h1>
            <p className="page-sub">Create a new investment order for processing</p>
          </div>
        </div>

        <div className="form-card">
          <form onSubmit={handleSubmit} className="order-form">
            {error && <div className="form-error">{error}</div>}
            {success && <div className="form-success">{success}</div>}

            <div className="form-row">
              <div className="form-group">
                <label htmlFor="asset">Asset Symbol</label>
                <input
                  id="asset"
                  type="text"
                  className="form-input"
                  placeholder="e.g. AAPL, BTC, NVDA"
                  value={asset}
                  onChange={(e) => setAsset(e.target.value)}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="quantity">Quantity</label>
                <input
                  id="quantity"
                  type="number"
                  className="form-input"
                  placeholder="e.g. 100"
                  min={1}
                  value={quantity}
                  onChange={(e) => setQuantity(e.target.value)}
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label>Order Type</label>
              <div className="type-selector">
                {(['MARKET', 'LIMIT', 'STOP'] as OrderType[]).map((t) => (
                  <button
                    key={t}
                    type="button"
                    className={`type-btn ${orderType === t ? 'active' : ''}`}
                    onClick={() => setOrderType(t)}
                  >
                    {t}
                  </button>
                ))}
              </div>
            </div>

            <div className="form-preview">
              <span>Preview:</span>
              <strong>{quantity || '0'}</strong> units of <strong>{asset.toUpperCase() || 'ASSET'}</strong> — <strong>{orderType}</strong>
            </div>

            <button type="submit" className="btn btn-primary btn-full" disabled={submitting}>
              {submitting ? 'Submitting…' : 'Submit Order'}
            </button>
          </form>
        </div>
      </main>
    </>
  );
}
