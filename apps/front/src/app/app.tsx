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

function useFetch(url: string) {
  const request = useRef<Promise<Response>>(null);
  const [data, setData] = useState<Response | null>(null);

  useEffect(() => {
    const audio = document.createElement('audio');

    for (const type of audioTypes) {
      console.log(type, audio.canPlayType(type));
    }
    if (!request.current) {
      request.current = fetch(url, {
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

export function App() {
  useFetch(new URL('/stream', CONFIG.streamServerURI).toString());
  return <div>Hello World</div>;
}

export default App;
