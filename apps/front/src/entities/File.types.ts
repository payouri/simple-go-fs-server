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
