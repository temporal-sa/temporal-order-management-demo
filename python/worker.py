import asyncio
import concurrent.futures

from temporalio.client import Client
from temporalio.worker import Worker

from activities import OrderActivities
from workflows import OrderWorkflow, OrderWorkflowScenarios, ShippingChildWorkflow

TASK_QUEUE = "orders"


async def main():
    client = await Client.connect("localhost:7233")

    activities = OrderActivities()

    with concurrent.futures.ThreadPoolExecutor(max_workers=100) as activity_executor:
        worker = Worker(
            client,
            task_queue=TASK_QUEUE,
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
            activity_executor=activity_executor
        )
        await worker.run()


if __name__ == "__main__":
    asyncio.run(main())