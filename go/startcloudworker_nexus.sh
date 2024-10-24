#!/bin/bash
source ../setcloudenv.sh

# make sure you have created your endoint, e.g.
# tcld nexus endpoint create --name shipping-endpoint --target-task-queue shipping --target-namespace helloworld.sdvdw --allow-namespace helloworld.sdvdw

go run nexus/worker/worker.go
