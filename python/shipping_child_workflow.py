from datetime import timedelta
import logging

from temporalio import workflow

from activities import OrderActivities
from shared_objects import OrderInput, OrderItem

logging.basicConfig(level=logging.INFO)


@workflow.defn
class ShippingChildWorkflow:

    def __init__(self):
        self.retry_policy = OrderActivities.retry_policy

    @workflow.run
    async def execute(self, input: OrderInput, item: OrderItem):
        workflow.logger.info("Shipping workflow started, orderId " + input.OrderId)

        await workflow.start_activity_method(
            OrderActivities.ship_order,
            args=[input, item],
            start_to_close_timeout=timedelta(seconds=5),
            retry_policy=self.retry_policy
        )
