# Temporal Order Management Demo - Java

This is an alternate implementation of the Temporal Order Management Demo backend
using the [Java SDK](https://github.com/temporalio/sdk-java).

All of the scenarios suppported in the Go backend, and outlined in the main README are implemented in
this Java version. This version is also fully compatible with the Python UI. See the main README for
instructions on how to run the UI, and the instructions below for running the Java backend.

NOTE - If the nexus shipping service is to be used then it is mandatory to ensure that the environment variable to set the endpoint is set.

```
$ export TEMPORAL_SHIPPING_ENDPOINT=<value used for the end point in setting up the service.>
```


(See scripts section for running the services from scripts!)

## Run Worker
```bash
cd java
./gradlew  :order-management:bootRun 
```

## Run Worker with Profile (If you want to use Temporal Cloud)

```bash
cd java
./gradlew :order-management:bootRun --args='--spring.profiles.active=tc'
```




```
# *************************************************************************************************
# Additional env vars for using nexus to demonstrate the shipping fulfillment in another namespace.
# Notes:
# - For self hosted this will default to using the default namespace
# - Nexus only allows unique endpoints in an account so you may need to change this to a 
#   Unique value to match the endpoint you created.  All code references are from this env var.
# *************************************************************************************************
export TEMPORAL_SHIPPING_TASK_QUEUE=shipping
export TEMPORAL_SHIPPING_NEXUS_ENDPOINT=shipping-endpoint
```


To try and simplify the setup there is an env-sample.sh file that details all the environment variables required to run the Order Management Demo.  Edit this file to provide the details you plan to use and then source the file prior to running the ui or worker.


# Scripts for running Workers.
It is possible to set the environment variables and run from the gradle command line however it may be easier to use scripts to start up the services.  For Nexus workloads it is necessary to startup one worker to handle the main/happy path workflows and if the Nexus option is selectted then it is necessary to startup the Nexus worker to progress the shipping workflow.

The Temporal command line can use an ["environment"](https://docs.temporal.io/cli/env/) to connect to the cloud service.   Multiple environments can be configured allowing quick switching to different namespaces.  Setup one environment to connect to the order namespace and another for the shipping namespace.

To start the services for running against a local Temporal Service use

`./startlocalworker.sh`
&
`./startlocalworker_nexus.sh`

To start the services for running against Temporal Cloud use

`./startcloudworker.sh <env name that connects to TCloud> <Shipping endpoint>`

& 

`./startcloudworker_nexus.sh <env name for shipping namespace on T Cloud>` 

Where the first parameter is used to specify the environment name that is used for the temporal CLI to connect to the namespace.  Cloud worker to the orders namespace and the nexus one for the shipping namespace.

The second parameter is required for the order management worker to specify the Nexus endpoint where the shipping nexus calls will be made.




