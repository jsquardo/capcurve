import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getPlayer, getCareerArc } from '@/api'
import PlayerPageSkeleton from '@/components/players/PlayerPageSkeleton'
import PlayerHero from '@/components/player/PlayerHero'
import CareerArcChart from '@/components/player/CareerArcChart'
import ProjectionPanel from '@/components/player/ProjectionPanel'
import SeasonStatsTable from '@/components/player/SeasonStatsTable'
import ComparablePlayersRow from '@/components/player/ComparablePlayersRow'

export default function PlayerPage() {
  const { id } = useParams<{ id: string }>()
  const playerId = Number(id)

  const { data: player, isLoading } = useQuery({
    queryKey: ['player', playerId],
    queryFn: () => getPlayer(playerId),
  })

  const {
    data: arcResponse,
    isLoading: isArcLoading,
    isError: isArcError,
  } = useQuery({
    queryKey: ['career-arc', playerId],
    queryFn: () => getCareerArc(playerId),
    enabled: !!playerId,
  })

  if (isLoading) return <PlayerPageSkeleton />
  if (!player) return <div className="text-danger">Player not found</div>

  const arcData = arcResponse?.data ?? null

  // Placeholder shown in the chart region while the arc request is in-flight or
  // if it fails — keeps the hero and stats table visible rather than hiding the page.
  const arcChartContent = (() => {
    if (isArcLoading) {
      return (
        <div className="h-[320px] animate-pulse rounded-[8px] bg-panel" />
      )
    }
    if (isArcError) {
      return (
        <div className="flex h-[160px] items-center justify-center rounded-[8px] border border-border bg-panel text-[13px] text-text-subtle">
          Career arc data could not be loaded.
        </div>
      )
    }
    if (arcData) {
      return <CareerArcChart arcData={arcData} />
    }
    return null
  })()

  return (
    <div className="space-y-6">
      <PlayerHero player={player} arcData={arcData} />
      {arcChartContent}
      {arcData && <ProjectionPanel projection={arcData.projection} />}
      <SeasonStatsTable seasons={player.career_stats ?? []} arcMeta={arcData?.arc} />
      {arcData && arcData.projection.comparables.length > 0 && (
        <ComparablePlayersRow comparables={arcData.projection.comparables} />
      )}
      <div className="rounded-[8px] border border-border bg-elevated px-5 py-5 sm:px-6 sm:py-6">
        <div className="flex flex-col items-start gap-3">
          <p className="text-[13px] font-medium uppercase tracking-widest text-text-subtle">
            Stats Playground
          </p>
          {/* TODO: swap for <Link to={`/playground?player=${playerId}`}> and remove aria-disabled/describedby when PlaygroundPage is built (Phase 4) */}
          <button
            type="button"
            aria-disabled="true"
            aria-describedby="playground-cta-hint"
            onClick={(e) => e.preventDefault()}
            className="cursor-not-allowed rounded-[8px] border border-border bg-panel px-5 py-2.5 text-[13px] font-medium text-text-subtle opacity-40"
          >
            Explore in Playground →
          </button>
          <p id="playground-cta-hint" className="text-[12px] text-text-subtle opacity-60">
            Coming soon — available in Phase 4
          </p>
        </div>
      </div>
    </div>
  )
}
