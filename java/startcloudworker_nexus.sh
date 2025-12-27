#!/bin/bash
source ../setcloudenv.sh

if [ -n "$TEMPORAL_NEXUS_API_KEY" ]; then
    export SPRING_PROFILES_ACTIVE=tc-apikey
elif [ -n "$TEMPORAL_NEXUS_TLS_CLIENT_CERT_PATH" ]; then
    export SPRING_PROFILES_ACTIVE=tc-mtls
elif [ -n "$TEMPORAL_API_KEY" ]; then
    export SPRING_PROFILES_ACTIVE=tc-apikey
else
    export SPRING_PROFILES_ACTIVE=tc-mtls
fi

./gradlew :shipping-service:bootRun
