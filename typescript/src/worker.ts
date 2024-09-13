import { NativeConnection, Runtime, Worker } from '@temporalio/worker';
import * as activities from './activities';
import { getWorkflowOptions, getConnectionOptions, getTelemetryOptions, namespace, taskQueue } from './env';

async function main() {
  const telemetryOptions = getTelemetryOptions();

  if(telemetryOptions) {
    Runtime.install(telemetryOptions);
  }

  const connectionOptions = await getConnectionOptions();
  const connection = await NativeConnection.connect(connectionOptions);

  const worker = await Worker.create({
    connection,
    namespace,
    taskQueue,
    activities: {...activities},
    ...getWorkflowOptions(),
  });
  try {
    console.info('🤖: Temporal Worker Online! Beep Boop Beep!');
    await worker.run();
  } finally {
    console.info('🤖: Temporal Worker Shutdown! Beep Boop Beep!');
    // await worker.shutdown();
  }
}

main().then(
  () => void process.exit(0),
  (err) => {
    console.error(err);
    process.exit(1);
  }
);
