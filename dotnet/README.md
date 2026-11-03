# Temporal Order Management Demo - .NET

An implementation of the Temporal Order Management Demo backend
using the [.NET SDK](https://github.com/temporalio/sdk-dotnet).

All of the scenarios outlined in the main [README](../README.md) are implemented in this .NET version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## NexusOperation Scenario

The .NET Worker implements the caller side of the NexusOperation scenario only. The `shipping-service` Nexus Service it
calls is handled by the Go Nexus Worker, so run both:

```bash
# create the nexus-demo Namespace and shipping-endpoint first (see the main README)
cd go && ./startlocalworker_nexus.sh
cd dotnet && ./startlocalworker.sh
```
