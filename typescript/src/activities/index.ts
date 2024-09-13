import { ApplicationFailure, Context, log } from '@temporalio/activity';
import type { OrderInput, OrderItem } from '../types';
import type { RetryPolicy } from '@temporalio/client';

export const ERROR_CHARGE_API_UNAVAILABLE = 'OrderWorkflowAPIFailure';
export const ERROR_INVALID_CREDIT_CARD = 'OrderWorkflowNonRecoverableFailure'
export const NO_ERROR = 'NoError';

export const DEFAULT_RETRY_POLICY:RetryPolicy = {
  initialInterval: '1s',
  backoffCoefficient: 2,
  maximumInterval: '30s',
}

async function simulateExternalOperation(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function simulateExternalOperationCharge(ms: number, type: string, attempt: number) {
  await simulateExternalOperation(ms  / attempt);
  return attempt < 5 ? type : NO_ERROR;
}

export async function getItems(): Promise<Array<OrderItem>> {
  log.info(`Getting list of items`);

  await simulateExternalOperation(100);

  const items = [
    {
      id: 654300,
      description: 'Table Top',
      quantity: 1
    }, 
    {
      id: 65321,
      description: 'Table Legs',
      quantity: 2
    },
    {
      id: 654322,
      description: 'Keypad',
      quantity: 1
    }
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

  switch(error) {
    case ERROR_CHARGE_API_UNAVAILABLE: 
      log.error(`Charge Customer API unavailable, ${attempt}`);
      throw ApplicationFailure.create({
        message: 'Charge Customer activity failed, API unavailable'
      });
    case ERROR_INVALID_CREDIT_CARD:
      throw ApplicationFailure.create({
        nonRetryable: true,
        message: 'Charge Customer activity failed, card is invalid',
        type: 'InvalidCreditCard'
      });
    default:
      break;
  }

  return input.OrderId;
}

export async function shipOrder(input: OrderInput, item: OrderItem) {
  log.info(`Ship Order activity started, ${input.OrderId}, ${item.id}, ${item.description}`);
  
  await simulateExternalOperation(1000);
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