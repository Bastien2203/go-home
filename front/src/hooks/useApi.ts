import { useState, useEffect, useCallback } from "react";

export function useApi<T>(apiFunction: () => Promise<T>) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiFunction();
      setData(result);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [apiFunction]);

  useEffect(() => {
    execute();
  }, [execute]);

  return {
    data,
    loading,
    error,
    refresh: execute,
  };
}