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
  SearchAttributes, 
  ActivityFailure,
  executeChild, 
  setHandler,
  log
} from '@temporalio/workflow';
import type * as activities from '../../activities';
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
  undoPrepareShipment } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5s',
  retry: DEFAULT_RETRY_POLICY
});

const GET_PROGRESS_QUERY = defineQuery<number>('getProgress');
const UPDATE_ORDER_SIGNAL = defineSignal<[UpdateOrderInput]>('UpdateOrder');
const UPDATE_ORDER_UPDATE = defineUpdate<UpdateOrderInput, [string]>('UpdateOrder');

export async function OrderWorkflowScenarios(input: OrderInput): Promise<SearchAttributes> {
  // Defining Workflow's Values
  let progress = 0;
  let orderStatus = '';

  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });

  const {workflowType} = workflowInfo();
  const compensations = [];
  
  // Defining Workflow's Handlers
  setHandler(GET_PROGRESS_QUERY, () => {
    return progress;
  })

  setHandler(UPDATE_ORDER_SIGNAL, (update_input: UpdateOrderInput) => {
    input.Address = update_input.Address;
  });
  const validator = (arg: number) => {
    if (arg < 0) {
      throw new Error('Argument must not be negative');
    }
  };

  setHandler(UPDATE_ORDER_UPDATE, (update_input: UpdateOrderInput) => {
    input.Address = update_input.Address;
    return `Updated address: ${input.Address}`;
  }, { validator });

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
    fn: () => undoPrepareShipment(input)
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
      fn: () => undoChargeCustomer(input)
    });

    await chargeCustomer(input, workflowType);
  } catch (err) {
    if (err instanceof ActivityFailure && err.cause instanceof ApplicationFailure) {
      log.error(err.cause.message);
    } else {
      log.error(`error while opening account: ${err}`);
    }
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

  if(workflowType == OrderWorkflowRecoverableFailure.toString()) {
    // Simulate a bug
    // throw new Error('Workflow bug!');
  }

  const shipOrderWorkflows = orderItems.map((anItem) => {
    return executeChild(ShippingWorkflow, {
      args:[input, anItem]
    });
  })

  await Promise.all(shipOrderWorkflows);

  // Order Completed
  progress = 100;
  orderStatus = 'Order Completed';
  upsertSearchAttributes({
    OrderStatus: [orderStatus]
  });

  const trackingId = uuid4();

  return workflowInfo().searchAttributes;
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