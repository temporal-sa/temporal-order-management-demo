# Temporal Order Management Demo - Java

This is an alternate implementation of the Temporal Order Management Demo backend
using the [Java SDK](https://github.com/temporalio/sdk-java).

All of the scenarios suppported in the Go backend, and outlined in the main README are implemented in
this Java version. This version is also fully compatible with the Python UI. See the main README for
instructions on how to run the UI, and the instructions below for running the Java backend.

## Run Worker

```bash
cd java
./gradlew bootRun
```

## Run Worker with Profile (If you want to use Temporal Cloud)

```bash
cd java
./gradlew bootRun --args='--spring.profiles.active=tc'
```


For the Nexus shipping workflow solution there is the need to setup the Nexus configuration for
either self-hosted or Temporal Cloud.  The instructions for this are held in the main repo [README](../README.md).  There are a number of additional environment variables that need to be set to specify the Nexus endpoint and the worker connection information to the "Shipping" namespace.  These additional environment variables are:


```
# *************************************************************************************************
# Additional env vars for using nexus to demonstrate the shipping fulfillment in another namespace.
# Notes:
# - For self hosted this will default to using the default namespace
# - Nexus only allows unique endpoints in an account so you may need to change this to a 
#   Unique value to match the endpoint you created.  All code references are from this env var.
# *************************************************************************************************
export TEMPORAL_SHIPPING_NAMESPACE=<Your SHIPPING Namespace>.<Account ID>
export TEMPORAL_SHIPPING_TASK_QUEUE=shipping
export TEMPORAL_SHIPPING_ADDRESS=<Namespace>.<account id>.tmprl.cloud:7233
export TEMPORAL_SHIPPING_CERT_PATH=<Path>/<To>/<Client public certificate for Shipping NS>.pem
export TEMPORAL_SHIPPING_KEY_PATH=<Path>/<To>/<Client private key file for Shipping NS>
export TEMPORAL_SHIPPING_CERT_RELOAD_PERIOD=30   # Refresh period for certificaPath>/<To>/<Client ptes (in minutes)
export TEMPORAL_SHIPPING_NEXUS_ENDPOINT=shipping-endpoint
```


To try and simplify the setup there is an env-sample.sh file that details all the environment variables requred to run the Order Management Demo.  Edit this file to provide the details you plan to use and then source the file prior to running the ui or worker.

