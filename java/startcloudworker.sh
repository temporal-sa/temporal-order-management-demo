#!/bin/bash
source ../setcloudenv.sh
./gradlew bootRun --args='--spring.profiles.active=tc'
