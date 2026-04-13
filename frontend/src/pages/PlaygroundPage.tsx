import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getPlayer, getPlaygroundQuery } from '@/api'
import type { PlaygroundQueryParams } from '@/types'
import PlaygroundLayout from '@/components/playground/PlaygroundLayout'
import PlaygroundFilterPanel from '@/components/playground/PlaygroundFilterPanel'
import PlaygroundResultsHeader from '@/components/playground/PlaygroundResultsHeader'
import PlaygroundResultsTable from '@/components/playground/PlaygroundResultsTable'

// Fields excluded from the active-filter count — they're pagination/sort controls,
// not meaningful user-set filters.
const EXCLUDED_PARAM_KEYS: (keyof PlaygroundQueryParams)[] = ['sort', 'page', 'page_size']

function countActiveFilters(params: PlaygroundQueryParams | null): number {
  if (!params) return 0
  return (Object.keys(params) as (keyof PlaygroundQueryParams)[]).filter(
    (k) => !EXCLUDED_PARAM_KEYS.includes(k) && params[k] !== undefined,
  ).length
}

export default function PlaygroundPage() {
  const [searchParams] = useSearchParams()
  const playerId = Number(searchParams.get('player')) || null

  const [committedParams, setCommittedParams] = useState<PlaygroundQueryParams | null>(null)
  const [sort, setSort] = useState('-value_score')
  const [page, setPage] = useState(1)

  // Resolve player name when ?player=:id is in the URL so we can pre-fill the
  // search field and auto-run the query on load.
  const { data: preloadPlayer } = useQuery({
    queryKey: ['player', playerId],
    queryFn: () => getPlayer(playerId!),
    enabled: !!playerId,
  })

  // Auto-run once the player resolves. Dep on .id avoids re-firing if the
  // object reference changes without the player actually changing.
  useEffect(() => {
    if (preloadPlayer) {
      setCommittedParams({ q: preloadPlayer.full_name })
      setPage(1)
    }
  }, [preloadPlayer?.id]) // eslint-disable-line react-hooks/exhaustive-deps

  const initialParams: PlaygroundQueryParams | undefined = preloadPlayer
    ? { q: preloadPlayer.full_name }
    : undefined

  const { data, isLoading } = useQuery({
    queryKey: ['playground-query', committedParams, sort, page],
    queryFn: () => getPlaygroundQuery({ ...committedParams, sort, page, page_size: 25 }),
    enabled: committedParams !== null,
  })

  function handleSearch(params: PlaygroundQueryParams) {
    setCommittedParams(params)
    setPage(1)
  }

  function handleReset() {
    setCommittedParams(null)
    setPage(1)
  }

  function handleSortChange(s: string) {
    setSort(s)
    setPage(1)
  }

  const items = data?.data ?? []
  const total = data?.meta.total ?? 0
  const totalPages = data?.meta.total_pages ?? 1
  const group = committedParams?.group

  const header = (
    <div className="border-b border-border bg-gradient-to-b from-brand/5 to-transparent">
      <div className="shell-container py-8">
        <h1 className="font-display text-[42px] tracking-[0.04em]">
          Stat <span className="text-brand">Playground</span>
        </h1>
        <p className="mt-1 text-[14px] text-text-subtle">
          Build custom queries, explore stats, and compare players.
        </p>
      </div>
    </div>
  )

  const sidebar = (
    <PlaygroundFilterPanel
      initialParams={initialParams}
      onSearch={handleSearch}
      onReset={handleReset}
    />
  )

  const main = committedParams === null ? (
    <div className="flex h-full min-h-[200px] items-center justify-center text-[13px] text-text-subtle">
      Set your filters and click Run Query to explore the data.
    </div>
  ) : (
    <div>
      <PlaygroundResultsHeader
        total={total}
        activeFilterCount={countActiveFilters(committedParams)}
        sort={sort}
        onSortChange={handleSortChange}
      />
      <PlaygroundResultsTable
        items={items}
        isLoading={isLoading}
        page={page}
        totalPages={totalPages}
        onPageChange={setPage}
        group={group}
      />
    </div>
  )

  return <PlaygroundLayout header={header} sidebar={sidebar} main={main} />
}

