import { useQuery } from '@tanstack/react-query';
import { fsServerClient } from '../useFsServerClient/useFsServerClient.hook';
import { useURLState, type UseURLStateReturn } from '../useURLState';

export type UseFilesParams = {
  path: string;
  initialLimit?: number;
  initialOffset?: number;
};

const GET_FILES_QUERY_KEY = 'fsServer/getFiles' as const;
function formatState<
  Result extends Record<string, unknown>,
  State extends Record<string, string>
>(
  params: UseURLStateReturn<State>,
  formatFn: (state: State) => Result
): [Result, UseURLStateReturn<State>[1]] {
  const [state, setState] = params;

  return [formatFn(state), setState];
}

export function useFiles(params: UseFilesParams) {
  const { path, initialLimit = 20, initialOffset = 0 } = params;
  const [{ limit, offset }, setPagination] = formatState(
    useURLState<{
      limit: string;
      offset: string;
    }>({
      initialState: {
        limit: initialLimit.toString(),
        offset: initialOffset.toString(),
      },
    }),
    (state) => ({
      limit: state.limit ? parseInt(state.limit, 10) : initialLimit,
      offset: state.offset ? parseInt(state.offset, 10) : initialOffset,
    })
  );

  const getFilesQueryResult = useQuery({
    queryKey: [GET_FILES_QUERY_KEY, path, limit, offset],
    async queryFn() {
      const response = await fsServerClient.getFiles(path, {
        limit,
        offset,
      });

      return response;
    },
  });

  return {
    getFiles: getFilesQueryResult,
    setFilesPagination: (params: { limit?: number; offset?: number }) => {
      if (
        typeof params.limit !== 'number' &&
        typeof params.offset !== 'number'
      ) {
        return;
      }
      const { limit, offset } = params;

      setPagination((prev) => ({
        ...prev,
        limit: limit?.toString(10) ?? prev.limit,
        offset: offset?.toString(10) ?? prev.offset,
      }));
    },
    pagination: {
      limit,
      offset,
    },
  };
}
