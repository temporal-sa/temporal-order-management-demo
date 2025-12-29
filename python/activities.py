import asyncio
import logging
import random
import time
from datetime import timedelta

from temporalio import activity
from temporalio.exceptions import ApplicationError
from temporalio.common import RetryPolicy

from shared_objects import OrderItem, OrderInput

logging.basicConfig(level=logging.INFO)


class OrderActivities:

    ERROR_CHARGE_API_UNAVAILABLE = "OrderWorkflowAPIFailure"
    ERROR_INVALID_CREDIT_CARD = "OrderWorkflowNonRecoverableFailure"

    retry_policy = RetryPolicy(initial_interval=timedelta(seconds=1), backoff_coefficient=2, maximum_interval=timedelta(seconds=30))

    async def simulate_external_operation(self, ms: int):
        try:
            await asyncio.sleep(ms / 1000)
        except InterruptedError as e:
            print(e.__traceback__)

    async def simulate_external_operation_charge(self, ms: int, type: str, attempt: int):
        await self.simulate_external_operation(ms / attempt)
        return type if attempt < 5 else "NoError"

    @activity.defn
    async def get_items(self) -> list[OrderItem]:
        activity.logger.info("Getting list of items")

        await self.simulate_external_operation(100)

        items = [
            OrderItem(654300, "Table Top", 1),
            OrderItem(654321, "Table Legs", 2),
            OrderItem(654322, "Keypad", 1)
            ]

        return items

    @activity.defn
    async def check_fraud(self, input: OrderInput) -> str:
        activity.logger.info(f"Check Fraud activity started, orderId = {input.OrderId}")

        # Simulate external API call
        await self.simulate_external_operation(1000)

        return input.OrderId

    @activity.defn
    async def prepare_shipment(self, input: OrderInput) -> str:
        activity.logger.info(f"Prepare Shipment activity started, orderId = {input.OrderId}")

        # Simulate external API call
        await self.simulate_external_operation(1000)

        return input.OrderId

    @activity.defn
    async def charge_customer(self, input: OrderInput, type: str) -> str:
        activity.logger.info(f"Charge Customer activity started, orderId = {input.OrderId}")
        attempt = activity.info().attempt

        # Simulate external API call
        error = await self.simulate_external_operation_charge(1000, type, attempt)
        activity.logger.info(f"Simulated call complete, type = {type}, error = {error}")
        match error:
            case self.ERROR_CHARGE_API_UNAVAILABLE:
                # a transient error, which can be retried
                activity.logger.info(f"Charge Customer API unavailable, attempt = {attempt}")
                raise ApplicationError("Charge Customer activity failed, API unavailable")
            case self.ERROR_INVALID_CREDIT_CARD:
                # a business error, which cannot be retried
                raise ApplicationError("Charge Customer activity failed, card is invalid", type="InvalidCreditCard", non_retryable=True)
            case _:
                # pass through, no error
                pass

        return input.OrderId

    @activity.defn
    async def ship_order(self, input: OrderInput, item: OrderItem) -> None:
        activity.logger.info(f"Ship Order activity started, orderId = {input.OrderId}, itemId = {item.id}, itemDescription = {item.description}")
        randTime = random.randint(1000, 4000)
        activity.logger.info(f"Shipping Delay Time:  {randTime}")
        # Simulate external API call
        await self.simulate_external_operation(randTime)

    @activity.defn
    async def undo_prepare_shipment(self, input: OrderInput) -> str:
        activity.logger.info(f"Undo Prepare Shipment activity started, orderId = {input.OrderId}")

        # Simulate external API call
        await self.simulate_external_operation(1000)

        return input.OrderId

    @activity.defn
    async def undo_charge_customer(self, input: OrderInput) -> str:
        activity.logger.info(f"Undo Charge Customer activity started, orderId = {input.OrderId}")

        # Simulate external API call
        await self.simulate_external_operation(1000)

        return input.OrderId
