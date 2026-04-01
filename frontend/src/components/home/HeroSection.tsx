import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import {
  Area,
  CartesianGrid,
  ComposedChart,
  Line,
  ReferenceArea,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

// Mock data — replace with featured-player API response when backend endpoint is ready
// TODO: wire to a /api/v1/players/featured (or similar) endpoint in Phase 3
const heroArcData = [
  { year: '2014', actual: 48, projected: null, lower: null, upper: null },
  { year: '2015', actual: 56, projected: null, lower: null, upper: null },
  { year: '2016', actual: 67, projected: null, lower: null, upper: null },
  { year: '2017', actual: 78, projected: null, lower: null, upper: null },
  { year: '2018', actual: 86, projected: null, lower: null, upper: null },
  { year: '2019', actual: 91, projected: null, lower: null, upper: null },
  { year: '2020', actual: 88, projected: null, lower: null, upper: null },
  { year: '2021', actual: 84, projected: null, lower: null, upper: null },
  { year: '2022', actual: 80, projected: null, lower: null, upper: null },
  { year: '2023', actual: 83, projected: null, lower: null, upper: null },
  { year: '2024', actual: 79, projected: 79, lower: 74, upper: 84 },
  { year: '2025', actual: null, projected: 76, lower: 68, upper: 84 },
  { year: '2026', actual: null, projected: 72, lower: 62, upper: 82 },
  { year: '2027', actual: null, projected: 67, lower: 55, upper: 79 },
  { year: '2028', actual: null, projected: 61, lower: 48, upper: 74 },
]

const featureStats = [
  { label: 'Peak Arc Score', value: '91.4', tone: 'accent' as const },
  { label: 'Current Score', value: '79.0', tone: 'default' as const },
  { label: 'Projection Arc', value: 'Age 35', tone: 'danger' as const },
] as const

type FeatureStatTone = (typeof featureStats)[number]['tone']

const statToneClassName: Record<FeatureStatTone, string> = {
  accent: 'text-accent',
  default: 'text-text',
  danger: 'text-danger',
}

function HeroTooltip({
  active,
  payload,
  label,
}: {
  active?: boolean
  payload?: Array<{ dataKey?: string; value?: number | null }>
  label?: string
}) {
  if (!active || !payload?.length) return null

  const actualPoint = payload.find((e) => e.dataKey === 'actual' && typeof e.value === 'number')
  const projectionPoint = payload.find((e) => e.dataKey === 'projected' && typeof e.value === 'number')
  const upperPoint = payload.find((e) => e.dataKey === 'upper' && typeof e.value === 'number')
  const lowerPoint = payload.find((e) => e.dataKey === 'lower' && typeof e.value === 'number')

  return (
    <div className="rounded-xl border border-border/80 bg-app/95 px-4 py-3 shadow-[var(--shadow-soft)] backdrop-blur-sm">
      <div className="font-mono text-[11px] uppercase tracking-[0.18em] text-text-subtle">{label}</div>
      {typeof actualPoint?.value === 'number' ? (
        <div className="mt-2 flex items-center justify-between gap-4 text-sm text-text">
          <span>Historical value score</span>
          <span className="stat-value text-accent">{actualPoint.value.toFixed(1)}</span>
        </div>
      ) : null}
      {typeof projectionPoint?.value === 'number' ? (
        <div className="mt-2 flex items-center justify-between gap-4 text-sm text-text">
          <span>Projected value score</span>
          <span className="stat-value text-projection">{projectionPoint.value.toFixed(1)}</span>
        </div>
      ) : null}
      {typeof lowerPoint?.value === 'number' && typeof upperPoint?.value === 'number' ? (
        <div className="mt-2 flex items-center justify-between gap-4 text-sm text-text-muted">
          <span>Confidence band</span>
          <span className="stat-value">
            {lowerPoint.value.toFixed(1)}-{upperPoint.value.toFixed(1)}
          </span>
        </div>
      ) : null}
    </div>
  )
}

export default function HeroSection() {
  return (
    <section className="relative min-h-[520px] overflow-hidden border-b border-border">
      <div className="pointer-events-none absolute -right-20 top-[-100px] h-[700px] w-[700px] rounded-full bg-[radial-gradient(circle,rgba(240,192,64,0.08)_0%,rgba(240,192,64,0.03)_35%,transparent_70%)]" />
      <div className="pointer-events-none absolute bottom-[-150px] left-[30%] h-[500px] w-[500px] rounded-full bg-[radial-gradient(circle,rgba(74,158,255,0.05)_0%,transparent_60%)]" />

      <div className="relative grid min-h-[520px] grid-cols-1 gap-0 xl:grid-cols-[1fr_1fr]">
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: 'easeOut' }}
          className="flex flex-col justify-center border-b border-border px-10 py-16 xl:border-b-0 xl:border-r xl:py-16 xl:pr-12 xl:pl-10"
        >
          <div className="mb-6 flex items-center gap-[10px]">
            <span className="h-[2px] w-6 bg-accent" />
            <span className="text-[11px] font-semibold uppercase tracking-[2px] text-accent">
              MLB Career Intelligence
            </span>
          </div>

          <h1 className="mb-7 font-display text-[86px] leading-[0.92] tracking-[2px] text-text max-xl:text-[68px] max-sm:text-[56px]">
            <span className="block">EVERY</span>
            <span className="block text-accent">CAREER.</span>
            <span className="block text-transparent [-webkit-text-stroke:1px_rgb(var(--color-text-subtle))]">
              CHARTED.
            </span>
          </h1>

          <p className="mb-8 max-w-[400px] text-[15px] leading-[1.75] text-text-muted">
            Deep career arc analysis for every MLB player: peak years, projection curves, and the comparable
            legends who came before them. Built for fans who want more than a box score.
          </p>

          <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
            <Link
              to="/players/660271"
              className="inline-flex items-center justify-center rounded-[8px] bg-accent px-7 py-3 text-[14px] font-semibold tracking-[0.2px] text-[#0a0d12] transition hover:bg-accent-strong"
            >
              Explore Players
            </Link>
            <Link
              to="/leaderboards"
              className="inline-flex items-center justify-center rounded-[8px] border border-border-strong bg-transparent px-6 py-3 text-[14px] font-medium text-text-muted transition hover:border-text-muted hover:text-text"
            >
              Open Leaderboards
            </Link>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.55, ease: 'easeOut', delay: 0.05 }}
          className="relative bg-elevated"
        >
          <div className="flex h-full flex-col p-8">
            <div className="mb-3 text-[10px] font-semibold uppercase tracking-[2px] text-text-subtle">
              Featured Arc · Today
            </div>
            <div className="mb-1 font-display text-[44px] leading-none tracking-[2px] text-accent">
              MOOKIE BETTS
            </div>
            <div className="mb-5 text-[12px] text-text-muted">
              RF · Los Angeles Dodgers · 2014-Present
            </div>

            <div className="flex-1">
              <div className="relative h-[260px] min-h-[200px] sm:h-[320px]">
                <ResponsiveContainer width="100%" height="100%">
                  <ComposedChart data={heroArcData} margin={{ top: 8, right: 8, bottom: 8, left: -24 }}>
                    <defs>
                      <linearGradient id="heroProjectionBand" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="rgb(var(--color-projection))" stopOpacity={0.24} />
                        <stop offset="100%" stopColor="rgb(var(--color-projection))" stopOpacity={0.03} />
                      </linearGradient>
                      <linearGradient id="heroHistoricalLine" x1="0" y1="0" x2="1" y2="0">
                        <stop offset="0%" stopColor="rgb(var(--color-accent))" stopOpacity={0.82} />
                        <stop offset="100%" stopColor="rgb(var(--color-accent-strong))" stopOpacity={1} />
                      </linearGradient>
                    </defs>

                    <CartesianGrid stroke="rgb(var(--color-border))" strokeDasharray="3 6" vertical={false} />
                    <ReferenceArea
                      x1="2018"
                      x2="2020"
                      fill="rgb(var(--color-accent))"
                      fillOpacity={0.08}
                      stroke="rgb(var(--color-accent))"
                      strokeOpacity={0.18}
                    />
                    <XAxis
                      dataKey="year"
                      tickLine={false}
                      axisLine={false}
                      minTickGap={20}
                      tick={{ fill: 'rgb(var(--color-text-subtle))', fontSize: 11, fontFamily: 'var(--font-mono)' }}
                    />
                    <YAxis
                      domain={[40, 100]}
                      tickCount={7}
                      tickLine={false}
                      axisLine={false}
                      tick={{ fill: 'rgb(var(--color-text-subtle))', fontSize: 11, fontFamily: 'var(--font-mono)' }}
                    />
                    <Tooltip
                      content={<HeroTooltip />}
                      cursor={{ stroke: 'rgb(var(--color-border-strong))', strokeDasharray: '4 4' }}
                    />
                    <Area
                      type="monotone"
                      dataKey="upper"
                      stroke="none"
                      fill="url(#heroProjectionBand)"
                      activeDot={false}
                      isAnimationActive={false}
                      legendType="none"
                    />
                    <Area
                      type="monotone"
                      dataKey="lower"
                      stroke="none"
                      fill="rgb(var(--color-bg-elevated))"
                      fillOpacity={1}
                      activeDot={false}
                      isAnimationActive={false}
                      legendType="none"
                    />
                    <Line
                      type="monotone"
                      dataKey="actual"
                      stroke="url(#heroHistoricalLine)"
                      strokeWidth={3}
                      dot={false}
                      activeDot={{ r: 4, fill: 'rgb(var(--color-accent))', strokeWidth: 0 }}
                      animationDuration={900}
                    />
                    <Line
                      type="monotone"
                      dataKey="projected"
                      stroke="rgb(var(--color-projection))"
                      strokeWidth={3}
                      strokeDasharray="7 7"
                      dot={false}
                      activeDot={{ r: 4, fill: 'rgb(var(--color-projection))', strokeWidth: 0 }}
                      animationDuration={900}
                    />
                  </ComposedChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="mt-5 grid overflow-hidden rounded-[8px] border border-border bg-border sm:grid-cols-3">
              {featureStats.map((stat) => (
                <div key={stat.label} className="border-r border-border bg-panel px-4 py-3 last:border-r-0">
                  <div className={`font-display text-3xl tracking-[0.08em] ${statToneClassName[stat.tone]}`}>
                    {stat.value}
                  </div>
                  <div className="mt-[2px] text-[10px] tracking-[0.3px] text-text-subtle">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
