import { Link } from 'react-router-dom'
import type { CareerArcProjection } from '@/types'

const COLOR_PROJECTION = '#9b72f8'

interface ProjectionPanelProps {
  projection: CareerArcProjection
}

export default function ProjectionPanel({ projection }: ProjectionPanelProps) {
  // Do not render for ineligible players (retired) — they have complete arcs.
  if (!projection.eligible) return null

  // Treat ready-but-empty the same as insufficient_data — there's nothing to show.
  const hasData =
    projection.status === 'ready' &&
    projection.points.length > 0 &&
    projection.points.some(p => p.value_score > 0)

  const insufficientReason =
    projection.reason.trim() ||
    (projection.status === 'insufficient_data'
      ? 'Not enough seasons to generate a projection.'
      : 'Projection data is not yet available.')

  // Confidence band is only meaningful when projected scores are non-zero.
  const bandMap = new Map(projection.confidence_band.map(b => [b.year, b]))
  const showBand = hasData

  const seasonCount = projection.points.length
  const compCount = projection.comparables.length

  return (
    <div
      className="rounded-[8px] border px-5 py-4 sm:px-6 sm:py-5"
      style={{
        background: 'rgba(155,114,248,0.04)',
        borderColor: 'rgba(155,114,248,0.18)',
      }}
    >
      {/* Header */}
      <div className="mb-3 flex items-center gap-2">
        <div
          className="h-2 w-2 shrink-0 rounded-full"
          style={{ background: COLOR_PROJECTION }}
        />
        <span className="text-[13px] font-medium" style={{ color: COLOR_PROJECTION }}>
          {hasData
            ? `${seasonCount}-Season Projection${compCount > 0 ? ` · ${compCount} comparable career${compCount === 1 ? '' : 's'}` : ''}`
            : 'Projection'}
        </span>
      </div>

      {/* Insufficient / empty state */}
      {!hasData && (
        <p className="text-[13px] text-text-subtle">{insufficientReason}</p>
      )}

      {/* Season grid */}
      {hasData && (
        <>
          <div className="flex gap-2 overflow-x-auto pb-1 [-webkit-overflow-scrolling:touch] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
            {projection.points.map(p => {
              const band = showBand ? bandMap.get(p.year) : undefined
              return (
                <div
                  key={p.year}
                  className="flex shrink-0 flex-col items-center rounded-[6px] border border-border bg-panel px-3 py-2.5 text-center"
                >
                  <div className="font-mono text-[10px] text-text-subtle">{p.year}</div>
                  <div
                    className="mt-0.5 font-display text-[22px] leading-none"
                    style={{ color: COLOR_PROJECTION }}
                  >
                    {p.value_score.toFixed(0)}
                  </div>
                  {band && (
                    <div className="mt-0.5 font-mono text-[9px] text-text-subtle">
                      {band.lower.toFixed(0)}–{band.upper.toFixed(0)}
                    </div>
                  )}
                </div>
              )
            })}
          </div>

          {/* Comparable chips */}
          {compCount > 0 && (
            <div className="mt-3 flex flex-wrap items-center gap-1.5 border-t pt-3" style={{ borderColor: 'rgba(37,45,66,0.6)' }}>
              <span className="text-[11px] text-text-subtle">Comps:</span>
              {projection.comparables.map(c => (
                <Link
                  key={c.player_id}
                  to={`/players/${c.player_id}`}
                  className="rounded-full border border-border bg-panel px-2.5 py-0.5 text-[11px] text-text-muted transition-colors hover:border-border-strong hover:text-text"
                >
                  {c.full_name}
                </Link>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}
