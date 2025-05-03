#!/bin/bash
echo "Starting Web UI on http://localhost:5000 ..."
poetry install --no-root
poetry run python app.py
