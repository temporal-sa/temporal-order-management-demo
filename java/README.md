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

## Run Worker with Profile

```bash
cd java
./gradlew bootRun --args='--spring.profiles.active=tc'
```


Adding Nexus support for Java
