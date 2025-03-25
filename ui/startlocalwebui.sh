#!/bin/bash


echo "Starting Web UI on http://localhost:5000 ..."
echo "Setting env vars for a local Temporal service."
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=orders
TEMPORAL_ADDRESS=localhost:7233
env | grep TEMP
poetry install --no-root
poetry run python app.py
