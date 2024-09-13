import { 
  proxyActivities, 
  proxyLocalActivities, 
  sleep, 
  workflowInfo, 
  defineQuery, 
  setHandler,
  log, 
  uuid4 } from '@temporalio/workflow';
import type * as activities from '../../activities';
import type { RetryPolicy } from '@temporalio/client';
import type { OrderInput, OrderOutput } from '../../types';

const DEFAULT_RETRY_POLICY:RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
}

const { getItems } = proxyLocalActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY
});

const { checkFraud, chargeCustomer, prepareShipment, shipOrder } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY
});

const GET_PROGRESS_QUERY = defineQuery<number>('getProgress');

export async function OrderWorkflow(input: OrderInput): Promise<OrderOutput> {
  // Defining Workflow's Values
  let progress = 0;
  const { workflowType } = workflowInfo();

  log.info(`Order workflow started, ${workflowType}, ${input.OrderId}`);

  // Defining Workflow's Handlers
  setHandler(GET_PROGRESS_QUERY, () => {
    return progress;
  })

  // Get Items
  const orderItems = await getItems();

  // Check Fraud
  await checkFraud(input);

  await sleep('1s');
  progress = 25;

  // Prepare Shipment
  await prepareShipment(input);

  await sleep('1s');
  progress = 50;

  // Charge Customer
  await chargeCustomer(input, workflowType);

  await sleep('3s');
  progress = 75;

  // Ship Orders 
  const shipOrderActivites = [];

  for(const anItem of orderItems) {
    log.info(`Shipping item: ${anItem.description}`);
    shipOrderActivites.push(
      shipOrder(input, anItem)
    )
  }

  await Promise.all(shipOrderActivites);

  await sleep(1);
  progress = 100;

  const trackingId = uuid4();

  return {trackingId, address: input.Address};
}

export const OrderWorkflowHappyPath = OrderWorkflow;