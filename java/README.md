# Temporal Order Management Demo - Java

This is an alternate implementation of the Temporal Order Management Demo backend using
the [Java SDK](https://github.com/temporalio/sdk-java).

All the scenarios supported in the Go backend, and outlined in the main README are implemented in this Java version.
This version is also fully compatible with the Python UI. See the main README for instructions on how to run the UI, and
the instructions below for running the Java backend.

NOTE - If the nexus shipping service is to be used then it is mandatory to ensure that the environment variable to set
the endpoint is set.

```
$ export TEMPORAL_NEXUS_SHIPPING_ENDPOINT=<value used for the end point in setting up the service.>
```

# Run Workers for Local

To start the services for running against a local Temporal Service use

`./startlocalworker.sh`
&
`./startlocalworker_nexus.sh`

# Run Workers for Cloud

To start the services for running against Temporal Cloud use

`./startcloudworker.sh`
&
`./startcloudworker_nexus.sh`
