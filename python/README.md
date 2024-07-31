# Temporal Order Management Demo - Python

This is an alternate implementation of the Temporal Order Management Demo backend
using the [Python SDK](https://github.com/temporalio/sdk-python).

All of the scenarios suppported in the Go backend, and outlined in the main README are implemented in
this Python version. This version is also fully compatible with the Python UI. See the main README for
instructions on how to run the UI, and the instructions below for running the Java backend.

## Enable Updates

In order to allow the temporal server permission for Workflow Execution Updates, start your temporal server as follows:
```bash
$ temporal server start-dev --dynamic-config-value frontend.enableUpdateWorkflowExecution=true
```

## Run Worker

```bash
$ cd python

$ poetry install --no-root

$ poetry run python worker.py
```

