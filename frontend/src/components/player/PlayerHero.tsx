import type { PlayerDetail, CareerArcData, CareerStatItem } from '@/types'

// ── Aggregation helpers ───────────────────────────────────────────────────────

function sumHitting(stats: CareerStatItem[], field: keyof NonNullable<CareerStatItem['hitting']>): number {
  return stats.reduce((acc, s) => acc + ((s.hitting?.[field] as number) ?? 0), 0)
}

function sumPitching(stats: CareerStatItem[], field: keyof NonNullable<CareerStatItem['pitching']>): number {
  return stats.reduce((acc, s) => acc + ((s.pitching?.[field] as number) ?? 0), 0)
}

/** Weighted average of a rate stat. Returns null if denominator is 0. */
function weightedRate(
  stats: CareerStatItem[],
  rateField: keyof NonNullable<CareerStatItem['hitting']> | keyof NonNullable<CareerStatItem['pitching']>,
  denomField: keyof NonNullable<CareerStatItem['hitting']> | keyof NonNullable<CareerStatItem['pitching']>,
  group: 'hitting' | 'pitching',
): number | null {
  let numerator = 0
  let denominator = 0
  for (const s of stats) {
    const grp = s[group]
    if (!grp) continue
    const rate = (grp as unknown as Record<string, number>)[rateField as string] ?? 0
    const denom = (grp as unknown as Record<string, number>)[denomField as string] ?? 0
    numerator += rate * denom
    denominator += denom
  }
  return denominator === 0 ? null : numerator / denominator
}

/** True career AVG: sum(hits) / sum(at_bats) */
function careerAvg(stats: CareerStatItem[]): number | null {
  const h = sumHitting(stats, 'hits')
  const ab = sumHitting(stats, 'at_bats')
  return ab === 0 ? null : h / ab
}

/** True career WHIP: (sum(hits_allowed) + sum(walks_allowed)) / sum(innings_pitched) */
function careerWhip(stats: CareerStatItem[]): number | null {
  const ha = sumPitching(stats, 'hits_allowed')
  const bb = sumPitching(stats, 'walks_allowed')
  const ip = sumPitching(stats, 'innings_pitched')
  return ip === 0 ? null : (ha + bb) / ip
}

/** Determine role from career_stats workload presence. */
function playerRole(stats: CareerStatItem[]): 'hitter' | 'pitcher' | 'two-way' {
  const hasHitting = stats.some(s => (s.hitting?.plate_appearances ?? 0) > 0)
  const hasPitching = stats.some(s => (s.pitching?.innings_pitched ?? 0) > 0)
  if (hasHitting && hasPitching) return 'two-way'
  if (hasPitching) return 'pitcher'
  return 'hitter'
}

function fmt(value: number | null, decimals: number): string {
  if (value === null || isNaN(value)) return '—'
  return value.toFixed(decimals)
}

function fmtAvg(value: number | null): string {
  if (value === null || isNaN(value)) return '—'
  // Format as .XXX (no leading zero)
  return value.toFixed(3).replace(/^0/, '')
}

function initials(first: string, last: string): string {
  return `${first[0] ?? ''}${last[0] ?? ''}`.toUpperCase()
}

// ── Stat card ────────────────────────────────────────────────────────────────

interface StatCardProps {
  label: string
  value: string
  accent?: boolean
}

function StatCard({ label, value, accent = false }: StatCardProps) {
  return (
    <div className="flex min-w-[80px] flex-col justify-center rounded-[6px] border border-border bg-panel px-3 py-3 sm:min-w-[96px]">
      <div className="text-[11px] uppercase tracking-[0.06em] text-text-subtle">{label}</div>
      <div className={`mt-1 font-mono text-[18px] font-medium tabular-nums ${accent ? 'text-accent' : 'text-text'}`}>
        {value}
      </div>
    </div>
  )
}

// ── Component ─────────────────────────────────────────────────────────────────

interface PlayerHeroProps {
  player: PlayerDetail
  arcData: CareerArcData | null
}

