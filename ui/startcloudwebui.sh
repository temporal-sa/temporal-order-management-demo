#!/bin/bash

source ../setcloudenv.sh

echo "Starting Web UI on http://localhost:5000 ..."
poetry run python app.py
