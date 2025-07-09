# Temporal Order Management Demo - Ruby

This is an alternate implementation of the Temporal Order Management Demo backend
using the [Ruby SDK](https://github.com/temporalio/sdk-ruby).

All of the scenarios supported in the Go backend, and outlined in the main README are implemented in
this Ruby version. This version is also fully compatible with the UI. See the main README for
instructions on how to run the UI, and the instructions below for running the Ruby backend.

Caveats: This SDK is still Alpha, and not all features exist. Of note, only one of update_order (signal) or update_orders_update(update) will work at a time, as they both cannot have the 'UpdateOrder' workflow_signal or update_signal annotation
Also, the validation is not working for the update, even though it's implemented per the docs. `workflow_update_validator(:update_order_update)`
## Run Worker

```bash
./startlocalworker.sh
```
