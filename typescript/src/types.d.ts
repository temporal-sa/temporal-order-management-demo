export interface OrderInput {
  OrderId: string;
  Address: string;
}

export interface OrderItem {
  id: number;
  description: string;
  quantity: number;
}

export interface ShippingInput {
  Order: OrderInput;
  Item: OrderItem;
  Index?: number; // item n ships in n*2 seconds; omitted means a random delay
}

export interface OrderOutput {
  trackingId: string;
  address: string;
}

export interface UpdateOrderInput {
  Address: string;
}
