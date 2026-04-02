interface PlayerListHeroProps {
  total: number | null
}

export default function PlayerListHero({ total }: PlayerListHeroProps) {
  return (
    <div className="border-b border-border pb-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="font-display text-5xl tracking-[0.04em] text-text">PLAYERS</h1>
          <p className="mt-1 text-[13px] text-text-muted">
            Browse and search the full CapCurve player database
          </p>
        </div>
        {total === null ? (
          <div className="h-[26px] w-24 animate-pulse rounded-full bg-border" />
        ) : (
          <span className="rounded-full border border-border bg-elevated px-3 py-1 font-mono text-[11px] text-text-subtle">
            {total.toLocaleString()} PLAYERS
          </span>
        )}
      </div>
    </div>
  )
}
