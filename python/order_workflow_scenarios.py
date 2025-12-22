import logging
from datetime import timedelta
import asyncio
from collections.abc import Sequence
from typing import Any

from temporalio import workflow
from temporalio.common import RawValue, SearchAttributeKey
from temporalio.exceptions import ApplicationError
from temporalio.workflow import ParentClosePolicy

from activities import OrderActivities
from shared_objects import OrderInput, OrderOutput, UpdateOrderInput, OrderItem
from shipping_child_workflow import ShippingChildWorkflow

logging.basicConfig(level=logging.INFO)


@workflow.defn(dynamic=True)
class OrderWorkflowScenarios:

    BUG = "OrderWorkflowRecoverableFailure"
    CHILD = "OrderWorkflowChildWorkflow"
    SIGNAL = "OrderWorkflowHumanInLoopSignal"
    UPDATE = "OrderWorkflowHumanInLoopUpdate"
    VISIBILITY = "OrderWorkflowAdvancedVisibility"

    ORDER_STATUS_SA = SearchAttributeKey.for_keyword("OrderStatus")

    def __init__(self) -> None:
        self.progress = 0
        self.updated_address = None
        self.retry_policy = OrderActivities.retry_policy

    @workflow.run
    async def execute(self, args: Sequence[RawValue]) -> Any:
        input = workflow.payload_converter().from_payload(args[0].payload, OrderInput)
        workflow_type = workflow.info().workflow_type
        workflow.logger.info(f"Dynamic Order workflow started, type = {workflow_type}, orderId = {input.OrderId}")

        # Saga compensations
        compensations = []

        # Get items
        order_items = await workflow.execute_local_activity_method(
            OrderActivities.get_items,
            start_to_close_timeout=timedelta(seconds=5)
        )

        await self.update_progress("Check Fraud", 0, 0)

        # Check fraud
        await workflow.execute_activity_method(
            OrderActivities.check_fraud,
            input,
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )

        await self.update_progress("Prepare Shipment", 25, 1)

        # Prepare shipment
        compensations.append(OrderActivities.undo_prepare_shipment)
        await workflow.execute_activity_method(
            OrderActivities.prepare_shipment,
            input,
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )

        await self.update_progress("Charge Customer", 50, 1)

        # Charge customer
        try:
            compensations.append(OrderActivities.undo_charge_customer)
            await workflow.execute_activity_method(
                OrderActivities.charge_customer,
                args=[input, workflow_type],
                start_to_close_timeout=timedelta(seconds=5),
                retry_policy=self.retry_policy
            )
        except Exception as ex:
            workflow.logger.error("Failed to charge customer", ex)
            for compensation in reversed(compensations):
                await workflow.execute_activity(
                    compensation,
                    input,
                    start_to_close_timeout=timedelta(seconds=10),
                    retry_policy=self.retry_policy
                )
            raise ex

        await self.update_progress("Ship Order", 75, 3)

        if self.BUG == workflow_type:
            # Simulate bug
            raise RuntimeError("Simulated bug - fix me!")
            pass

        if self.SIGNAL == workflow_type or self.UPDATE == workflow_type:
            # Await message to update address
            await self.wait_for_updated_address_or_timeout(input)

        # Ship order items
        handles = []
        for item in order_items:
            workflow.logger.info(f"Shipping item: {item.description}")
            handles.append(self.ship_item_async(input, item, workflow_type))

        # Wait for all items to ship
        await asyncio.gather(*handles)

        await self.update_progress("Order Completed", 100, 0)

        # Generate trackingId
        tracking_id = str(workflow.uuid4())
        return OrderOutput(tracking_id, input.Address)

    def ship_item_async(self, input: OrderInput, item: OrderItem, workflow_type: str) -> Any:
        if self.CHILD == workflow_type:
            return workflow.execute_child_workflow(
                ShippingChildWorkflow.execute,
                args=[input, item],
                id="shipment-" + input.OrderId + "-" + str(item.id),
                parent_close_policy=ParentClosePolicy.TERMINATE
            )

        else:
            return workflow.execute_activity_method(
                OrderActivities.ship_order,
                args=[input, item],
                start_to_close_timeout=timedelta(seconds=5),
                retry_policy=self.retry_policy
            )

    async def wait_for_updated_address_or_timeout(self, input: OrderInput):
        workflow.logger.info("Waiting up to 60 seconds for updated address")
        try:
            await workflow.wait_condition(lambda: self.updated_address != None, timeout=timedelta(seconds=60))
            input.Address = self.updated_address
        except asyncio.TimeoutError:
            pass
            # raise ApplicationFailure("Updated address was not received within 60 seconds.", type="timeout")

    async def update_progress(self, order_status: str, progress: int, sleep: int):
        self.progress = progress
        if sleep > 0:
            await asyncio.sleep(sleep)
        if self.VISIBILITY == workflow.info().workflow_type:
            workflow.upsert_search_attributes([self.ORDER_STATUS_SA.value_set(order_status)])

    @workflow.query(name="getProgress")
    def query_progress(self) -> int:
        return self.progress

    @workflow.signal(name="UpdateOrder")
    def update_order_signal(self, update_input: UpdateOrderInput):
        workflow.logger.info(f"Received update order signal with address: {update_input.Address}")
        self.updated_address = update_input.Address

    @workflow.update(name="UpdateOrder")
    def update_order_update(self, update_input: UpdateOrderInput) -> str:
        workflow.logger.info(f"Received update order signal with address: {update_input.Address}")
        self.updated_address = update_input.Address
        return "Updated address: " + update_input.Address

    @update_order_update.validator
    def update_order_validator(self, update_input: UpdateOrderInput):
        if not update_input.Address[0].isdigit():
            workflow.logger.info(f"Rejecting order update, invalid address: {update_input.Address}")
            raise ApplicationError("Address must start with a digit", type="invalid-address")
        workflow.logger.info(f"Order update address is valid: {update_input.Address}")
