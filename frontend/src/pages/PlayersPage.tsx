import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getPlayers } from '@/api'
import PlayerListHero from '@/components/players/PlayerListHero'
import PlayerFilters from '@/components/players/PlayerFilters'
import PlayerCard from '@/components/players/PlayerCard'
import PlayerListSkeleton from '@/components/players/PlayerListSkeleton'
import LeaderboardPagination from '@/components/leaderboards/LeaderboardPagination'

const PAGE_SIZE = 25

export default function PlayersPage() {
  const [q, setQ] = useState('')
  const [debouncedQ, setDebouncedQ] = useState('')
  const [active, setActive] = useState<boolean | undefined>(undefined)
  const [position, setPosition] = useState('')
  const [sort, setSort] = useState('name')
  const [page, setPage] = useState(1)

  // Debounce the text search — avoid firing a query on every keystroke.
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQ(q), 300)
    return () => clearTimeout(timer)
  }, [q])

  // Reset to page 1 whenever the debounced query changes.
  useEffect(() => {
    setPage(1)
  }, [debouncedQ])

  function handleActiveChange(val: boolean | undefined) {
    setActive(val)
    setPage(1)
  }

  function handlePositionChange(val: string) {
    setPosition(val)
    setPage(1)
  }

  function handleSortChange(val: string) {
    setSort(val)
    setPage(1)
  }

  const { data, isLoading, isError } = useQuery({
    queryKey: ['players', debouncedQ, active, position, sort, page],
    queryFn: () =>
      getPlayers({
        q: debouncedQ || undefined,
        active,
        position: position || undefined,
        sort,
        page,
        page_size: PAGE_SIZE,
      }),
  })

  const players = data?.data ?? []
  const totalPages = data?.meta.total_pages ?? 1
  // Show null while loading so the hero badge displays a pulse rather than "0 PLAYERS".
  const total = isLoading ? null : (data?.meta.total ?? null)

  return (
    <div className="shell-container space-y-6 py-8">
      <PlayerListHero total={total} />
      <PlayerFilters
        q={q}
        active={active}
        position={position}
        sort={sort}
        onQChange={setQ}
        onActiveChange={handleActiveChange}
        onPositionChange={handlePositionChange}
        onSortChange={handleSortChange}
      />
      {isLoading ? (
        <PlayerListSkeleton rows={PAGE_SIZE} />
      ) : isError ? (
        <div className="rounded-[8px] border border-border bg-elevated px-6 py-10 text-center text-[13px] text-text-subtle">
          Could not load players. Please try again later.
        </div>
      ) : players.length === 0 ? (
        <div className="rounded-[8px] border border-border bg-elevated px-6 py-10 text-center text-[13px] text-text-subtle">
          No players match your filters.
        </div>
      ) : (
        <>
          <div className="space-y-1.5">
            {players.map(player => (
              <PlayerCard key={player.id} player={player} />
            ))}
          </div>
          <LeaderboardPagination
            page={page}
            totalPages={totalPages}
            onPrev={() => setPage(p => p - 1)}
            onNext={() => setPage(p => p + 1)}
          />
        </>
      )}
    </div>
  )
}
