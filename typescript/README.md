# Temporal Order Management Demo - Typescript

An implementation of the Temporal Order Management Demo backend
using the [Typescript SDK](https://github.com/temporalio/sdk-typescript)

All of the scenarios outlined in the main [README](../README.md) are implemented in this Typescript version, except where noted.
See the main README for instructions on how to run the UI, and the Workers.

## (Optional) Run Worker in Productionize Build

1. `npm run build` to build out the worker and activites.
1. `NODE_ENV=production node lib/worker.js` to run the production Worker.

## (Optional) Using VSCode Debugger

1. Install the [VSCode Debugger](https://temporal.io/blog/temporal-for-vs-code)
1. `cd typescript/`
1. `Command + Shift + P` Select `Temporal: Open Panel`

Just be mindful, that the VSCode Debugger doesn't support the `OrderWorkflowHumanInLoopUpdate` use case.
