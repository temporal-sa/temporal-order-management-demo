#!/bin/bash
poetry install --no-root
poetry run python worker.py
