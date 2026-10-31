# Temporal Order Management Demo - Python

An implementation of the Temporal Order Management Demo backend
using the [Python SDK](https://github.com/temporalio/sdk-python).

All of the scenarios outlined in the main [README](../README.md) are implemented in this Python version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## Updating the Temporal SDK

The worker scripts start the Worker with `uv run`, which installs the SDK version pinned in `uv.lock` automatically, so pulling the latest changes is usually all you need.

If you run the Worker another way (for example, `python worker.py` from an activated virtual environment), update the SDK yourself after pulling:

```bash
uv sync
# or, if you use pip instead of uv
pip install -U temporalio
```

## NexusOperation Scenario

The Python Worker implements the caller side of the NexusOperation scenario only. The `shipping-service`
Nexus Service it calls is handled by the Go Nexus Worker, so run both:

```bash
# create the nexus-demo Namespace and shipping-endpoint first (see the main README)
cd go && ./startlocalworker_nexus.sh
cd python && ./startlocalworker.sh
```
