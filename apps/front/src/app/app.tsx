import { useLocation, useNavigate } from 'react-router-dom';
import { FileList } from './components/FileList/FileList.component';
import { useFiles } from './hooks/useFiles/useFiles.hook';
import type { FileMetadata } from './hooks/useFsServerClient/useFsServerClient.hook';
import { splitPath } from '../helpers/splitPath';

export function App() {
  const location = useLocation();
  const navigate = useNavigate();
  const { pathname: path } = location;
  const { getFiles, pagination, setFilesPagination } = useFiles({
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

  const { data: getFilesData } = getFiles;
  const { files = [], total = 0, page = 0, maxPage = 0 } = getFilesData || {};
  return (
    <div className="flex flex-col h-screen px-4 py-2">
      <FileList
        files={files}
        totalFiles={total}
        totalPages={maxPage}
        currentPage={page}
        isLoadingFiles={getFiles.isLoading}
        path={decodeURIComponent(path)}
        onFileClick={onFileClick}
        onPageChange={onPageChange}
      />
    </div>
  );
}

export default App;
