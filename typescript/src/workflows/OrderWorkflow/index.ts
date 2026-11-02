import {
  proxyActivities,
  proxyLocalActivities,
  sleep,
  workflowInfo,
  defineQuery,
  setHandler,
  log,
  uuid4,
} from '@temporalio/workflow';
import type * as activities from '../../activities/index';
import type { RetryPolicy } from '@temporalio/client';
import type { OrderInput, OrderOutput } from '../../types';

const DEFAULT_RETRY_POLICY: RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
};

const { getItems } = proxyLocalActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY,
});

const { checkFraud, chargeCustomer, prepareShipment } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY,
});

const { shipOrder } = proxyActivities<typeof activities>({
  // Item n takes n*2s; leave headroom so start-to-close doesn't fire first
  startToCloseTimeout: '15s',
  // A dead worker is detected in ~2s instead of waiting out the full start-to-close
  heartbeatTimeout: '2s',
  retry: DEFAULT_RETRY_POLICY,
});

const getProgressQuery = defineQuery<number>('getProgress');

export async function OrderWorkflow(input: OrderInput): Promise<OrderOutput> {
  let progress = 0;
  const { workflowType } = workflowInfo();
  log.info(`Order workflow started, ${workflowType}, ${input.OrderId}`);

  // getProgress query handler
  setHandler(getProgressQuery, () => {
    return progress;
  });

  // Get Items
  const orderItems = await getItems();

  // Check Fraud
  await checkFraud(input);
  progress = await doSleep(1, 25);

  // Prepare Shipment
  await prepareShipment(input);
  progress = await doSleep(1, 50);

  // Charge Customer
  await chargeCustomer(input, workflowType);
  progress = await doSleep(3, 75);

  // Ship order items. 
  const shipOrderPromises = [];
  for (const [i, orderItem] of orderItems.entries()) {
    log.info(`Shipping item: ${orderItem.description}`);
    shipOrderPromises.push(shipOrder({ Order: input, Item: orderItem, Index: i + 1 }));
  }

  // Wait for all items to ship
  await Promise.all(shipOrderPromises);
  progress = await doSleep(0, 100);

  const trackingId = uuid4();
  return { trackingId, address: input.Address };
}

async function doSleep(seconds: number, progress: number): Promise<number> {
  if (seconds > 0) {
    await sleep(`${seconds}s`);
  }
  return progress;
}

export const OrderWorkflowHappyPath = OrderWorkflow;
