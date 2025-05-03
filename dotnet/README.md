# temporal-order-management-dotnet

An order management demo written using the Temporal .NET SDK, which is compatable with the Python UI.

See the main [README](../README.md) for instructions for starting the UI

## Run Worker Locally
```bash
./startlocalworker.sh
```

## Start Worker on Temporal Cloud
If you haven't created the setcloudenv.sh file, then copy setcloundenv.example to setcloudenv.sh
and edit as needed

```bash
cd ..
cp setcloudenv.example setcloudenv.sh
vi setcloudenv.sh # edit setcloudenv.sh to match your Temporal Cloud account
```

```bash
# run the worker
./startcloudworker.sh
```
