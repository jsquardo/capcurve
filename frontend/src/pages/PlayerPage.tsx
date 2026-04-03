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

  const { data: arcResponse } = useQuery({
    queryKey: ['career-arc', playerId],
    queryFn: () => getCareerArc(playerId),
    enabled: !!playerId,
  })

  if (isLoading) return <PlayerPageSkeleton />
  if (!player) return <div className="text-danger">Player not found</div>

  const arcData = arcResponse?.data ?? null

  return (
    <div className="space-y-6">
      <PlayerHero player={player} arcData={arcData} />
      {arcData && <CareerArcChart arcData={arcData} />}
      {arcData && <ProjectionPanel projection={arcData.projection} />}
      <SeasonStatsTable seasons={player.career_stats ?? []} arcMeta={arcData?.arc} />
      {arcData && arcData.projection.comparables.length > 0 && (
        <ComparablePlayersRow comparables={arcData.projection.comparables} />
      )}
    </div>
  )
}
