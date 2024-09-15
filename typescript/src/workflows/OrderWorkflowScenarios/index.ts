import {
  proxyActivities,
  proxyLocalActivities,
  sleep,
  workflowInfo,
  defineQuery,
  defineSignal,
  defineUpdate,
  upsertSearchAttributes,
  ActivityFailure,
  executeChild,
  setHandler,
  log,
  uuid4,
  condition,
  ParentClosePolicy,
} from '@temporalio/workflow';
import type * as activities from '../../activities/index';
import type { RetryPolicy } from '@temporalio/client';
import type { OrderInput, OrderItem, OrderOutput, UpdateOrderInput } from '../../types';
import { ShippingWorkflow } from '../Shipping';

const DEFAULT_RETRY_POLICY: RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
};

const { getItems } = proxyLocalActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY,
});

const { checkFraud, chargeCustomer, undoChargeCustomer, prepareShipment, undoPrepareShipment, shipOrder } =
  proxyActivities<typeof activities>({
    startToCloseTimeout: '5s',
    retry: DEFAULT_RETRY_POLICY,
  });

const getProgressQuery = defineQuery<number>('getProgress');
const updateOrderSignal = defineSignal<[UpdateOrderInput]>('UpdateOrder');
const updateOrderUpdate = defineUpdate<string, [UpdateOrderInput]>('UpdateOrder');

interface Compensation {
  message: string;
  fn: () => Promise<void>;
}

export async function OrderWorkflowScenarios(input: OrderInput): Promise<OrderOutput> {
  let progress = 0;
  let updatedAddress = '';

  const { workflowType } = workflowInfo();
  log.info(`Scenario Order workflow started, ${workflowType}, ${input.OrderId}`);

  // Saga compensations
  const compensations = [];

  // getProgress query handler
  setHandler(getProgressQuery, () => {
    return progress;
  });

  setHandler(updateOrderSignal, (update_input: UpdateOrderInput) => {
    log.info(`Received update order signal with address: ${update_input.Address}`);
    updatedAddress = update_input.Address;
  });

  setHandler(
    updateOrderUpdate,
    (update_input: UpdateOrderInput) => {
      log.info(`Received update order update with address: ${update_input.Address}`);
      updatedAddress = update_input.Address;
      return `Updated address: ${update_input.Address}`;
    },
    {
      validator: (update_input: UpdateOrderInput): void => {
        const { Address } = update_input;

        if (Address.length > 0) {
          const firstChar = Address[0];

          if (Number(firstChar)) {
            log.info(`Order update address is valid: ${Address}`);
          } else {
            log.error(`Rejecting order update, invalid address: ${Address}`);
            throw new Error(`Address must start with a digit`);
          }
        } else {
          log.error(`Rejecting order update, invalid address: ${Address}`);
          throw new Error(`Address can not be blank`);
        }
      },
    },
  );

  // Get Items
  const orderItems = await getItems();

  progress = await updateProgress('Check Fraud', 0, 0);

  // Check Fraud
  await checkFraud(input);

  progress = await updateProgress('Prepare Shipment', 25, 1);

  // Prepare Shipment
  compensations.unshift({
    message: prettyErrorMessage('reversing shipment.'),
    fn: async () => {
      await undoPrepareShipment(input);
    },
  });
  await prepareShipment(input);

  progress = await updateProgress('Charge Customer', 50, 1);

  // Charge Customer
  try {
    compensations.unshift({
      message: prettyErrorMessage('reversing charge.'),
      fn: async () => {
        await undoChargeCustomer(input);
      },
    });
    await chargeCustomer(input, workflowType);
  } catch (err) {
    log.error(`Failed to charge customer ${err}`);
    // an error occurred so call compensations
    await compensate(compensations);
    throw err;
  }

  progress = await updateProgress('Ship Order', 75, 3);

  if (WF_TYPES.BUG == workflowType) {
    // Simulate a bug
    throw new Error('Simulated bug - fix me!');
  }

  if (WF_TYPES.SIGNAL == workflowType || WF_TYPES.UPDATE == workflowType) {
    // Await message to update address
    await waitForUpdatedAddressOrTimeout();
  }

  // Ship Order
  const shipOrderPromises = [];
  for (const orderItem of orderItems) {
    log.info(`Shipping item: ${orderItem.description}`);
    shipOrderPromises.push(shipItemAsync(input, orderItem, workflowType));
  }

  // Wait for all items to ship
  await Promise.all(shipOrderPromises);

  progress = await updateProgress('Order Completed', 100, 1);

  const trackingId = uuid4();
  return { trackingId, address: input.Address };

  async function waitForUpdatedAddressOrTimeout() {
    log.info('Waiting up to 60 seconds for updated address');

    if (await condition(() => updatedAddress != '', '60s')) {
      input.Address = updatedAddress;
    } else {
      // Do nothing - use the original address
      // In other cases, you may want to throw an exception on timeout, e.g.
      //   throw ApplicationFailure.create({message : 'Updated address was not received within 60 seconds.', type: 'timeout'});
    }
  }
}

async function shipItemAsync(input: OrderInput, orderItem: OrderItem, type: string): Promise<void> {
  if (WF_TYPES.CHILD == type) {
    return executeChild(ShippingWorkflow, {
      args: [input, orderItem],
      workflowId: `shipment-${input.OrderId}-${orderItem.id}`,
      parentClosePolicy: ParentClosePolicy.PARENT_CLOSE_POLICY_TERMINATE,
    });
  } else {
    return shipOrder(input, orderItem);
  }
}

async function updateProgress(orderStatus: string, progress: number, seconds: number): Promise<number> {
  if (seconds > 0) {
    await sleep(`${seconds}s`);
  }
  if (WF_TYPES.VISIBILITY == workflowInfo().workflowType) {
    upsertSearchAttributes({
      OrderStatus: [orderStatus],
    });
  }
  return progress;
}

function prettyErrorMessage(message: string, err?: any) {
  let errMessage = err && err.message ? err.message : '';
  if (err && err instanceof ActivityFailure) {
    errMessage = `${err.cause?.message}`;
  }
  return `${message}: ${errMessage}`;
}

async function compensate(compensations: Compensation[] = []) {
  if (compensations.length > 0) {
    log.info('failures encountered during account opening - compensating');
    for (const comp of compensations) {
      try {
        log.error(comp.message);
        await comp.fn();
      } catch (err) {
        log.error(`failed to compensate: ${prettyErrorMessage('', err)}`, { err });
        // swallow errors
      }
    }
  }
}

// Exported Workflow Scenarios
export const OrderWorkflowRecoverableFailure = OrderWorkflowScenarios;
export const OrderWorkflowChildWorkflow = OrderWorkflowScenarios;
export const OrderWorkflowHumanInLoopSignal = OrderWorkflowScenarios;
export const OrderWorkflowHumanInLoopUpdate = OrderWorkflowScenarios;
export const OrderWorkflowAdvancedVisibility = OrderWorkflowScenarios;
export const OrderWorkflowAPIFailure = OrderWorkflowScenarios;
export const OrderWorkflowNonRecoverableFailure = OrderWorkflowScenarios;

const WF_TYPES = {
  BUG: 'OrderWorkflowRecoverableFailure',
  CHILD: 'OrderWorkflowChildWorkflow',
  SIGNAL: 'OrderWorkflowHumanInLoopSignal',
  UPDATE: 'OrderWorkflowHumanInLoopUpdate',
  VISIBILITY: 'OrderWorkflowAdvancedVisibility',
} as const;
