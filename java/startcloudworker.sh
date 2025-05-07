#!/bin/bash
source ../setcloudenv.sh
./gradlew :core:bootRun --args='--spring.profiles.active=tc'
