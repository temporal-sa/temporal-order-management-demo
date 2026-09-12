import dataclasses
import os
from contextlib import AsyncExitStack
from temporalio.client import Client
from temporalio.converter import DataConverter, ExternalStorage
from temporalio.envconfig import ClientConfig

_exit_stack = AsyncExitStack()

async def get_client() -> Client:
    connect_config = ClientConfig.load_client_connect_config()
    if os.getenv("EXTERNAL_STORAGE") == "true":
        connect_config["data_converter"] = await _external_storage_converter()
    client = await Client.connect(**connect_config)
    print(f"✅ Client connected to {client.service_client.config.target_host} in namespace '{client.namespace}'")
    return client

async def close_client() -> None:
    await _exit_stack.aclose()

# Retrieve only - at the default threshold the UI offloads nothing itself
async def _external_storage_converter() -> DataConverter:
    import aioboto3
    from temporalio.contrib.aws.s3driver import S3StorageDriver
    from temporalio.contrib.aws.s3driver.aioboto3 import new_aioboto3_client

    s3 = await _exit_stack.enter_async_context(
        aioboto3.Session().client(
            "s3", endpoint_url=os.getenv("S3_ENDPOINT"), region_name=os.getenv("AWS_REGION")
        )
    )
    driver = S3StorageDriver(
        client=new_aioboto3_client(s3),
        bucket=os.getenv("S3_BUCKET", "temporal-payloads"),
    )
    return dataclasses.replace(
        DataConverter.default, external_storage=ExternalStorage(drivers=[driver])
    )
