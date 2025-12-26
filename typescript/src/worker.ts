import { NativeConnection, Runtime, Worker } from '@temporalio/worker';
import * as activities from './activities/index';
import { getWorkflowOptions, getTelemetryOptions, taskQueue } from './env';
import { loadClientConnectConfig } from '@temporalio/envconfig';
// import { createApiKeyServer } from './apikey-server';

async function main() {
  const telemetryOptions = getTelemetryOptions();
  if (telemetryOptions) {
    Runtime.install(telemetryOptions);
  }

  const config = loadClientConnectConfig();
  const connection = await NativeConnection.connect(config.connectionOptions);
  console.info(`✅ Client connected to ${config.connectionOptions.address} in namespace '${config.namespace}'`);

  // if (process.env.TEMPORAL_APIKEY) {
  //   createApiKeyServer(connection).listen(3333, () => {
  //     console.log('API Key server is running on http://localhost:3333');
  //   });
  // }

  try {
    const worker = await Worker.create({
      connection,
      namespace: config.namespace,
      taskQueue,
      activities: { ...activities },
      ...getWorkflowOptions(),
    });

    console.info('🤖: Temporal Worker Online! Beep Boop Beep!');
    await worker.run();
  } finally {
    await connection.close();
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
