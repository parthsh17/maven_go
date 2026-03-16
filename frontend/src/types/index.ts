

export type OrderType = 'MARKET' | 'LIMIT' | 'STOP';

export type OrderState =
  | 'CREATED'
  | 'VALIDATED'
  | 'QUEUED'
  | 'EXECUTING'
  | 'COMPLETED'
  | 'FAILED'
  | 'RETRYING';

export interface Order {
  id: string;
  asset: string;
  quantity: number;
  order_type: OrderType;
  state: OrderState;
  retry_count: number;
  created_at: string;
  updated_at: string;
}

export interface OrderEvent {
  order_id: string;
  previous_state: string;
  new_state: string;
  timestamp: string;
  message?: string;
}

export interface Metrics {
  total_orders: number;
  processing_orders: number;
  completed_orders: number;
  failed_orders: number;
  worker_count: number;
  [key: string]: number;
}

export interface CreateOrderPayload {
  asset: string;
  quantity: number;
  order_type: OrderType;
}

export interface ApiError {
  error: string;
}
