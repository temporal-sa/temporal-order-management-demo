#!/bin/bash
source ../setcloudenv.sh

if [ -n "$TEMPORAL_API_KEY" ]; then
    export SPRING_PROFILES_ACTIVE=tc-apikey
else
    export SPRING_PROFILES_ACTIVE=tc-mtls
fi

./gradlew :core:bootRun
