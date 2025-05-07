#!/bin/bash
source ../setcloudenv.sh
./gradlew :shipping-service:bootRun --args='--spring.profiles.active=tc'
