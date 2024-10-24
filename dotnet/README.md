# temporal-order-management-dotnet

An order management demo  mwritten using the Temporal .NET SDK, which is compatable with the Python UI.

See the main [README](../README.md) for instructions for starting the UI

## Prerequisites
Install the .NET SDK
```bash
brew isntall -cask dotnet
```

## Run Worker Locally
```bash
./startlocalworker.sh
```

## Start Worker on Temporal Cloud

If you haven't created the setcloudenv.sh file, the setcloundenv.example to setcloudenv.sh

```bash
cd ..
cp setcloudenv.example setcloudenv.sh
```
Edit setcloudenv.sh to match your Temporal Cloud account:

```bash
export TEMPORAL_ADDRESS=<namespace>.<accountId>.tmprl.cloud:7233
export TEMPORAL_NAMESPACE=<namespace>.<accountId>
export TEMPORAL_CERT_PATH="/path/to/cert.pem"
export TEMPORAL_KEY_PATH="/path/to/key.key"
export TEMPORAL_TASK_QUEUE=orders
```

Then start the UX. Instructions can be found in the [README](../README.md).

```bash
# run the worker
./startcloudworker.sh
```
