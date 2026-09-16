# Temporal Order Management Demo - Typescript

An implementation of the Temporal Order Management Demo backend
using the [Typescript SDK](https://github.com/temporalio/sdk-typescript)

All of the scenarios outlined in the main [README](../README.md) are implemented in this Typescript version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## External Storage

The ExternalStorage scenario offloads payloads to S3 and writes claim checks into Event History. See the main
[README](../README.md#external-storage) for what it demonstrates and how to run it, and
[`../setcloudenv.example`](../setcloudenv.example) for Temporal Cloud. Defaults live in [src/env.ts](src/env.ts).

When `S3_ENDPOINT` is set the code supplies MinIO's region and credentials itself, because an ambient `AWS_PROFILE`
takes precedence over `AWS_ACCESS_KEY_ID` and would otherwise sign the demo's requests with real AWS credentials.

The debug replayer uses the same data converter, so replaying a history with references needs external storage enabled
and access to the bucket.

> [!CAUTION]
> The codec server has no authentication - anything that can reach it can read any payload in the bucket. Keep it on
> localhost. See [Securing a codec server](https://docs.temporal.io/codec-server#securing-a-codec-server).

> [!NOTE]
> External storage is experimental, so the `@temporalio/*` packages are pinned to exactly `1.23.0` and must move
> together - `external-storage-s3` exact-depends on `common`, so a caret range resolves a second copy of
> `ExternalStorage`. The Web UI must read the format `1.23.0` writes, hence `temporalio[aioboto3]>=1.32.0`.

## (Optional) Run Worker in Productionize Build

1. `npm run build` to build out the worker and activites.
1. `NODE_ENV=production node lib/worker.js` to run the production Worker.

## (Optional) Using VSCode Debugger

1. Install the [VSCode Debugger](https://temporal.io/blog/temporal-for-vs-code)
1. `cd typescript/`
1. `Command + Shift + P` Select `Temporal: Open Panel`

Just be mindful, that the VSCode Debugger doesn't support the `OrderWorkflowHumanInLoopUpdate` use case.
