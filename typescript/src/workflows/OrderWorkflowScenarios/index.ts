import { 
  proxyActivities, 
  proxyLocalActivities, 
  sleep, 
  defineQuery,
  defineSignal, 
  defineUpdate,
  workflowInfo, 
  uuid4, 
  upsertSearchAttributes,
  ActivityFailure,
  executeChild, 
  setHandler,
  log,
  condition,
} from '@temporalio/workflow';
import type * as activities from '../../activities/index';
import type { RetryPolicy } from '@temporalio/client';
import type { OrderInput, OrderOutput, UpdateOrderInput } from '../../types';
import { ShippingWorkflow } from '../Shipping';

const DEFAULT_RETRY_POLICY:RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
}

const { getItems } = proxyLocalActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY
});

const { 
  checkFraud, 
  chargeCustomer,
  undoChargeCustomer,
  prepareShipment, 
  undoPrepareShipment,
  shipOrder } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY
});

const GET_PROGRESS_QUERY = defineQuery<number>('getProgress');
const UPDATE_ORDER_SIGNAL = defineSignal<[UpdateOrderInput]>('UpdateOrder');
const UPDATE_ORDER_UPDATE = defineUpdate<string, [UpdateOrderInput]>('UpdateOrder');

interface Compensation {
  message: string;
  fn: () => Promise<void>;
}

export async function OrderWorkflowScenarios(input: OrderInput): Promise<OrderOutput>{
  // Defining Workflow's Values
  let progress = 0;
  let orderStatus = '';
  let signalFired = false;
  let updateFired = false;

  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });

  const {workflowType} = workflowInfo();
  const compensations = [];

  log.info(`Dynamic Order workflow started, ${workflowType}, ${input.OrderId}`);

  // Defining Workflow's Handlers
  setHandler(GET_PROGRESS_QUERY, () => {
    return progress;
  })

  setHandler(UPDATE_ORDER_SIGNAL, (update_input: UpdateOrderInput) => {
    log.info(`Received update order signal with address: ${update_input.Address}`);
    signalFired = true;
    input.Address = update_input.Address;
  });

  setHandler(UPDATE_ORDER_UPDATE, (update_input: UpdateOrderInput) => {
    log.info(`Received update order update with address: ${update_input.Address}`);
    input.Address = update_input.Address;
    return `Updated address: ${input.Address}`;
  }, { 
    validator: (update_input: UpdateOrderInput): void => {
      const { Address } = update_input;

      if(Address.length > 0) {
        const firstChar = Address[0];

        if(Number(firstChar)) {
          log.info(`Order update address is valid: ${Address}`);
          updateFired = true;
        } else {
          log.error(`Rejecting order update, invalid address: ${Address}`);
          
          throw new Error(`Address must start with a digit`);
        }
      } else {
        log.error(`Rejecting order update, invalid address: ${Address}`);

        throw new Error(`Address can not be blank`);
      }
    } 
  });

  // Start of the Workflow

  // Get Items
  const orderItems = await getItems();

  // Check Fraud
  progress = 0;
  orderStatus = "Check Fraud";
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });

  await checkFraud(input);

  // Prepare Shipment
  progress = 25;
  orderStatus = 'Prepare Shipment';
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });
  await sleep('1s');

  compensations.unshift({
    message: prettyErrorMessage('reversing shipment.'),
    fn: async () => { 
      await undoPrepareShipment(input) 
    }
  });

  await prepareShipment(input);

  // Charge Customer
  progress = 50;
  orderStatus = 'Charge Customer';
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });
  await sleep('1s');

  try {
    compensations.unshift({
      message: prettyErrorMessage('reversing charge.'),
      fn: async () => {
        await undoChargeCustomer(input)
      }
    });

    await chargeCustomer(input, workflowType);
  } catch (err) {
    log.error(`Failed to charge customer ${err}`);
    /*if (err instanceof ActivityFailure && err.cause instanceof ApplicationFailure) {
      log.error(err.message);
    }*/
    // an error occurred so call compensations
    await compensate(compensations);
    throw err;
  }

  // Ship Order
  progress = 75;
  orderStatus = 'Ship Order';
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });
  await sleep('3s');

  if(workflowType == 'OrderWorkflowRecoverableFailure') {
    // Simulate a bug
    throw new Error('Simulated bug - fix me!');
  }
  
  // wait_for_updated_address_or_timeout
  if(workflowType == 'OrderWorkflowHumanInLoopSignal' || 
    workflowType == 'OrderWorkflowHumanInLoopUpdate') {
    log.info('Waiting up to 60 seconds for updated address');

    if(await condition(() => signalFired || updateFired, '60s')) {
      // Signal and Update was fired
    } else {
      //throw ApplicationFailure.create({message : 'Updated address was not received within 60 seconds.', type: 'timeout'});
    }
  }

  const shipOrderWorkflows = orderItems.map((anItem) => {
    log.info(`Shipping item: ${anItem.description}`);

    if(workflowType == 'OrderWorkflowChildWorkflow') {
      return executeChild(ShippingWorkflow, {
        args:[input, anItem]
      });
    } else {
      return shipOrder(input, anItem);
    }
  })

  await Promise.all(shipOrderWorkflows);

  // Order Completed
  progress = 100;
  orderStatus = 'Order Completed';
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });

  const trackingId = uuid4();
  return {trackingId, address: input.Address};
}

/**
 * Helper Methods  
 */
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


/**
 * "Dyanmic Workflows"
 */
export const OrderWorkflowRecoverableFailure = OrderWorkflowScenarios;
export const OrderWorkflowChildWorkflow = OrderWorkflowScenarios;
export const OrderWorkflowHumanInLoopSignal = OrderWorkflowScenarios;
export const OrderWorkflowHumanInLoopUpdate = OrderWorkflowScenarios;
export const OrderWorkflowAdvancedVisibility = OrderWorkflowScenarios;
export const OrderWorkflowAPIFailure = OrderWorkflowScenarios;
export const OrderWorkflowNonRecoverableFailure = OrderWorkflowScenarios;