#!/bin/bash
# pass in the temporal environment to use to pick up the connection strings as the first parameter to this script or 
# ensure that the environment variable TEMPORAL_ENV is set appropriately.
source ../setcloudenv.sh $1

echo "Starting Web UI on http://localhost:5000 ..."
poetry run python app.py
