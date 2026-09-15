import { S3Client } from '@aws-sdk/client-s3';
import { ExternalStorage } from '@temporalio/common';
import type { DataConverter } from '@temporalio/common';
import { S3StorageDriver } from '@temporalio/external-storage-s3';
import { AwsSdkS3StorageDriverClient } from '@temporalio/external-storage-s3-aws-sdk';
import { awsRegion, externalStorage, s3Bucket, s3Endpoint } from './env';

const EXTERNAL_STORAGE_WORKFLOW_TYPE = 'OrderWorkflowExternalStorage';

export function getDataConverter(): DataConverter | undefined {
  if (!externalStorage) {
    console.info('🤖: No External Storage');
    return undefined;
  }

  // forcePathStyle is required for MinIO
  const client = new S3Client({
    region: awsRegion,
    ...(s3Endpoint ? { endpoint: s3Endpoint, forcePathStyle: true } : {}),
  });

  const driver = new S3StorageDriver({
    client: new AwsSdkS3StorageDriverClient(client),
    bucket: s3Bucket,
  });

  console.info('🤖: External Storage 📦', s3Bucket);

  return {
    externalStorage: new ExternalStorage({
      drivers: [driver],
      payloadSizeThreshold: 0,
      // this worker serves every scenario - only offload the one that demos claim checks
      driverSelector: (context) => (context.target?.type == EXTERNAL_STORAGE_WORKFLOW_TYPE ? driver : null),
    }),
  };
}
