import { Link } from 'react-router-dom'
import type { ComparablePlayer } from '@/types'

function initials(fullName: string): string {
  const parts = fullName.trim().split(/\s+/)
  const first = parts[0]?.[0] ?? ''
  const last = parts[parts.length - 1]?.[0] ?? ''
  return (first + (parts.length > 1 ? last : '')).toUpperCase()
}

interface ComparablePlayersRowProps {
  comparables: ComparablePlayer[]
}

export default function ComparablePlayersRow({ comparables }: ComparablePlayersRowProps) {
  if (comparables.length === 0) return null

  return (
    <div className="rounded-[8px] border border-border bg-elevated px-5 py-5 sm:px-6 sm:py-6">
      <h2 className="mb-4 font-display text-2xl tracking-[0.04em] text-text">
        Similar <span className="text-accent">Arcs</span>
      </h2>

      <div className="flex gap-3 overflow-x-auto pb-2 [scrollbar-width:thin] [&::-webkit-scrollbar]:h-[3px] [&::-webkit-scrollbar-thumb]:rounded [&::-webkit-scrollbar-thumb]:bg-border">
        {comparables.map(c => (
          <Link
            key={c.player_id}
            to={`/players/${c.player_id}`}
            className="group flex min-w-[150px] max-w-[180px] shrink-0 items-center gap-3 rounded-[8px] border border-border bg-panel px-3 py-3 transition-colors hover:border-accent"
          >
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-border bg-elevated font-mono text-[12px] text-text-muted transition-colors group-hover:border-accent/50">
              {initials(c.full_name)}
            </div>
            <div className="min-w-0">
              <div className="truncate text-[13px] font-medium text-text transition-colors group-hover:text-accent">
                {c.full_name}
              </div>
              <div className="text-[11px] text-text-subtle">{c.position}</div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
