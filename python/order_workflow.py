import logging
from datetime import timedelta
import asyncio

from temporalio import workflow

from activities import OrderActivities
from shared_objects import OrderInput, OrderOutput

logging.basicConfig(level=logging.INFO)


@workflow.defn(name="OrderWorkflowHappyPath")
class OrderWorkflow:

    def __init__(self) -> None:
        self.progress = 0
        self.retry_policy = OrderActivities.retry_policy

    @workflow.run
    async def execute(self, input: OrderInput) -> OrderOutput:
        workflow_type = workflow.info().workflow_type
        workflow.logger.info("Order workflow started, " + workflow_type + ", " + input.OrderId)

        order_items = await workflow.execute_local_activity_method(
            OrderActivities.get_items,
            start_to_close_timeout=timedelta(seconds=5)
        )

        await workflow.execute_activity_method(
            OrderActivities.check_fraud,
            input,
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )
        await self.sleep(1, 25)

        await workflow.execute_activity_method(
            OrderActivities.prepare_shipment,
            input,
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )
        await self.sleep(1, 50)

        await workflow.execute_activity_method(
            OrderActivities.charge_customer,
            args=[input, workflow_type],
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )
        await self.sleep(3, 75)

        handles = []
        for item in order_items:
            workflow.logger.info("Shipping item: " + item.description)
            handles.append(
                workflow.execute_activity_method(
                    OrderActivities.ship_order,
                    args=[input, item],
                    start_to_close_timeout=timedelta(seconds=5),
                    retry_policy=self.retry_policy
                )
            )

        await asyncio.gather(*handles)
        await self.sleep(1, 100)

        tracking_id = str(workflow.uuid4())
        return OrderOutput(tracking_id, input.Address)

    @workflow.query(name="getProgress")
    def query_progress(self) -> int:
        return self.progress

    async def sleep(self, seconds: int, progress: int):
        await asyncio.sleep(seconds)
        self.progress = progress
