import { Worker } from '@temporalio/worker';
import * as activities from './activities';

async function main() {
  const worker = await Worker.create({
    workflowsPath: require.resolve('./workflows'),
    activities,
    taskQueue: 'interceptors-opentelemetry-example',
  });
  try {
    await worker.run();
  } finally {
    await worker.shutdown();
  }
}

main().then(
  () => void process.exit(0),
  (err) => {
    console.error(err);
    process.exit(1);
  }
);
