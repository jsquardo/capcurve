import { Link } from 'react-router-dom'
import type { PlaygroundQueryItem } from '@/types'

// ── Formatters (same helpers as SeasonStatsTable) ─────────────────────────────

function fmtRate(v: number | null | undefined): string {
  if (v == null || isNaN(v)) return '—'
  return v.toFixed(3).replace(/^0/, '')
}

function fmtDec(v: number | null | undefined, dp: number): string {
  if (v == null || isNaN(v)) return '—'
  return v.toFixed(dp)
}

function fmtInt(v: number | null | undefined): string {
  if (v == null || isNaN(v)) return '—'
  return Math.round(v).toString()
}

// ── Column definitions ────────────────────────────────────────────────────────

interface ColDef {
  key: string
  label: string
  render: (item: PlaygroundQueryItem) => string
  numeric: boolean
}

const HITTING_COLS: ColDef[] = [
  { key: 'pa',  label: 'PA',  render: (i) => fmtInt(i.hitting?.plate_appearances), numeric: true },
  { key: 'hr',  label: 'HR',  render: (i) => fmtInt(i.hitting?.home_runs),         numeric: true },
  { key: 'avg', label: 'AVG', render: (i) => fmtRate(i.hitting?.batting_avg),      numeric: true },
  { key: 'obp', label: 'OBP', render: (i) => fmtRate(i.hitting?.obp),              numeric: true },
  { key: 'slg', label: 'SLG', render: (i) => fmtRate(i.hitting?.slg),              numeric: true },
  { key: 'ops', label: 'OPS', render: (i) => fmtRate(i.hitting?.ops),              numeric: true },
]

const PITCHING_COLS: ColDef[] = [
  { key: 'w',    label: 'W',    render: (i) => fmtInt(i.pitching?.wins),                numeric: true },
  { key: 'ip',   label: 'IP',   render: (i) => fmtDec(i.pitching?.innings_pitched, 1),  numeric: true },
  { key: 'era',  label: 'ERA',  render: (i) => fmtDec(i.pitching?.era, 2),              numeric: true },
  { key: 'whip', label: 'WHIP', render: (i) => fmtDec(i.pitching?.whip, 2),             numeric: true },
  { key: 'k9',   label: 'K/9',  render: (i) => fmtDec(i.pitching?.strikeouts_per_9, 1), numeric: true },
]

// Determine which stat columns to show based on group and actual data.
// When group is "all", scan items to pick the dominant type; add a Pos column
// so mixed rows are identifiable.
function resolveStatCols(
  items: PlaygroundQueryItem[],
  group: Props['group'],
): { cols: ColDef[]; showPosCol: boolean } {
  if (group === 'hitting') return { cols: HITTING_COLS, showPosCol: false }
  if (group === 'pitching') return { cols: PITCHING_COLS, showPosCol: false }

  if (items.length === 0) return { cols: HITTING_COLS, showPosCol: false }

  const pitcherCount = items.filter((i) => i.pitching !== null).length
  const isPitcherMajority = pitcherCount > items.length / 2
  const isMixed = pitcherCount > 0 && pitcherCount < items.length

  return {
    cols: isPitcherMajority ? PITCHING_COLS : HITTING_COLS,
    showPosCol: isMixed,
  }
}

// ── Arc score bar ─────────────────────────────────────────────────────────────

// Maps a value_score (roughly 0–100) to a bar width capped at 80px.
function ArcScoreCell({ score }: { score: number }) {
  const barWidth = Math.min(Math.max(score, 0), 100) * 0.6 // 0–60px
  return (
    <div className="flex items-center gap-2">
      <div
        className="h-[5px] rounded-full bg-accent"
        style={{ width: `${barWidth}px`, minWidth: 3 }}
      />
      <span className="font-mono tabular-nums text-[12px] text-accent">
        {fmtDec(score, 1)}
      </span>
    </div>
  )
}

// ── Skeleton rows ─────────────────────────────────────────────────────────────

function SkeletonRows({ cols }: { cols: number }) {
  return (
    <>
      {Array.from({ length: 8 }).map((_, i) => (
        <tr key={i}>
          {Array.from({ length: cols }).map((_, j) => (
            <td key={j} className="px-3 py-3">
              <div className="h-3 animate-pulse rounded bg-elevated" />
            </td>
          ))}
        </tr>
      ))}
    </>
  )
}

