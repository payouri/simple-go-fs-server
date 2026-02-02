export const CONFIG = {
  fsServerURI: 'http://localhost:5008',
  streamServerURI: 'http://localhost:5001',
  apiKey: import.meta.env.VITE_FS_SERVER_API_KEY,
} as const;

export const audioTypes = [
  'audio/mpeg',
  'audio/mp4',
  'audio/ogg',
  'audio/wav',
  'audio/aac',
  'audio/webm',
];
