import { NativeConnection } from '@temporalio/worker';
import { createServer, IncomingMessage, Server, ServerResponse } from 'http';

const getRequestBody = (req: IncomingMessage): Promise<string> => {
  return new Promise((resolve, reject) => {
    let body: string = '';

    req.on('data', (chunk) => {
      body += chunk.toString();
    });

    req.on('end', () => {
      try {
        resolve(body);
      } catch (error) {
        reject(error);
      }
    });

    req.on('error', (error) => {
      reject(error);
    });
  });
};

export function createApiKeyServer (connection: NativeConnection): Server<any, any> {
    const requestHandler = async (req: IncomingMessage, res: ServerResponse): Promise<any> => {
        // CORS
        res.setHeader('Access-Control-Allow-Origin', '*');
        res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
        res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
      
        if (req.method === 'OPTIONS') {
          res.writeHead(204);
          res.end();
          return;
        }

        if (req.method === 'PUT') {
          try {
            const newApiKey: string = await getRequestBody(req);
            res.statusCode = 202;
            connection.setApiKey(newApiKey);
            res.end();
          } catch (error) {
            res.statusCode = 400;
            res.end();
          }
        } else {
          res.statusCode = 404;
          res.end();
        }
      };

    return createServer(requestHandler);
}