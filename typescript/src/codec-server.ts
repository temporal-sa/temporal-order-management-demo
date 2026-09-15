import type { Payload } from '@temporalio/common';
import { ExternalStorageRunner, isReferencePayload } from '@temporalio/common/lib/internal-non-workflow';
import { temporal } from '@temporalio/proto';
import { createServer, IncomingMessage, ServerResponse } from 'http';
import { getDataConverter } from './data-converter';
import { webUiOrigin } from './env';

const { Payloads } = temporal.api.common.v1;
const port = 8081;

const externalStorage = getDataConverter()?.externalStorage;
if (!externalStorage) {
  console.error('🤖: Codec Server needs EXTERNAL_STORAGE=true');
  process.exit(1);
}

const runner = new ExternalStorageRunner(externalStorage);

async function getRequestBody(req: IncomingMessage): Promise<string> {
  const chunks: Buffer[] = [];
  for await (const chunk of req) {
    chunks.push(chunk as Buffer);
  }
  return Buffer.concat(chunks).toString();
}

const requestHandler = async (req: IncomingMessage, res: ServerResponse): Promise<void> => {
  // CORS
  res.setHeader('Access-Control-Allow-Origin', webUiOrigin);
  res.setHeader('Access-Control-Allow-Methods', 'POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'content-type, x-namespace');

  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method !== 'POST') {
    res.writeHead(404);
    res.end();
    return;
  }

  const url = new URL(req.url ?? '/', `http://localhost:${port}`);

  try {
    const { payloads } = Payloads.fromObject(JSON.parse(await getRequestBody(req)));
    let result: Payload[];

    switch (url.pathname) {
      case '/decode':
        result = url.searchParams.get('preserveStorageRefs') == 'true' ? payloads : await runner.retrieve(payloads);
        break;
      case '/download':
        if (!payloads.every(isReferencePayload)) {
          res.writeHead(400);
          res.end();
          return;
        }
        result = await runner.retrieve(payloads);
        break;
      default:
        res.writeHead(404);
        res.end();
        return;
    }

    res.writeHead(200, { 'content-type': 'application/json' });
    res.end(
      JSON.stringify(
        result.length
          ? Payloads.toObject(Payloads.create({ payloads: result }), { bytes: String, longs: String })
          : { payloads: [] },
      ),
    );
  } catch (error) {
    console.error(error);
    res.writeHead(500);
    res.end();
  }
};

createServer(requestHandler).listen(port, () => {
  console.info(`🤖: Codec Server on http://localhost:${port}`);
});
