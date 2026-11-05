#!/bin/bash
# MinIO stands in for S3. mc ships inside the server image, so creating the bucket needs no second container.
docker rm -fv minio >/dev/null 2>&1
docker run -d --name minio -p 9000:9000 -p 9001:9001 -v minio-data:/data \
  -e MINIO_ROOT_USER=minioadmin -e MINIO_ROOT_PASSWORD=minioadmin -e MINIO_REGION_NAME=us-east-1 \
  docker.io/pgsty/silo:RELEASE.2026-10-26T00-00-00Z server /data --console-address ":9001" >/dev/null || exit 1

# Retrying the bucket doubles as the readiness check - it fails until the server is accepting requests.
for _ in $(seq 30); do
  if docker exec -e MC_HOST_local=http://minioadmin:minioadmin@localhost:9000 minio \
      mc mb --ignore-existing local/temporal-payloads >/dev/null 2>&1; then
    echo "MinIO console on http://localhost:9001 (minioadmin/minioadmin) ..."
    exit 0
  fi
  sleep 1
done

echo "MinIO did not come up - check: docker logs minio" >&2
exit 1
