import asyncio

import os

from temporalio.client import Client
from temporalio.worker import Worker

from activities import OrderActivities
from order_workflow_scenarios import OrderWorkflowScenarios
from order_workflow import OrderWorkflow
from shipping_child_workflow import ShippingChildWorkflow


async def main():
    client = await Client.connect("localhost:7233")

    activities = OrderActivities()

    worker = Worker(
        client,
        task_queue=os.getenv("TEMPORAL_TASK_QUEUE", "orders"),
        workflows=[
            OrderWorkflow,
            OrderWorkflowScenarios,
            ShippingChildWorkflow
        ],
        activities=[
            activities.get_items,
            activities.check_fraud,
            activities.prepare_shipment,
            activities.charge_customer,
            activities.ship_order,
            activities.undo_prepare_shipment,
            activities.undo_charge_customer
        ],
    )
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
