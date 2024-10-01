import fs from 'fs/promises';
import type * as client from '@temporalio/client';
import type { RuntimeOptions, WorkerOptions } from '@temporalio/worker';
import { NativeConnectionOptions } from '@temporalio/worker';

// Common set of connection options that can be used for both the client and worker connections.
export type ConnectionOptions = Pick<NativeConnectionOptions, 'tls' | 'address' | 'apiKey' | 'metadata'>;

export function getenv(key: string, defaultValue?: string): string {
  const value = process.env[key];
  if (!value) {
    if (defaultValue != null) {
      return defaultValue;
    }
    throw new Error(`missing env var: ${key}`);
  }
  return value;
}

export async function getConnectionOptions(): Promise<ConnectionOptions> {
  const address = getenv('TEMPORAL_HOST_URL', 'localhost:7233');

  let tls: ConnectionOptions['tls'] = undefined;
  let apiKey: string | undefined = undefined;
  let metadata: Record<string, string> = {};

  if (process.env.TEMPORAL_MTLS_TLS_CERT && process.env.TEMPORAL_MTLS_TLS_KEY) {
    const crt = await fs.readFile(getenv('TEMPORAL_MTLS_TLS_CERT'));
    const key = await fs.readFile(getenv('TEMPORAL_MTLS_TLS_KEY'));

    tls = { clientCertPair: { crt, key } };
    console.info('🤖: Connecting to Temporal Cloud (mTLS) ⛅');
  } else if (process.env.TEMPORAL_APIKEY) {
    apiKey = getenv('TEMPORAL_APIKEY');
    tls = true;
    metadata = {
      'temporal-namespace': getenv('TEMPORAL_NAMESPACE'),
    }

    console.info('🤖: Connecting to Temporal Cloud (API key) ⛅');
  } else {
    console.info('🤖: Connecting to Local Temporal');
  }

  return {
    address,
    tls,
    apiKey,
    metadata
  };
}

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
  let telemetryOptions = {};

  switch (metrics) {
    case 'PROMETHEUS':
      const bindAddress = getenv('TEMPORAL_METRICS_PROMETHEUS_ADDRESS', `0.0.0.0:${port}`);
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

export const namespace = getenv('TEMPORAL_NAMESPACE', 'default');
export const taskQueue = getenv('TEMPORAL_TASK_QUEUE', 'orders');
export const env = getenv('NODE_ENV', 'development');
