import { ApplicationFailure, Context, log } from '@temporalio/activity';
import type { OrderInput, OrderItem, ShippingInput } from '../types';
import type { RetryPolicy } from '@temporalio/client';

export const ERROR_CHARGE_API_UNAVAILABLE = 'OrderWorkflowAPIFailure';
export const ERROR_INVALID_CREDIT_CARD = 'OrderWorkflowNonRecoverableFailure';
export const NO_ERROR = 'NoError';

export const DEFAULT_RETRY_POLICY: RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
};

// Cancellation-aware sleep, so Ctrl+C on the worker interrupts in-flight activities instead of draining them
async function simulateExternalOperation(ms: number) {
  return Context.current().sleep(ms);
}

async function simulateExternalOperationCharge(ms: number, type: string, attempt: number) {
  await simulateExternalOperation(ms / attempt);
  return attempt < 5 ? type : NO_ERROR;
}

export async function getItems(): Promise<Array<OrderItem>> {
  log.info(`Getting list of items`);

  await simulateExternalOperation(100);

  const items = [
    {
      id: 654300,
      description: 'Table Top',
      quantity: 1,
    },
    {
      id: 65321,
      description: 'Table Legs',
      quantity: 2,
    },
    {
      id: 654322,
      description: 'Keypad',
      quantity: 1,
    },
  ];

  return items;
}

export async function checkFraud(input: OrderInput): Promise<string> {
  log.info(`Check Fraud activity started, ${input.OrderId}`);

  await simulateExternalOperation(1000);

  return input.OrderId;
}

export async function prepareShipment(input: OrderInput): Promise<string> {
  log.info(`"Prepare Shipment activity started, ${input.OrderId}`);

  await simulateExternalOperation(1000);

  return input.OrderId;
}

export async function chargeCustomer(input: OrderInput, type: string): Promise<string> {
  log.info(`Charge Customer activity started, ${input.OrderId}`);

  const context = Context.current();
  const { attempt } = context.info;

  const error = await simulateExternalOperationCharge(1000, type, attempt);
  log.info(`Simulated call complete, ${type}, ${error}`);

  switch (error) {
    case ERROR_CHARGE_API_UNAVAILABLE:
      log.error(`Charge Customer API unavailable, ${attempt}`);
      throw ApplicationFailure.create({
        message: 'Charge Customer activity failed, API unavailable',
      });
    case ERROR_INVALID_CREDIT_CARD:
      throw ApplicationFailure.create({
        nonRetryable: true,
        message: 'Charge Customer activity failed, card is invalid',
        type: 'InvalidCreditCard',
      });
    default:
      break;
  }

  return input.OrderId;
}

export async function shipOrder(input: ShippingInput) {
  log.info(`Ship Order activity started, ${input.Order.OrderId}, ${input.Item.id}, ${input.Item.description}`);

  // Item n ships in n*2 seconds, so a fan-out finishes one item at a time; no index means random
  const delayMs = input.Index ? input.Index * 2000 : Math.floor(Math.random() * 3001) + 1000;

  // On a retry, pick up from the last second this item reported
  const context = Context.current();
  const secondsDone: number = context.info.heartbeatDetails ?? 0;
  log.info(`Shipping Delay Time: ${delayMs}, resuming at ${secondsDone}s`);

  // Simulate external API call, heartbeating progress each second so a retry can resume
  for (let second = secondsDone; second < Math.floor(delayMs / 1000); second++) {
    await simulateExternalOperation(1000);
    context.heartbeat(second + 1);
  }
  await simulateExternalOperation(delayMs % 1000);
}

export async function undoPrepareShipment(input: OrderInput): Promise<string> {
  log.info(`Undo Prepare Shipment activity started, ${input.OrderId}`);

  await simulateExternalOperation(1000);

  return input.OrderId;
}

export async function undoChargeCustomer(input: OrderInput): Promise<string> {
  log.info(`Undo Charge Customer activity started, ${input.OrderId}`);

  await simulateExternalOperation(1000);

  return input.OrderId;
}