export default function PlayerHero({ player, arcData }: PlayerHeroProps) {
  const stats = player.career_stats ?? []
  const role = playerRole(stats)

  // Current arc score: latest historical timeline item with a non-zero value_score.
  // Sub-threshold seasons at the tail of a career (score = 0) are skipped — showing
  // 0.0 for a player whose last 3 seasons were injury/decline years would be misleading.
  const latestHistorical = arcData?.timeline
    .filter(t => !t.is_projection && t.value_score > 0)
    .slice(-1)[0]
  const currentArcScore = latestHistorical?.value_score ?? null

  // Peak arc score from arcMeta
  const peakArcScore = arcData?.arc?.peak_value_score ?? null

  // Career span
  const years = stats.map(s => s.year)
  const careerSpan =
    years.length > 0
      ? years[0] === years[years.length - 1]
        ? `${years[0]}`
        : `${years[0]}–${years[years.length - 1]}`
      : null

  // Position-appropriate cards (4 stats after the two arc cards)
  const positionCards: StatCardProps[] =
    role === 'pitcher'
      ? [
          { label: 'Career W', value: fmt(sumPitching(stats, 'wins'), 0) },
          { label: 'Career ERA', value: fmt(weightedRate(stats, 'era', 'innings_pitched', 'pitching'), 2) },
          { label: 'Career WHIP', value: fmt(careerWhip(stats), 2) },
          { label: 'Career K/9', value: fmt(weightedRate(stats, 'strikeouts_per_9', 'innings_pitched', 'pitching'), 1) },
        ]
      : role === 'two-way'
      ? [
          { label: 'Career HR', value: fmt(sumHitting(stats, 'home_runs'), 0) },
          { label: 'Career AVG', value: fmtAvg(careerAvg(stats)) },
          { label: 'Career ERA', value: fmt(weightedRate(stats, 'era', 'innings_pitched', 'pitching'), 2) },
          { label: 'Career WHIP', value: fmt(careerWhip(stats), 2) },
        ]
      : [
          { label: 'Career HR', value: fmt(sumHitting(stats, 'home_runs'), 0) },
          { label: 'Career AVG', value: fmtAvg(careerAvg(stats)) },
          { label: 'Career OBP', value: fmtAvg(weightedRate(stats, 'obp', 'plate_appearances', 'hitting')) },
          { label: 'Career RBI', value: fmt(sumHitting(stats, 'rbi'), 0) },
        ]

  const latestTeam = player.latest_season?.team_name ?? null

  return (
    <div className="rounded-[8px] border border-border bg-elevated px-5 py-5 sm:px-6 sm:py-6">
      <div className="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">

        {/* Avatar + identity */}
        <div className="flex items-center gap-4">
          {player.image_url ? (
            <img
              src={player.image_url}
              alt={`${player.full_name} headshot`}
              className="h-20 w-20 shrink-0 rounded-full border border-border bg-panel object-cover object-top"
            />
          ) : (
            <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-full border border-border bg-panel font-mono text-xl text-text-muted">
              {initials(player.first_name, player.last_name)}
            </div>
          )}

          <div>
            <h1 className="font-display text-4xl leading-none tracking-[0.04em] text-text sm:text-5xl">
              {player.full_name.toUpperCase()}
            </h1>
            <p className="mt-1.5 text-[13px] text-text-muted">
              {player.position}
              {latestTeam ? ` · ${latestTeam}` : ''}
            </p>
            <div className="mt-1.5 flex items-center gap-2">
              <span
                className={`rounded-full border px-2 py-0.5 text-[10px] font-medium ${
                  player.active
                    ? 'border-success text-success'
                    : 'border-border text-text-subtle'
                }`}
              >
                {player.active ? 'Active' : 'Retired'}
              </span>
              {careerSpan && (
                <span className="text-[11px] text-text-subtle">{careerSpan}</span>
              )}
            </div>
          </div>
        </div>

        {/* Stat cards — 2 arc + 4 position-specific */}
        <div className="flex flex-wrap gap-2 sm:flex-1 sm:flex-nowrap sm:justify-between">
          <StatCard
            label="Arc Score"
            value={currentArcScore !== null ? fmt(currentArcScore, 1) : '—'}
            accent
          />
          <StatCard
            label="Peak Arc"
            value={peakArcScore !== null ? fmt(peakArcScore, 1) : '—'}
            accent
          />
          {positionCards.map(card => (
            <StatCard key={card.label} label={card.label} value={card.value} />
          ))}
        </div>

      </div>
    </div>
  )
}
