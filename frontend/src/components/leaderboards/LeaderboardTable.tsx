import type { LeaderboardCategory, LeaderboardEntry } from '@/types'
import LeaderRow from '@/components/home/LeaderRow'

// Column header label for the value column varies by category.
const VALUE_LABELS: Record<LeaderboardCategory, string> = {
  peak_arc: 'VALUE SCORE',
  hr:       'HR',
  avg:      'AVG',
  era:      'ERA',
  k9:       'K/9',
}

// Format a numeric value for display in the value column.
function formatValue(value: number, category: LeaderboardCategory): string {
  switch (category) {
    case 'avg': return value.toFixed(3)
    case 'era': return value.toFixed(2)
    case 'k9':  return value.toFixed(1)
    case 'hr':  return String(Math.round(value))
    case 'peak_arc': return value.toFixed(1)
  }
}

// Compute bar fill percentage. ERA sorts ascending (lower = better),
// so we invert the scale for that category.
function barPct(entry: LeaderboardEntry, leaders: LeaderboardEntry[], category: LeaderboardCategory): number {
  if (leaders.length === 0) return 0
  const values = leaders.map(l => l.value)
  const max = Math.max(...values)
  const min = Math.min(...values)
  if (max === min) return 100
  return category === 'era'
    ? ((max - entry.value) / (max - min)) * 100
    : (entry.value / max) * 100
}

interface LeaderboardTableProps {
  leaders: LeaderboardEntry[]
  category: LeaderboardCategory
}

export default function LeaderboardTable({ leaders, category }: LeaderboardTableProps) {
  const valueLabel = VALUE_LABELS[category]

  if (leaders.length === 0) {
    return (
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-10 text-center text-[13px] text-text-subtle">
        No data available for this category.
      </div>
    )
  }

  return (
    <div>
      {/* Column headers */}
      <div className="mb-2 flex items-center gap-3 px-[14px]">
        <span className="w-[18px] shrink-0" />
        <span className="flex-1 text-[10px] uppercase tracking-wider text-text-subtle">Player</span>
        {/* spacer to align with bar */}
        <span className="w-24 shrink-0 text-right text-[10px] uppercase tracking-wider text-text-subtle" />
        <span className="w-[52px] shrink-0 text-right text-[10px] uppercase tracking-wider text-text-subtle">
          {valueLabel}
        </span>
      </div>

      {/* Rows */}
      <div className="space-y-1.5">
        {leaders.map(entry => (
          <LeaderRow
            key={entry.player_id}
            rank={entry.rank}
            name={entry.player_name}
            team={entry.team}
            position={entry.position}
            value={formatValue(entry.value, category)}
            barPct={barPct(entry, leaders, category)}
            // TODO: pass playerId={entry.player_id} once the page is wired to the live API;
            // mock player_ids are not real DB IDs and would produce broken /players/:id links.
          />
        ))}
      </div>
    </div>
  )
}
