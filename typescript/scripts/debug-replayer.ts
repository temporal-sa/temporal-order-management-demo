import { startDebugReplayer } from '@temporalio/worker';
import { getDataConverter } from '../src/data-converter';

startDebugReplayer({
  workflowsPath: require.resolve('../src/workflows/index'),
  dataConverter: getDataConverter(),
});
