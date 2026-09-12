#!/bin/bash
docker compose up -d minio
docker compose run --rm minio-init
echo "MinIO console on http://localhost:9001 (minioadmin/minioadmin) ..."
