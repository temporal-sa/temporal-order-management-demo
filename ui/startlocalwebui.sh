#!/bin/bash
export TEMPORAL_ADDRESS=localhost:7233
export TEMPORAL_NAMESPACE=default
export EXTERNAL_STORAGE=true
export S3_ENDPOINT=http://localhost:9000
echo "Starting Web UI on http://localhost:5000 ..."
uv run app.py
