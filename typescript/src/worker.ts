import { NativeConnection, Runtime, Worker } from '@temporalio/worker';
import * as activities from './activities/index';
import { getWorkflowOptions, getConnectionOptions, getTelemetryOptions, namespace, taskQueue } from './env';

async function main() {
  const telemetryOptions = getTelemetryOptions();

  if (telemetryOptions) {
    Runtime.install(telemetryOptions);
  }

  const connectionOptions = await getConnectionOptions();
  const connection = await NativeConnection.connect(connectionOptions);

  const worker = await Worker.create({
    connection,
    namespace,
    taskQueue,
    activities: { ...activities },
    ...getWorkflowOptions(),
  });

  console.info('🤖: Temporal Worker Online! Beep Boop Beep!');
  await worker.run();
}

main().then(
  () => void process.exit(0),
  (err) => {
    console.error(err);
    process.exit(1);
  },
);
