import { useLocation, useNavigate } from 'react-router-dom';
import { FileList } from './components/FileList/FileList.component';
import { useFiles } from './hooks/useFiles/useFiles.hook';
import { splitPath } from '../helpers/splitPath';
import type { FileMetadata } from '@/entities/File.types';

export function App() {
  const location = useLocation();
  const navigate = useNavigate();
  const { pathname: path } = location;
  const { getFiles, metadata, pagination, setFilesPagination } = useFiles({
    path,
    initialLimit: 20,
    initialOffset: 0,
  });

  function onPageChange(page: number) {
    const newOffset = (page - 1) * pagination.limit;
    if (newOffset !== pagination.offset) {
      setFilesPagination({ offset: newOffset });
    }
  }

  function onFileClick(file: FileMetadata) {
    if (file.isDir) {
      navigate(
        {
          pathname:
            path === '/'
              ? file.name
              : [...splitPath(path), file.name].join('/'),
        },
        {
          relative: 'route',
        }
      );
      return;
    }
  }

  const { data: getFilesData, isLoading } = getFiles;
  const { files = [] } = getFilesData || {};
  const total = getFilesData?.total ?? metadata?.total ?? 0;
  const page = getFilesData?.page ?? metadata?.page ?? 0;
  const maxPage = getFilesData?.maxPage ?? metadata?.maxPage ?? 0;

  return (
    <div className="flex flex-col h-screen">
      <div className="flex flex-col flex-[1_1_auto] px-4 py-2">
        <FileList
          files={files}
          totalFiles={total}
          totalPages={maxPage}
          currentPage={page}
          isLoadingFiles={isLoading}
          path={decodeURIComponent(path)}
          onFileClick={onFileClick}
          onPageChange={onPageChange}
        />
      </div>
      <div className="flex flex-col flex-0_0_auto"></div>
    </div>
  );
}

export default App;
