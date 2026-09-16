import { S3Client } from '@aws-sdk/client-s3';
import { ExternalStorage } from '@temporalio/common';
import type { DataConverter } from '@temporalio/common';
import { S3StorageDriver } from '@temporalio/external-storage-s3';
import { AwsSdkS3StorageDriverClient } from '@temporalio/external-storage-s3-aws-sdk';
import { awsRegion, externalStorage, s3Bucket, s3Endpoint } from './env';

const EXTERNAL_STORAGE_WORKFLOW_TYPE = 'OrderWorkflowExternalStorage';

// MinIO is addressed by path rather than by virtual host, and its region and credentials are fixed. They are set
// here rather than exported by the start scripts because an ambient AWS_PROFILE takes precedence over
// AWS_ACCESS_KEY_ID, which would otherwise sign the local demo's requests with real AWS credentials and 403.
const MINIO = {
  forcePathStyle: true,
  region: 'us-east-1',
  credentials: { accessKeyId: 'minioadmin', secretAccessKey: 'minioadmin' },
};

export function getDataConverter(): DataConverter | undefined {
  if (!externalStorage) {
    console.info('🤖: No External Storage');
    return undefined;
  }

  const client = new S3Client({
    region: awsRegion,
    ...(s3Endpoint ? { endpoint: s3Endpoint, ...MINIO } : {}),
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
