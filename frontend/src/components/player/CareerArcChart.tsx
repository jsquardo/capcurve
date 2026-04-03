import { useState } from 'react'
import {
  ComposedChart,
  Line,
  Area,
  ReferenceArea,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { CareerArcData } from '@/types'

// ── Color constants (raw values from CSS variables) ───────────────────────────
const COLOR_ACCENT = '#f0c040'
const COLOR_PROJECTION = '#9b72f8'
const COLOR_BAND_FILL = 'rgba(155,114,248,0.15)'
const COLOR_PEAK_FILL = 'rgba(240,192,64,0.08)'
const COLOR_GRID = 'rgba(37,45,66,0.7)'
const COLOR_AXIS_TEXT = 'rgb(107,122,153)'

// ── Merged chart point ─────────────────────────────────────────────────────────

interface ChartPoint {
  year: number
  age: number | null
  team_name: string | null
  // Historical arc score (null on pure projection points)
  historical: number | null
  // Projected arc score (null on pure historical points)
  projected: number | null
  // Confidence band — only populated for projection points
  lower: number | null
  upper: number | null
  is_peak: boolean
  is_projection: boolean
}

function buildChartData(arcData: CareerArcData): ChartPoint[] {
  const historicalPoints: ChartPoint[] = arcData.timeline
    .filter(t => !t.is_projection)
    .map(t => ({
      year: t.year,
      age: t.age,
      team_name: t.team_name,
      historical: t.value_score,
      projected: null,
      lower: null,
      upper: null,
      is_peak: t.is_peak,
      is_projection: false,
    }))

  const { points: projPoints, confidence_band } = arcData.projection

  // Build a map of year → confidence band for quick lookup
  const bandMap = new Map(confidence_band.map(b => [b.year, b]))

  const projectionPoints: ChartPoint[] = projPoints.map(p => ({
    year: p.year,
    age: p.age,
    team_name: null,
    historical: null,
    projected: p.value_score,
    lower: bandMap.get(p.year)?.lower ?? null,
    upper: bandMap.get(p.year)?.upper ?? null,
    is_peak: false,
    is_projection: true,
  }))

  // Duplicate the last historical point as the start of the projection series so the
  // gold line and purple dashed line share a seamless join at the handoff year.
  const lastHistorical = historicalPoints[historicalPoints.length - 1]
  if (lastHistorical && projectionPoints.length > 0) {
    projectionPoints.unshift({
      ...lastHistorical,
      historical: null,
      projected: lastHistorical.historical,
      lower: projectionPoints[0].lower,
      upper: projectionPoints[0].upper,
      is_projection: true,
    })
  }

  const allPoints = [...historicalPoints, ...projectionPoints]
  allPoints.sort((a, b) => a.year - b.year)
  return allPoints
}

// ── Custom tooltip ─────────────────────────────────────────────────────────────

interface TooltipPayloadEntry {
  dataKey: string
  value: number
  payload: ChartPoint
}

interface CustomTooltipProps {
  active?: boolean
  payload?: TooltipPayloadEntry[]
  label?: number
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (!active || !payload || payload.length === 0) return null
  const point = payload[0].payload as ChartPoint

  const score = point.is_projection ? point.projected : point.historical

  return (
    <div className="rounded-[8px] border border-border bg-overlay px-3 py-2.5 shadow-lg">
      <div className="mb-1 flex items-center gap-2">
        <span className="font-mono text-[11px] text-text-subtle">{point.year}</span>
        {point.age !== null && (
          <span className="font-mono text-[10px] text-text-subtle">Age {point.age}</span>
        )}
        {point.is_peak && (
          <span className="rounded px-1.5 py-0.5 font-mono text-[9px] font-medium"
            style={{ background: COLOR_PEAK_FILL, color: COLOR_ACCENT, border: `1px solid rgba(240,192,64,0.3)` }}>
            PEAK
          </span>
        )}
        {point.is_projection && (
          <span className="rounded px-1.5 py-0.5 font-mono text-[9px] font-medium"
            style={{ background: 'rgba(155,114,248,0.12)', color: COLOR_PROJECTION, border: `1px solid rgba(155,114,248,0.3)` }}>
            PROJ
          </span>
        )}
      </div>

      {score !== null && (
        <div
          className="font-display text-[26px] leading-none"
          style={{ color: point.is_projection ? COLOR_PROJECTION : COLOR_ACCENT }}
        >
          {score.toFixed(1)}
        </div>
      )}

      {point.is_projection && point.lower !== null && point.upper !== null && (
        <div className="mt-1 font-mono text-[10px]" style={{ color: COLOR_PROJECTION }}>
          {point.lower.toFixed(1)}–{point.upper.toFixed(1)} range
        </div>
      )}

      {!point.is_projection && point.team_name && (
        <div className="mt-1 text-[11px] text-text-subtle">{point.team_name}</div>
      )}
    </div>
  )
}

// ── Toggle button ──────────────────────────────────────────────────────────────

interface ToggleProps {
  label: string
  on: boolean
  variant: 'gold' | 'proj'
  onClick: () => void
}

function Toggle({ label, on, variant, onClick }: ToggleProps) {
  const activeClass =
    variant === 'gold'
      ? 'border-accent text-accent bg-accent/[0.08]'
      : 'border-projection text-projection bg-projection/[0.08]'
  const inactiveClass = 'border-border text-text-subtle'

  return (
    <button
      onClick={onClick}
      className={`rounded-[6px] border px-3 py-1 text-xs font-medium transition-colors ${on ? activeClass : inactiveClass}`}
    >
      {label}
    </button>
  )
}

// ── Legend swatch ──────────────────────────────────────────────────────────────

function LegendItem({ children, swatch }: { children: React.ReactNode; swatch: React.ReactNode }) {
  return (
    <div className="flex items-center gap-1.5 text-[11px] text-text-subtle">
      {swatch}
      {children}
    </div>
  )
}

// ── Main component ─────────────────────────────────────────────────────────────

interface CareerArcChartProps {
  arcData: CareerArcData
}

export default function CareerArcChart({ arcData }: CareerArcChartProps) {
  const [showPeak, setShowPeak] = useState(true)
  const [showProjection, setShowProjection] = useState(true)

  const data = buildChartData(arcData)

  // Deduplicate years for the X axis tick array — the join point intentionally
  // duplicates the last historical year as the first projection point, which would
  // produce a double label if we let Recharts derive ticks from the data itself.
  const uniqueYearTicks = [...new Set(data.map(d => d.year))]

  // On narrow viewports, showing every year for a long career crowds the axis.
  // Show every other tick when there are more than 8 unique years, but always
  // keep the first and last year so the line doesn't appear to start from nowhere.
  const displayTicks =
    uniqueYearTicks.length > 8
      ? [
          uniqueYearTicks[0],
          ...uniqueYearTicks.slice(1, -1).filter((_, i) => i % 2 === 0),
          uniqueYearTicks[uniqueYearTicks.length - 1],
        ]
      : uniqueYearTicks

  if (data.length === 0) {
    return (
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-10 text-center text-sm text-text-subtle">
        No arc data available yet.
      </div>
    )
  }

  const arcMeta = arcData.arc
  const projEligible = arcData.projection.eligible && arcData.projection.status === 'ready'
  // Don't show the confidence band when all projected scores are 0 — this happens when
  // the comparable pool is empty/broken and produces meaningless band data.
  const projHasData = arcData.projection.points.some(p => p.value_score > 0)
  const projVisible = projEligible && showProjection
  const bandVisible = projVisible && projHasData

  return (
    <div className="rounded-[8px] border border-border bg-elevated px-5 py-5 sm:px-6 sm:py-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-display text-2xl tracking-[0.04em] text-text">
          Career <span className="text-accent">Arc</span>
        </h2>
        <div className="flex gap-2">
          {arcMeta && (
            <Toggle
              label="Peak Window"
              on={showPeak}
              variant="gold"
              onClick={() => setShowPeak(p => !p)}
            />
          )}
          {projEligible && (
            <Toggle
              label="Projection"
              on={showProjection}
              variant="proj"
              onClick={() => setShowProjection(p => !p)}
            />
          )}
        </div>
      </div>

      {/* Chart */}
      <div className="h-[220px] md:h-[320px]">
        <ResponsiveContainer width="100%" height="100%">
          <ComposedChart data={data} margin={{ top: 8, right: 4, bottom: 0, left: 0 }}>
            <CartesianGrid
              vertical={false}
              stroke={COLOR_GRID}
              strokeDasharray="4 4"
            />
            <XAxis
              dataKey="year"
              type="number"
              domain={['dataMin', 'dataMax']}
              ticks={displayTicks}
              tick={{ fill: COLOR_AXIS_TEXT, fontFamily: 'DM Mono', fontSize: 11 }}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              domain={[0, 100]}
              ticks={[0, 25, 50, 75, 100]}
              tick={{ fill: COLOR_AXIS_TEXT, fontFamily: 'DM Mono', fontSize: 10 }}
              tickLine={false}
              axisLine={false}
              width={28}
            />
            <Tooltip
              content={<CustomTooltip />}
              cursor={{ stroke: COLOR_GRID, strokeWidth: 1 }}
            />

            {/* Peak window shading */}
            {arcMeta && showPeak && (
              <ReferenceArea
                x1={arcMeta.peak_year_start}
                x2={arcMeta.peak_year_end}
                fill={COLOR_PEAK_FILL}
                stroke="none"
              />
            )}

            {/* Confidence band — rendered behind the projection line */}
            {bandVisible && (
              <Area
                dataKey="upper"
                fill={COLOR_BAND_FILL}
                stroke="none"
                isAnimationActive={false}
                connectNulls
              />
            )}
            {bandVisible && (
              <Area
                dataKey="lower"
                fill="rgba(10,13,18,1)"
                stroke="none"
                isAnimationActive={false}
                connectNulls
              />
            )}

            {/* Historical arc score line */}
            <Line
              dataKey="historical"
              stroke={COLOR_ACCENT}
              strokeWidth={2.5}
              dot={false}
              activeDot={{ r: 4, fill: COLOR_ACCENT, stroke: 'rgba(10,13,18,1)', strokeWidth: 2 }}
              connectNulls
              isAnimationActive
              animationDuration={800}
              animationEasing="ease-out"
            />

            {/* Projected arc score line */}
            {projVisible && (
              <Line
                dataKey="projected"
                stroke={COLOR_PROJECTION}
                strokeWidth={2}
                strokeDasharray="6 4"
                dot={false}
                activeDot={{ r: 4, fill: COLOR_PROJECTION, stroke: 'rgba(10,13,18,1)', strokeWidth: 2 }}
                connectNulls
                isAnimationActive={false}
              />
            )}
          </ComposedChart>
        </ResponsiveContainer>
      </div>

      {/* Legend */}
      <div className="mt-4 flex flex-wrap gap-x-4 gap-y-2 border-t border-border pt-3">
        <LegendItem
          swatch={<div className="h-[3px] w-[18px] rounded-full" style={{ background: COLOR_ACCENT }} />}
        >
          Arc Score
        </LegendItem>
        {projVisible && (
          <>
            <LegendItem
              swatch={
                <div className="w-[18px] border-t-2 border-dashed" style={{ borderColor: COLOR_PROJECTION }} />
              }
            >
              Projection
            </LegendItem>
            {bandVisible && (
              <LegendItem
                swatch={<div className="h-3 w-[18px] rounded-sm" style={{ background: COLOR_BAND_FILL, border: `1px solid rgba(155,114,248,0.3)` }} />}
              >
                Conf. Band
              </LegendItem>
            )}
          </>
        )}
        {arcMeta && showPeak && (
          <LegendItem
            swatch={<div className="h-3 w-[18px] rounded-sm" style={{ background: COLOR_PEAK_FILL, border: `1px solid rgba(240,192,64,0.2)` }} />}
          >
            Peak Window
          </LegendItem>
        )}
      </div>
    </div>
  )
}
