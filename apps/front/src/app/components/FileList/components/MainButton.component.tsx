import { useLayoutEffect, useReducer, useRef } from 'react';
import { Button } from '../../../../components/ui/button';
import type { FileMetadata } from '../../../../entities/File.types';
import { streamServerClient } from '../../../hooks/useStreamServerClient/useStreamServerClient.hook';

export type MainButtonProps = {
  file: FileMetadata;
  path: string;
};

export type AudioFileButtonProps = Omit<MainButtonProps, 'file'> & {
  file: Extract<FileMetadata, { isDir: false }>;
};

function AudioFileButton(props: AudioFileButtonProps) {
  const { file, path } = props;
  const audioRef = useRef<HTMLAudioElement>(null);
  useLayoutEffect(() => {
    if (!audioRef.current) {
      return;
    }

    audioRef.current.addEventListener(
      'error',
      (e) => {
        console.log('error', e, e?.currentTarget?.error);
      },
      {
        once: true,
      }
    );
  }, []);

  const filePath = streamServerClient.getFileUri([path, file.name].join('/'));
  const [state, dispatch] = useReducer<
    {
      isPlaying: boolean;
    },
    [
      {
        type: 'toggleIsPlaying';
        data?: boolean;
      }
    ]
  >(
    (state, action) => {
      if (action.type === 'toggleIsPlaying') {
        if (!audioRef.current) {
          return state;
        }

        const newIsPlaying = action.data ?? !state.isPlaying;
        if (newIsPlaying === state.isPlaying) {
          return state;
        }
        if (newIsPlaying) {
          audioRef.current?.play().catch(console.error);
        } else {
          audioRef.current?.pause();
        }

        return {
          isPlaying: action.data ?? !state.isPlaying,
        };
      }
      return state;
    },
    {
      isPlaying: false,
    }
  );
  return (
    <div>
      <audio ref={audioRef} preload="auto" controls>
        <source src={filePath} type={'audio/webm'} />
      </audio>
      <a href={filePath} download={file.name}>
        {/* <Button onClick={() => dispatch({ type: 'toggleIsPlaying' })}></Button> */}
        Download
      </a>
      <Button onClick={() => dispatch({ type: 'toggleIsPlaying' })}>
        {state.isPlaying ? 'Pause' : 'Play'}
      </Button>
    </div>
  );
}

function isDirFile(
  file: FileMetadata
): file is Extract<FileMetadata, { isDir: true }> {
  return file.isDir;
}

export function MainButton(props: MainButtonProps) {
  const { file, ...rest } = props;
  if (isDirFile(file)) {
    return null;
  }

  if (file.mimeType.startsWith('audio/')) {
    return <AudioFileButton {...rest} file={file} />;
  }

  return null;
}
