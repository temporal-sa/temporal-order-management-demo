import dataclasses
import os
from contextlib import AsyncExitStack
from temporalio.client import Client
from temporalio.converter import DataConverter, ExternalStorage
from temporalio.envconfig import ClientConfig

_exit_stack = AsyncExitStack()

# MinIO's region and credentials are fixed, and are set here rather than exported by the start scripts so that
# nothing in the presenter's AWS environment can reach the local demo.
MINIO = {
    "region_name": "us-east-1",
    "aws_access_key_id": "minioadmin",
    "aws_secret_access_key": "minioadmin",
}

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
    from aiobotocore.session import AioSession
    from temporalio.contrib.aws.s3driver import S3StorageDriver
    from temporalio.contrib.aws.s3driver.aioboto3 import new_aioboto3_client

    endpoint = os.getenv("S3_ENDPOINT")
    if endpoint:
        # Drop AWS_PROFILE for MinIO. Explicit credentials are not enough on their own - botocore resolves the
        # profile's config before it looks at credentials at all, so a stale one raises ProfileNotFound and the
        # Web UI never finishes starting.
        session = aioboto3.Session(botocore_session=AioSession(session_vars={"profile": (None, None, None, None)}))
        options = MINIO
    else:
        session = aioboto3.Session()
        options = {"region_name": os.getenv("AWS_REGION", "us-east-1")}

    s3 = await _exit_stack.enter_async_context(session.client("s3", endpoint_url=endpoint, **options))
    driver = S3StorageDriver(
        client=new_aioboto3_client(s3),
        bucket=os.getenv("S3_BUCKET", "temporal-payloads"),
    )
    return dataclasses.replace(
        DataConverter.default, external_storage=ExternalStorage(drivers=[driver])
    )
