import * as nexus from 'nexus-rpc';
import type { ShippingInput } from './types';

// Nexus service contract for shipping. The names must match the handler
// (see go/app/service.go): service "shipping-service", operation "ship-item".
export const shippingService = nexus.service('shipping-service', {
  shipItem: nexus.operation<ShippingInput, string>({ name: 'ship-item' }),
});
