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

const { checkFraud, chargeCustomer, prepareShipment, shipOrder } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
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

  // Ship Order
  const shipOrderPromises = [];
  for (const orderItem of orderItems) {
    log.info(`Shipping item: ${orderItem.description}`);
    shipOrderPromises.push(shipOrder(input, orderItem));
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
