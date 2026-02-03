import { CONFIG } from '../../../config';

export const streamServerClient = (function buildFsServerClient() {
  const { streamServerURI, apiKey } = CONFIG;
  if (!streamServerURI || !apiKey) {
    throw new Error('Missing config');
  }
  function sanitizePath(path: string) {
    return path.replace(/\/+/g, '/');
  }
  function getFileUri(path: string) {
    console.log({
      path,
    });
    return new URL(sanitizePath(['stream', path].join('/')), streamServerURI)
      .href;
  }

  function buildStreamFile() {
    const mountPoint = new URL('stream', streamServerURI);
    return async function streamFileRequest(path: string) {
      console.log('streamFileRequest', { mountPoint, path });
    };
  }

  return {
    streamFile: buildStreamFile(),
    getFileUri,
  };
})();
