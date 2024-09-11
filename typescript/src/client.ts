import { Connection, Client } from '@temporalio/client';
import { example } from './workflows';

async function run() {
  const connection = await Connection.connect();
  const client = new Client({
    connection,
  });
  try {
    const result = await client.workflow.execute(example, {
      taskQueue: 'interceptors-opentelemetry-example',
      workflowId: 'otel-example-0',
      args: ['Temporal'],
    });
    console.log(result); // Hello, Temporal!
  } finally {
    await client.connection.close();
  }
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
