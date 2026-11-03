# Temporal Order Management Demo - Ruby

An implementation of the Temporal Order Management Demo backend
using the [Ruby SDK](https://github.com/temporalio/sdk-ruby).

All of the scenarios outlined in the main [README](../README.md) are implemented in this Ruby version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## Updating the Temporal SDK

The worker scripts run `bundle install` before starting the Worker, so pulling the latest changes is usually all you
need. If you start the Worker another way, run `bundle install` yourself after pulling.

## NexusOperation Scenario

The Ruby Worker implements the caller side of the NexusOperation scenario only. The `shipping-service` Nexus Service it
calls is handled by the Go Nexus Worker, so run both:

```bash
# create the nexus-demo Namespace and shipping-endpoint first (see the main README)
cd go && ./startlocalworker_nexus.sh
cd ruby && ./startlocalworker.sh
```

The Ruby SDK's Nexus support is experimental. The Workflow sends the Nexus input as a plain hash in the JSON shape the
Go handler expects, rather than as a Struct, which would serialize in Ruby's `json_class` format.
