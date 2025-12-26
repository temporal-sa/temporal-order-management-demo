from temporalio.client import Client, TLSConfig
from typing import Optional
import os

async def get_client()-> Client:

    if (os.getenv("TEMPORAL_API_KEY") is not None):
        print(os.getenv("TEMPORAL_API_KEY"), os.getenv("TEMPORAL_NAMESPACE"), os.getenv("TEMPORAL_ADDRESS"))

        client = await Client.connect(
            os.getenv("TEMPORAL_ADDRESS"),
            namespace=os.getenv("TEMPORAL_NAMESPACE"),
            rpc_metadata={"temporal-namespace": os.getenv("TEMPORAL_NAMESPACE")},
            api_key=os.getenv("TEMPORAL_API_KEY"),
            tls=True,
        )
    elif (
        os.getenv("TEMPORAL_TLS_CLIENT_CERT_PATH")
        and os.getenv("TEMPORAL_TLS_CLIENT_KEY_PATH") is not None
    ):
        server_root_ca_cert: Optional[bytes] = None
        with open(os.getenv("TEMPORAL_TLS_CLIENT_CERT_PATH"), "rb") as f:
            client_cert = f.read()

        with open(os.getenv("TEMPORAL_TLS_CLIENT_KEY_PATH"), "rb") as f:
            client_key = f.read()

        # Start client
        client = await Client.connect(
            os.getenv("TEMPORAL_ADDRESS"),
            namespace=os.getenv("TEMPORAL_NAMESPACE"),
            tls=TLSConfig(
                server_root_ca_cert=server_root_ca_cert,
                client_cert=client_cert,
                client_private_key=client_key,
            ),
            #data_converter=dataclasses.replace(
            #    temporalio.converter.default(), payload_codec=EncryptionCodec()
            #),
        )
    else:
        client = await Client.connect(
            "localhost:7233",
        )

    return client