// ── Pagination ────────────────────────────────────────────────────────────────

function Pagination({
  page,
  totalPages,
  onPageChange,
}: {
  page: number
  totalPages: number
  onPageChange: (p: number) => void
}) {
  if (totalPages <= 1) return null
  return (
    <div className="mt-4 flex items-center justify-center gap-4 text-[13px]">
      <button
        type="button"
        disabled={page <= 1}
        onClick={() => onPageChange(page - 1)}
        className="text-text-muted disabled:opacity-30 hover:text-text transition-colors"
      >
        ← Prev
      </button>
      <span className="text-text-subtle">
        Page {page} of {totalPages}
      </span>
      <button
        type="button"
        disabled={page >= totalPages}
        onClick={() => onPageChange(page + 1)}
        className="text-text-muted disabled:opacity-30 hover:text-text transition-colors"
      >
        Next →
      </button>
    </div>
  )
}

// ── Component ─────────────────────────────────────────────────────────────────

interface Props {
  items: PlaygroundQueryItem[]
  isLoading: boolean
  page: number
  totalPages: number
  onPageChange: (page: number) => void
  group?: 'all' | 'hitting' | 'pitching'
}

export default function PlaygroundResultsTable({
  items,
  isLoading,
  page,
  totalPages,
  onPageChange,
  group,
}: Props) {
  const { cols: statCols, showPosCol } = resolveStatCols(items, group)

  // Total column count: player name + year + team + age + arc + optional pos + stat cols
  const totalCols = 5 + (showPosCol ? 1 : 0) + statCols.length

  const thClass = 'px-3 py-2.5 text-left text-[10px] font-semibold uppercase tracking-[0.5px] text-text-subtle whitespace-nowrap font-mono border-b border-border'
  const tdClass = 'px-3 py-2.5 border-b border-border/40 text-[13px]'

  return (
    <div>
      <div className="overflow-x-auto rounded-[8px] border border-border">
        <table className="w-full border-collapse text-[13px]">
          <thead className="bg-panel">
            <tr>
              <th className={thClass}>Player</th>
              <th className={`${thClass} text-right`}>Year</th>
              <th className={thClass}>Team</th>
              <th className={`${thClass} text-right`}>Age</th>
              <th className={thClass}>Arc</th>
              {showPosCol && <th className={thClass}>Pos</th>}
              {statCols.map((col) => (
                <th key={col.key} className={`${thClass} text-right`}>
                  {col.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <SkeletonRows cols={totalCols} />
            ) : items.length === 0 ? (
              <tr>
                <td colSpan={totalCols} className="px-3 py-12 text-center text-[13px] text-text-subtle">
                  No results — try adjusting your filters
                </td>
              </tr>
            ) : (
              items.map((item, idx) => (
                <tr key={`${item.player.id}-${item.season.year}-${idx}`} className="hover:bg-panel/60 transition-colors">
                  <td className={tdClass}>
                    <Link
                      to={`/players/${item.player.id}`}
                      className="font-medium text-text hover:text-accent transition-colors"
                    >
                      {item.player.full_name}
                    </Link>
                    <span className="ml-2 text-[11px] text-text-subtle">
                      {item.player.position}
                    </span>
                  </td>
                  <td className={`${tdClass} text-right font-mono tabular-nums text-text-muted`}>
                    {item.season.year}
                  </td>
                  <td className={`${tdClass} text-text-muted`}>
                    {item.season.team_name}
                  </td>
                  <td className={`${tdClass} text-right font-mono tabular-nums text-text-muted`}>
                    {item.season.age}
                  </td>
                  <td className={tdClass}>
                    <ArcScoreCell score={item.season.value_score} />
                  </td>
                  {showPosCol && (
                    <td className={`${tdClass} text-[11px] text-text-subtle`}>
                      {item.pitching !== null ? 'P' : 'H'}
                    </td>
                  )}
                  {statCols.map((col) => (
                    <td key={col.key} className={`${tdClass} text-right font-mono tabular-nums text-text-muted`}>
                      {col.render(item)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Pagination page={page} totalPages={totalPages} onPageChange={onPageChange} />
    </div>
  )
}
