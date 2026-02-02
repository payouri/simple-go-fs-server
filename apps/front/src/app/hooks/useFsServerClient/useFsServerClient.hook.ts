import { useEffect, useRef, useState } from 'react';
import { audioTypes, CONFIG } from '../../../config';

export type FileMetadata =
  | {
      isDir: true;
      name: string;
      lastModified: string;
      permissions: `${number}${number}${number}`;
    }
  | {
      isDir: false;
      name: string;
      bytesSize: number;
      lastModified: string;
      mimeType: string;
      permissions: `${number}${number}${number}`;
    };

type GetFilesResponse200 = {
  files: FileMetadata[];
  total: number;
  page: number;
  maxPage: number;
};

export const fsServerClient = (function buildFsServerClient() {
  const { fsServerURI, apiKey } = CONFIG;
  if (!fsServerURI || !apiKey) {
    throw new Error('Missing config');
  }
  function sanitizePath(path: string) {
    return path.replace(/\/+/g, '/');
  }

  function buildGetFilesRequest() {
    const mountPoint = new URL('list', fsServerURI);
    console.log(mountPoint);
    return async function getFilesRequest(
      path: string,
      paginationOptions?: {
        limit?: number;
        offset?: number;
      }
    ) {
      const { limit = 20, offset = 0 } = paginationOptions || {};
      let response: Response;
      try {
        const sanitizedPath = sanitizePath(`/${mountPoint.pathname}/${path}`);
        response = await fetch(
          `${mountPoint.protocol}//${mountPoint.host}${sanitizedPath}?limit=${limit}&offset=${offset}`,
          {
            headers: {
              contentType: 'application/json',
              Authorization: apiKey,
            },
          }
        );
      } catch (error) {
        console.error(error);
        return null;
      }
      if (!response.ok) {
        console.error(response);
        return null;
      }
      if (response.status > 299) {
        console.error(response);
        return null;
      }
      if (response.headers.get('content-type') !== 'application/json') {
        console.error(response);
        throw new Error('Invalid content type');
      }

      let json: GetFilesResponse200;
      try {
        json = await response.json();
      } catch (error) {
        console.error(error);
        return null;
      }

      return json;
    };
  }

  return {
    getFiles: buildGetFilesRequest(),
  };
})();

export function useFetch(url: string) {
  const request = useRef<Promise<Response>>(null);
  const [data, setData] = useState<Response | null>(null);

  useEffect(() => {
    const audio = document.createElement('audio');

    for (const type of audioTypes) {
      console.log(type, audio.canPlayType(type));
    }
    if (!request.current) {
      request.current = fetch(new URL(url), {
        headers: {
          Authorization: `${CONFIG.apiKey}`,
        },
      });
      request.current.then((response) => {
        setData(response);
      });
    }
  }, [url]);
  return data;
}

export function useFsServerClient() {
  const { current } = useRef(fsServerClient);

  return current;
}
