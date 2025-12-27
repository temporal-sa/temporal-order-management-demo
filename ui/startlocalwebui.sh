#!/bin/bash
export TEMPORAL_ADDRESS=localhost:7233
export TEMPORAL_NAMESPACE=default
echo "Starting Web UI on http://localhost:5000 ..."
uv run app.py
