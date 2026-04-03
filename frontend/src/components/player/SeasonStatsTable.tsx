import { useState } from 'react'
import type { CareerStatItem, CareerArcMeta } from '@/types'

// ── Helpers ───────────────────────────────────────────────────────────────────

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

type ColKey =
  | 'year' | 'team_name' | 'age' | 'value_score'
  | 'home_runs' | 'batting_avg' | 'obp' | 'slg' | 'ops' | 'plate_appearances'
  | 'wins' | 'era' | 'whip' | 'innings_pitched' | 'strikeouts_per_9'

interface ColDef {
  key: ColKey
  label: string
  render: (s: CareerStatItem) => string
  numeric: boolean
}

const BASE_COLS: ColDef[] = [
  { key: 'year',        label: 'Year',  render: s => s.year.toString(),           numeric: false },
  { key: 'team_name',   label: 'Team',  render: s => s.team_name || '—',          numeric: false },
  { key: 'age',         label: 'Age',   render: s => s.age.toString(),            numeric: true  },
  { key: 'value_score', label: 'Arc',   render: s => fmtDec(s.value_score, 1),    numeric: true  },
]

const HITTING_COLS: ColDef[] = [
  { key: 'plate_appearances', label: 'PA',  render: s => fmtInt(s.hitting?.plate_appearances),  numeric: true },
  { key: 'home_runs',         label: 'HR',  render: s => fmtInt(s.hitting?.home_runs),          numeric: true },
  { key: 'batting_avg',       label: 'AVG', render: s => fmtRate(s.hitting?.batting_avg),       numeric: true },
  { key: 'obp',               label: 'OBP', render: s => fmtRate(s.hitting?.obp),               numeric: true },
  { key: 'slg',               label: 'SLG', render: s => fmtRate(s.hitting?.slg),               numeric: true },
  { key: 'ops',               label: 'OPS', render: s => fmtRate(s.hitting?.ops),               numeric: true },
]

const PITCHING_COLS: ColDef[] = [
  { key: 'wins',             label: 'W',    render: s => fmtInt(s.pitching?.wins),              numeric: true },
  { key: 'innings_pitched',  label: 'IP',   render: s => fmtDec(s.pitching?.innings_pitched, 1), numeric: true },
  { key: 'era',              label: 'ERA',  render: s => fmtDec(s.pitching?.era, 2),            numeric: true },
  { key: 'whip',             label: 'WHIP', render: s => fmtDec(s.pitching?.whip, 2),           numeric: true },
  { key: 'strikeouts_per_9', label: 'K/9',  render: s => fmtDec(s.pitching?.strikeouts_per_9, 1), numeric: true },
]

// ── Sort helpers ──────────────────────────────────────────────────────────────

type SortDir = 'asc' | 'desc'

function sortValue(s: CareerStatItem, key: ColKey): number | string {
  switch (key) {
    case 'year':               return s.year
    case 'age':                return s.age
    case 'value_score':        return s.value_score
    case 'plate_appearances':  return s.hitting?.plate_appearances ?? -1
    case 'home_runs':          return s.hitting?.home_runs ?? -1
    case 'batting_avg':        return s.hitting?.batting_avg ?? -1
    case 'obp':                return s.hitting?.obp ?? -1
    case 'slg':                return s.hitting?.slg ?? -1
    case 'ops':                return s.hitting?.ops ?? -1
    case 'wins':               return s.pitching?.wins ?? -1
    case 'innings_pitched':    return s.pitching?.innings_pitched ?? -1
    case 'era':                return s.pitching?.era ?? -1
    case 'whip':               return s.pitching?.whip ?? -1
    case 'strikeouts_per_9':   return s.pitching?.strikeouts_per_9 ?? -1
    case 'team_name':          return s.team_name
    default:                   return 0
  }
}

function sortedSeasons(seasons: CareerStatItem[], col: ColKey, dir: SortDir): CareerStatItem[] {
  return [...seasons].sort((a, b) => {
    const av = sortValue(a, col)
    const bv = sortValue(b, col)
    if (av < bv) return dir === 'asc' ? -1 : 1
    if (av > bv) return dir === 'asc' ? 1 : -1
    return 0
  })
}

