import type { RuntimeOptions, WorkerOptions } from '@temporalio/worker';

export function getWorkflowOptions(): Pick<WorkerOptions, 'workflowBundle' | 'workflowsPath'> {
  const workflowBundlePath = getenv('WORKFLOW_BUNDLE_PATH', 'lib/workflow-bundle.js');

  if (workflowBundlePath && env == 'production') {
    return { workflowBundle: { codePath: workflowBundlePath } };
  } else {
    return { workflowsPath: require.resolve('./workflows/index') };
  }
}

export function getTelemetryOptions(): RuntimeOptions {
  const metrics = getenv('TEMPORAL_WORKER_METRIC', 'PROMETHEUS');
  const port = getenv('TEMPORAL_WORKER_METRICS_PORT', '9464');
  let bindAddress: string;
  let telemetryOptions = {};

  switch (metrics) {
    case 'PROMETHEUS':
      bindAddress = getenv('TEMPORAL_METRICS_PROMETHEUS_ADDRESS', `0.0.0.0:${port}`);
      telemetryOptions = {
        metrics: {
          prometheus: {
            bindAddress,
          },
        },
      };
      console.info('🤖: Prometheus Metrics 🔥', bindAddress);
      break;
    case 'OTEL':
      telemetryOptions = {
        metrics: {
          otel: {
            url: getenv('TEMPORAL_METRICS_OTEL_URL'),
            headers: {
              'api-key': getenv('TEMPORAL_METRICS_OTEL_API_KEY'),
            },
          },
        },
      };
      console.info('🤖: OTEL Metrics 📈');
      break;
    default:
      console.info('🤖: No Metrics');
      break;
  }

  return { telemetryOptions };
}

export const taskQueue = getenv('TEMPORAL_TASK_QUEUE', 'orders');
export const env = getenv('NODE_ENV', 'development');

function getenv(key: string, defaultValue?: string): string {
  const value = process.env[key];
  if (!value) {
    if (defaultValue != null) {
      return defaultValue;
    }
    throw new Error(`missing env var: ${key}`);
  }
  return value;
}
