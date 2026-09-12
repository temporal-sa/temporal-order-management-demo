#!/bin/bash
export EXTERNAL_STORAGE=true
export S3_ENDPOINT=http://localhost:9000
export S3_BUCKET=temporal-payloads
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=minioadmin
export AWS_SECRET_ACCESS_KEY=minioadmin
npm install
npm run codec-server
