import { useEffect, useRef, useState } from 'react';
const CONFIG = {
  fsServerURI: 'http://localhost:5008',
  streamServerURI: 'http://localhost:5001',
  apiKey: import.meta.env.VITE_FS_SERVER_API_KEY,
} as const;

const audioTypes = [
  'audio/mpeg',
  'audio/mp4',
  'audio/ogg',
  'audio/wav',
  'audio/aac',
  'audio/webm',
];

const fsServerClient = (function buildFsServerClient() {
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
    return async function getFilesRequest(path: string) {
      let response: Response;
      try {
        const sanitizedPath = sanitizePath(`/${mountPoint.pathname}/${path}`);
        response = await fetch(
          `${mountPoint.protocol}//${mountPoint.host}${sanitizedPath}`,
          {
            headers: {
              contentType: 'application/json',
              Authorization: `${apiKey}`,
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

      let json: any;
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

function useFetch(url: string) {
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

function useFsServerClient() {
  const { current } = useRef(fsServerClient);

  return current;
}

export function App() {
  useFetch(new URL('/stream', CONFIG.streamServerURI).toString());
  const fsServerClient = useFsServerClient();

  console.log(fsServerClient.getFiles('/Downloads'));
  return <div>Hello World</div>;
}

export default App;