// ── Component ─────────────────────────────────────────────────────────────────

interface SeasonStatsTableProps {
  seasons: CareerStatItem[]
  arcMeta?: CareerArcMeta | null
}

export default function SeasonStatsTable({ seasons, arcMeta }: SeasonStatsTableProps) {
  const [expanded, setExpanded] = useState(true)
  const [sortCol, setSortCol] = useState<ColKey>('year')
  const [sortDir, setSortDir] = useState<SortDir>('desc')

  // Scan ALL seasons to determine which column sets to show
  const hasHitting  = seasons.some(s => (s.hitting?.plate_appearances ?? 0) > 0)
  const hasPitching = seasons.some(s => (s.pitching?.innings_pitched ?? 0) > 0)

  const columns: ColDef[] = [
    ...BASE_COLS,
    ...(hasHitting  ? HITTING_COLS  : []),
    ...(hasPitching ? PITCHING_COLS : []),
  ]

  const rows = sortedSeasons(seasons, sortCol, sortDir)

  function handleSort(key: ColKey) {
    if (key === sortCol) {
      setSortDir(d => d === 'asc' ? 'desc' : 'asc')
    } else {
      setSortCol(key)
      // Default new sort direction: year → desc, everything else → desc
      setSortDir('desc')
    }
  }

  function isPeakYear(year: number): boolean {
    if (!arcMeta) return false
    return year >= arcMeta.peak_year_start && year <= arcMeta.peak_year_end
  }

  return (
    <div className="rounded-[8px] border border-border bg-elevated px-5 py-5 sm:px-6 sm:py-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-display text-2xl tracking-[0.04em] text-text">
          Season <span className="text-accent">Stats</span>
        </h2>
        <button
          onClick={() => setExpanded(e => !e)}
          className="text-[12px] text-link transition-colors hover:text-text"
        >
          {expanded ? '▾ Collapse' : '► Expand'}
        </button>
      </div>

      {/* Table */}
      {expanded && (
        <div className="overflow-x-auto [-webkit-overflow-scrolling:touch] [scrollbar-width:thin]">
          <table className="w-full border-collapse text-[13px]" style={{ minWidth: 680 }}>
            <thead>
              <tr>
                {columns.map(col => (
                  <th
                    key={col.key}
                    onClick={() => handleSort(col.key)}
                    className={`select-none whitespace-nowrap border-b border-border pb-2 pt-1 font-mono text-[11px] font-medium uppercase tracking-[0.04em] text-text-subtle transition-colors hover:text-text ${col.numeric ? 'text-right' : 'text-left'} cursor-pointer px-3`}
                  >
                    {col.label}
                    {sortCol === col.key && (
                      <span className="ml-1 text-[9px] text-accent">
                        {sortDir === 'desc' ? '▼' : '▲'}
                      </span>
                    )}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map(s => {
                const isPeak = isPeakYear(s.year)
                return (
                  <tr
                    key={`${s.year}-${s.team_name}`}
                    className="group transition-colors hover:bg-panel"
                    style={isPeak ? { background: 'rgba(240,192,64,0.04)' } : undefined}
                  >
                    {columns.map(col => {
                      const isArcCol = col.key === 'value_score'
                      const isYearCol = col.key === 'year'
                      return (
                        <td
                          key={col.key}
                          className={`whitespace-nowrap border-b border-border/40 px-3 py-2.5 last-of-type:border-0 ${col.numeric ? 'text-right' : 'text-left'} ${isArcCol ? 'font-display text-[17px]' : 'font-mono'} ${isArcCol && isPeak ? 'text-accent' : isArcCol ? 'text-text-muted' : 'text-text'}`}
                        >
                          {col.render(s)}
                          {isYearCol && isPeak && (
                            <span
                              className="ml-1.5 rounded px-1 py-0.5 font-mono text-[9px] font-medium"
                              style={{ background: 'rgba(240,192,64,0.12)', color: '#f0c040', border: '1px solid rgba(240,192,64,0.25)' }}
                            >
                              PEAK
                            </span>
                          )}
                        </td>
                      )
                    })}
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
