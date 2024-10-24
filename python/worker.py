import asyncio
import os

from temporalio.client import Client, TLSConfig
from temporalio.worker import Worker

from activities import OrderActivities
from order_workflow_scenarios import OrderWorkflowScenarios
from order_workflow import OrderWorkflow
from shipping_child_workflow import ShippingChildWorkflow


async def main():
    address = os.getenv("TEMPORAL_ADDRESS","127.0.0.1:7233")
    namespace = os.getenv("TEMPORAL_NAMESPACE","default")
    tlsCertPath = os.getenv("TEMPORAL_CERT_PATH","")
    tlsKeyPath = os.getenv("TEMPORAL_KEY_PATH","")
    tls = None

    if tlsCertPath and tlsKeyPath:
        with open(tlsCertPath,"rb") as f:
            cert = f.read()
        with open(tlsKeyPath,"rb") as f:
            key = f.read()

        tls = TLSConfig(client_cert=cert,
                        client_private_key=key)

    client = await Client.connect(
        target_host=address,
        namespace=namespace,
        tls=tls
    )

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
    print(f"Connecting to Temporal on {address}")
    print("Python order management worker starting...")
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
