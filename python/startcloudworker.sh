#!/bin/bash
source ../setcloudenv.sh
poetry install --no-root
poetry run python worker.py
