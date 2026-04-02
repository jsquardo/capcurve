interface LeaderboardHeroProps {
  season: number
}

export default function LeaderboardHero({ season }: LeaderboardHeroProps) {
  return (
    <div className="border-b border-border pb-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="font-display text-5xl tracking-[0.04em] text-text">LEADERBOARDS</h1>
          <p className="mt-1 text-[13px] text-text-muted">
            Rankings across career arc peaks and seasonal stat categories
          </p>
        </div>
        <span className="rounded-full border border-border bg-elevated px-3 py-1 font-mono text-[11px] text-text-subtle">
          {season} SEASON
        </span>
      </div>
    </div>
  )
}
