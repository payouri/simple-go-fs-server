import { useQuery } from '@tanstack/react-query';

const BASE_QUERY_KEY = 'file' as const;

export type UseFileParams = {
  uri: string;
};

export type UseFileReturn = Readonly<[null | undefined, boolean]>;

export function useFile(params: UseFileParams): UseFileReturn {
  const { uri } = params;

  const useFileQuery = useQuery({
    queryKey: [BASE_QUERY_KEY, uri],
    queryFn() {
      return null;
    },
  });

  return [useFileQuery.data, useFileQuery.isLoading] as const;
}
