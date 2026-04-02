import { Link } from 'react-router-dom'
import type { PlayerListItem } from '@/types'

function initials(first: string, last: string): string {
  return `${first[0] ?? ''}${last[0] ?? ''}`.toUpperCase()
}

interface PlayerCardProps {
  player: PlayerListItem
}

export default function PlayerCard({ player }: PlayerCardProps) {
  const season = player.latest_season

  return (
    <Link
      to={`/players/${player.id}`}
      className="group flex items-center gap-3 rounded-[8px] border border-border bg-elevated px-[14px] py-[10px] transition-colors hover:border-border-strong"
    >
      {/* Initials avatar */}
      <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-border bg-panel font-mono text-[11px] text-text-muted">
        {initials(player.first_name, player.last_name)}
      </div>

      {/* Name + position · team */}
      <div className="min-w-0 flex-1">
        <div className="truncate text-[13px] font-medium text-text transition-colors group-hover:text-accent">
          {player.full_name}
        </div>
        <div className="text-[10px] text-text-subtle">
          {player.position}
          {season ? ` · ${season.team_name}` : ''}
        </div>
      </div>

      {/* Active / Retired badge */}
      <span
        className={`shrink-0 rounded-full border px-2 py-0.5 text-[10px] font-medium ${
          player.active
            ? 'border-success text-success'
            : 'border-border text-text-subtle'
        }`}
      >
        {player.active ? 'Active' : 'Retired'}
      </span>

      {/* Value score + season year */}
      <div className="w-[52px] shrink-0 text-right">
        <div className="font-mono text-[14px] font-medium text-accent">
          {season ? season.value_score.toFixed(1) : '—'}
        </div>
        <div className="text-[10px] text-text-subtle">
          {season ? season.year : 'no data'}
        </div>
      </div>
    </Link>
  )
}
