# Temporal Order Management Demo - Typescript

An implementation of the Temporal Order Management Demo backend
using the [Typescript SDK](https://github.com/temporalio/sdk-typescript)

All of the scenarios outlined in the main [README](../README.md) are implemented in this Typescript version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## External Storage

The ExternalStorage scenario offloads payloads to S3 and writes claim checks into Event History. See the main
[README](../README.md#external-storage) for what the scenario demonstrates and what to look for in the Web UI.

MinIO stands in for S3 locally and runs under Docker Compose from the root of the repo, since it is not specific to
this SDK.

1. `../startminio.sh` to start MinIO and create the bucket
1. `./startlocalworker.sh` to run the Worker with external storage enabled
1. (Optional) `./startcodecserver.sh` to run the codec server on http://localhost:8081

The Worker, the codec server, and the Python Web UI all read the same environment variables:

| Variable           | Default             | Description                                                                  |
| :----------------- | :------------------ | :--------------------------------------------------------------------------- |
| `EXTERNAL_STORAGE` | `false`             | `true` enables external storage, any other value turns it off entirely       |
| `S3_ENDPOINT`      | unset               | S3 endpoint, e.g. `http://localhost:9000` for MinIO. Unset means real AWS S3 |
| `S3_BUCKET`        | `temporal-payloads` | Bucket the payloads are written to                                           |
| `AWS_REGION`       | `us-east-1`         | Region used by the S3 client                                                 |

The codec server reads one more of its own:

| Variable        | Default                 | Description                                                 |
| :-------------- | :---------------------- | :---------------------------------------------------------- |
| `WEB_UI_ORIGIN` | `http://localhost:8233` | Origin allowed by CORS - the Web UI the codec server serves |

For Temporal Cloud, `./startcloudcodecserver.sh` sources `../setcloudenv.sh` the same way the Worker and Web UI cloud
scripts do, so set `WEB_UI_ORIGIN=https://cloud.temporal.io` there along with the bucket and region.

> [!CAUTION]
> The codec server has no authentication, so anything that can reach it can read any payload in the bucket. That is
> fine on localhost and not fine exposed. Before putting it anywhere reachable, terminate TLS in front of it - the
> Cloud Web UI is served over https and browsers block plaintext subresource requests from a secure page, with
> `http://localhost` the only exemption - and verify the access token the Web UI can forward when
> "Pass access token" is enabled.

`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are picked up from the environment by the AWS SDK, and by aioboto3 in
the Web UI. The local scripts set them to the MinIO defaults.

> [!NOTE]
> External storage is experimental and breaking changes are already queued upstream, so the `@temporalio/*` packages
> are pinned to exactly `1.23.0` rather than a caret range. The reference format written by `1.23.0` needs
> `temporalio>=1.27.0` on the Python side to be resolved, which is why `ui/pyproject.toml` requires
> `temporalio[aioboto3]>=1.32.0`.

`protobufjs` is a direct dependency of this package so that every `@temporalio/*` package resolves to one hoisted copy.
Without it those packages each nest their own, `instanceof` fails across the package boundary, and the first offloaded
payload throws `TypeError: type must be a Type`. A separate older copy under `@grpc/proto-loader` is expected and
harmless - it is not on the payload path.

The debug replayer uses the same data converter, so replaying a history that contains references needs external storage
enabled and live access to the bucket.

## (Optional) Run Worker in Productionize Build

1. `npm run build` to build out the worker and activites.
1. `NODE_ENV=production node lib/worker.js` to run the production Worker.

## (Optional) Using VSCode Debugger

1. Install the [VSCode Debugger](https://temporal.io/blog/temporal-for-vs-code)
1. `cd typescript/`
1. `Command + Shift + P` Select `Temporal: Open Panel`

Just be mindful, that the VSCode Debugger doesn't support the `OrderWorkflowHumanInLoopUpdate` use case.
