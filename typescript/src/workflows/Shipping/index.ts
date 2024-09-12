import { proxyActivities, proxyLocalActivities, sleep, workflowInfo, uuid4 } from '@temporalio/workflow';
import type * as activities from '../../activities';
import type { OrderInput, OrderItem, OrderOutput } from '../../types';

const { checkFraud, chargeCustomer, prepareShipment, shipOrder } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: {
    initialInterval: '1s',
    backoffCoefficient: 2.0,
    maximumInterval: '30s'
  }
});

export async function ShippingWorkflow(input: OrderInput, item: OrderItem) {
  await shipOrder(input, item);
}

