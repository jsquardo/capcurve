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
  const trimmed = query.trim();
  const [debouncedQuery, setDebouncedQuery] = useState("");

  useEffect(() => {
    // Clear immediately when input drops below threshold — no stale results.
    // Trimming before the length check prevents whitespace-only input from
    // passing the threshold and firing a request.
    if (trimmed.length < MIN_QUERY_LENGTH) {
      setDebouncedQuery("");
      return;
    }
    const timer = setTimeout(() => setDebouncedQuery(trimmed), DEBOUNCE_MS);
    return () => clearTimeout(timer);
  }, [trimmed]);

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
