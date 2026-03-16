import type { Order, OrderEvent, Metrics, CreateOrderPayload } from '../types';

const BASE_URL = '/api';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.error ?? `HTTP ${res.status}`);
  }

  return data as T;
}

export const api = {
  createOrder: (payload: CreateOrderPayload): Promise<Order> =>
    request<Order>('/orders', { method: 'POST', body: JSON.stringify(payload) }),

  listOrders: (): Promise<Order[]> =>
    request<Order[]>('/orders'),

  getOrder: (id: string): Promise<Order> =>
    request<Order>(`/orders/${id}`),

  getOrderEvents: (id: string): Promise<OrderEvent[]> =>
    request<OrderEvent[]>(`/orders/${id}/events`),

  getMetrics: (): Promise<Metrics> =>
    request<Metrics>('/metrics'),
};
