export interface OrderInput {
  OrderId: string;
  Address: string;
}

export interface OrderItem {
  id: number;
  description: string;
  quantity: number;
}

export interface OrderOutput {
  trackingId: string;
  address: string;
}

export interface UpdateOrderInput {
  Address: string;
}
