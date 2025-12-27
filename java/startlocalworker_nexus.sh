#!/bin/bash
export TEMPORAL_NEXUS_ADDRESS=localhost:7233
export TEMPORAL_NEXUS_NAMESPACE=nexus-demo
./gradlew :shipping-service:bootRun
