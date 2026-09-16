#!/bin/bash
export EXTERNAL_STORAGE=true
export S3_ENDPOINT=http://localhost:9000
npm install
npm run codec-server
