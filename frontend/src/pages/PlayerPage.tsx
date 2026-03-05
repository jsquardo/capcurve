import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getPlayer, getCareerArc } from '@/api'

export default function PlayerPage() {
  const { id } = useParams<{ id: string }>()
  const playerId = Number(id)

  const { data: player, isLoading } = useQuery({
    queryKey: ['player', playerId],
    queryFn: () => getPlayer(playerId),
  })

  const { data: _arcData } = useQuery({
    queryKey: ['career-arc', playerId],
    queryFn: () => getCareerArc(playerId),
    enabled: !!playerId,
  })

  if (isLoading) return <div className="text-neutral">Loading...</div>
  if (!player) return <div className="text-overpaid">Player not found</div>

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-5xl text-white">
          {player.first_name.toUpperCase()} {player.last_name.toUpperCase()}
        </h1>
        <p className="text-neutral mt-1">{player.position} · {player.active ? 'Active' : 'Retired'}</p>
      </div>
      <p className="text-neutral">Charts coming soon...</p>
    </div>
  )
}
