function SparkLine({ bars }: { bars: number[] }) {
  const max = Math.max(...bars)
  const n = bars.length
  const w = 100
  const h = 38
  const pts = bars.map((v, i) => `${(i / (n - 1)) * w},${h - (v / max) * h}`)
  const area = pts.join(' ')
  const fillPts = `0,${h} ${area} ${w},${h}`
  const lastY = h - (bars[n - 1] / max) * h

  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} className="overflow-visible">
      <polygon points={fillPts} fill="rgba(240,192,64,0.08)" />
      <polyline
        points={area}
        fill="none"
        stroke="rgb(var(--color-accent))"
        strokeWidth={2}
        strokeLinejoin="round"
        strokeLinecap="round"
      />
      <circle cx={w} cy={lastY} r={3} fill="rgb(var(--color-accent))" />
    </svg>
  )
}

interface TrendingCardProps {
  rank: number
  name: string
  team: string
  delta: string
  bars: number[]
  isViews?: boolean
}

export default function TrendingCard({ rank, name, team, delta, bars, isViews = false }: TrendingCardProps) {
  return (
    <div className="relative cursor-pointer overflow-hidden rounded-xl border border-border bg-elevated p-5 transition-all duration-200 hover:-translate-y-0.5 hover:border-accent">
      <div className="mb-3 font-mono text-[10px] font-semibold uppercase tracking-[1px] text-text-subtle">
        #{rank} · 14 day {isViews ? 'views' : 'arc move'}
      </div>
      <div className="mb-[2px] text-[15px] font-semibold text-text">{name}</div>
      <div className="mb-4 text-[11px] text-text-subtle">{team}</div>
      <div className="mb-3">
        <SparkLine bars={bars} />
      </div>
      <div className="inline-flex items-center gap-[5px] rounded-[6px] bg-success/[0.12] px-[10px] py-1 font-mono text-[13px] font-medium text-success">
        {delta}
      </div>
    </div>
  )
}
