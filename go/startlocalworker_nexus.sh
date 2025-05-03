#!/bin/bash

# make sure you have created your endoint, e.g.
# temporal operator nexus endpoint create --name shipping-endpoint --target-namespace default --target-task-queue shipping

go run nexus/worker/worker.go
