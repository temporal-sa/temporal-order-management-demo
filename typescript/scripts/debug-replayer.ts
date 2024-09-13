import { startDebugReplayer } from '@temporalio/worker';

startDebugReplayer({
  workflowsPath: require.resolve('../src/workflows/index'),
});