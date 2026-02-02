import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { FileMetadata } from '../../hooks/useFsServerClient/useFsServerClient.hook';
import { DirectoryBreadcrumb } from '../DirectoryBreadcrumb/DirectoryBreadcrumb.component';
import { splitPath } from '../../../helpers/splitPath';

export type FileListProps = {
  path: string;
  files: Array<FileMetadata>;
  currentPage: number;
  totalPages: number;
  totalFiles: number;
  isLoadingFiles: boolean;
  onFileClick: (file: FileMetadata) => void;
  onPageChange: (page: number) => void;
};

type LoadingStateProps = {
  isLoadingFiles: boolean;
};
function LoadingState(props: LoadingStateProps) {
  if (!props.isLoadingFiles) {
    return null;
  }

  return (
    <TableBody
      className="position-absolute h-full w-full h-full overflow-hidden"
      style={{
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        zIndex: 10,
      }}
    >
      <tr>
        <div
          className="flex items-center justify-center top-0 left-0 right-0 bottom-0"
          style={{
            position: 'absolute',
          }}
        >
          Loading...
        </div>
      </tr>
    </TableBody>
  );
}
export function FileList(props: FileListProps) {
  const {
    path,
    files,
    currentPage,
    totalPages,
    totalFiles,
    isLoadingFiles,
    onFileClick,
    onPageChange,
  } = props;

  function handlePageChange(page: number | 'prev' | 'next') {
    if (page === 'prev') {
      onPageChange(Math.max(currentPage - 1, 1));
    } else if (page === 'next') {
      onPageChange(Math.min(currentPage + 1, totalPages));
    } else {
      onPageChange(Math.min(Math.max(page, 1), totalPages));
    }
  }

  return (
    <div className="position-relative h-full flex flex-col">
      <DirectoryBreadcrumb path={splitPath(path)} />
      {/* <div className="flex flex-col"> */}
      <Table className="position-relative">
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
          </TableRow>
        </TableHeader>
        <LoadingState isLoadingFiles={isLoadingFiles} />
        <TableBody>
          {files.map((file, index) => (
            <TableRow key={index} onClick={() => onFileClick(file)}>
              <TableCell>{file.name}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {/* </div> */}
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              disabled={currentPage === 1}
              onClick={() => handlePageChange('prev')}
            />
          </PaginationItem>
          {Array.from({ length: totalPages }, (_, index) => (
            <PaginationItem key={index}>
              <PaginationLink
                isActive={currentPage === index + 1}
                onClick={() => handlePageChange(index + 1)}
              >
                {index + 1}
              </PaginationLink>
            </PaginationItem>
          ))}
          <PaginationItem>
            <PaginationNext
              disabled={currentPage === totalPages}
              onClick={() => handlePageChange('next')}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  );
}
