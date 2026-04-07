import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { searchPlayers } from "@/api";
import type { PlayerListItem } from "@/types";

const MIN_QUERY_LENGTH = 2;
const DEBOUNCE_MS = 300;
const RESULT_LIMIT = 8;

export function usePlayerSearch(query: string): {
  results: PlayerListItem[];
  isLoading: boolean;
} {
  const [debouncedQuery, setDebouncedQuery] = useState("");

  useEffect(() => {
    // Clear immediately when input drops below threshold — no stale results
    if (query.length < MIN_QUERY_LENGTH) {
      setDebouncedQuery("");
      return;
    }
    const timer = setTimeout(() => setDebouncedQuery(query), DEBOUNCE_MS);
    return () => clearTimeout(timer);
  }, [query]);

  const { data, isLoading } = useQuery({
    queryKey: ["player-search", debouncedQuery],
    queryFn: () => searchPlayers(debouncedQuery, RESULT_LIMIT),
    enabled: debouncedQuery.length >= MIN_QUERY_LENGTH,
    staleTime: 60_000,
  });

  return {
    results: data ?? [],
    isLoading,
  };
}
