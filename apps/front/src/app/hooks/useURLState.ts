import type { SetStateAction } from 'react';
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom';

function getQueryFields<Fields extends string[]>(params: {
  query: URLSearchParams;
  fields: Fields;
}): {
  [K in Fields[number]]: string | undefined;
} {
  const { query, fields } = params;

  const result: {
    [K in Fields[number]]: string | undefined;
  } = {} as any;

  for (const field of fields) {
    const value = query.get(field);
    if (value !== null) {
      result[field as Fields[number]] = value;
    }
  }

  return result;
}

function areFieldsEqual<T extends Record<string, string>>(params: {
  keys: (keyof T)[];
  fields: Partial<T>;
  newFields: Partial<T>;
}) {
  const { keys, fields, newFields } = params;
  for (const key of keys) {
    if (fields[key] !== newFields[key]) {
      return false;
    }
  }

  return true;
}
function updateURLState<T extends Record<string, string>>(params: {
  state: T;
  query: URLSearchParams;
}) {
  const { state, query } = params;
  const queryParams = new URLSearchParams(query);

  for (const key of Object.keys(state)) {
    queryParams.set(key, state[key]);
  }

  return queryParams;
}

export type UseURLStateParams<T extends Record<string, string>> = {
  initialState: T;
  matchingPath?: string;
};
export type UseURLStateReturn<T extends Record<string, string>> = Readonly<
  [T, (partialState: SetStateAction<Partial<T>>, replace?: boolean) => void]
>;

export function useURLState<T extends Record<string, string>>(
  params: UseURLStateParams<T>
) {
  const { initialState } = params;
  const [search, setSearch] = useSearchParams(initialState);

  const urlFields = getQueryFields({
    fields: Object.keys(initialState),
    query: search,
  });
  const state: T = {
    ...initialState,
    ...urlFields,
  };

  function setState(partialState: SetStateAction<Partial<T>>, replace = false) {
    const newState =
      typeof partialState === 'function'
        ? { ...state, ...partialState(state) }
        : {
            ...state,
            ...partialState,
          };

    setSearch(
      (prev) =>
        updateURLState({
          state: newState,
          query: prev,
        }),
      {
        replace,
      }
    );
  }

  if (
    !areFieldsEqual({
      keys: Object.keys(initialState),
      newFields: state,
      fields: urlFields,
    })
  ) {
    setState(state);
  }

  return [state, setState] as const;
}
