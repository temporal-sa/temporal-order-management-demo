import { proxyActivities, log } from '@temporalio/workflow';
import type * as activities from '../../activities/index';
import type { OrderInput, OrderItem } from '../../types';

const { shipOrder } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: {
    initialInterval: '1s',
    backoffCoefficient: 2.0,
    maximumInterval: '30s'
  }
});

export async function ShippingWorkflow(input: OrderInput, item: OrderItem) {
  log.info(`Shipping workflow started, orderId ${input.OrderId}`);

  await shipOrder(input, item);
}

