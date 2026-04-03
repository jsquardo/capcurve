import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import type { LeaderboardCategory } from '@/types'
import { getLeaderboards } from '@/api'
import LeaderboardHero from '@/components/leaderboards/LeaderboardHero'
import LeaderboardCategoryTabs from '@/components/leaderboards/LeaderboardCategoryTabs'
import LeaderboardTable from '@/components/leaderboards/LeaderboardTable'
import LeaderboardPagination from '@/components/leaderboards/LeaderboardPagination'
import LeaderboardSkeleton from '@/components/leaderboards/LeaderboardSkeleton'

const PAGE_SIZE = 25

// A season is considered complete once October 1 of that year has passed
// (mirrors the backend's completed-season cutoff rule).
const _now = new Date()
const CURRENT_SEASON = _now.getMonth() >= 9 ? _now.getFullYear() : _now.getFullYear() - 1

export default function LeaderboardsPage() {
  const [activeCategory, setActiveCategory] = useState<LeaderboardCategory>('peak_arc')
  const [page, setPage] = useState(1)

  const { data, isLoading, isError } = useQuery({
    queryKey: ['leaderboards', activeCategory, page],
    queryFn: () => getLeaderboards({ category: activeCategory, page, page_size: PAGE_SIZE }),
  })

  // Reset to page 1 whenever the category changes.
  function handleSelectCategory(category: LeaderboardCategory) {
    setActiveCategory(category)
    setPage(1)
  }

  const leaders = data?.data.leaders ?? []
  const totalPages = data?.data.meta.total_pages ?? 1
  // Use a stable constant for the season badge — never read from leaders[0] so the
  // badge doesn't flip to ALL-TIME during loading/error states for seasonal categories.
  const heroSeason = activeCategory === 'peak_arc' ? null : CURRENT_SEASON

  return (
    <div className="shell-container space-y-6 py-8">
      <LeaderboardHero season={heroSeason} />
      <LeaderboardCategoryTabs
        activeCategory={activeCategory}
        onSelect={handleSelectCategory}
      />
      {isLoading ? (
        <LeaderboardSkeleton rows={PAGE_SIZE} />
      ) : isError ? (
        <div className="rounded-[8px] border border-border bg-elevated px-6 py-10 text-center text-[13px] text-text-subtle">
          Could not load leaderboard data. Please try again later.
        </div>
      ) : (
        <>
          <LeaderboardTable leaders={leaders} category={activeCategory} />
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
