#!/bin/bash
source ../setcloudenv.sh donald-demo
env | grep TEMP
exit 1
./gradlew bootRun --args='--spring.profiles.active=tc'
