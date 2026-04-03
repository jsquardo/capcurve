export default function PlayerPageSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">

      {/* PlayerHero block */}
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-6">
        <div className="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between">
          {/* Avatar + name */}
          <div className="flex items-center gap-4">
            <div className="h-20 w-20 shrink-0 rounded-full bg-border" />
            <div className="space-y-2.5">
              <div className="h-7 w-48 rounded-sm bg-border" />
              <div className="h-4 w-32 rounded-sm bg-border opacity-60" />
              <div className="h-3 w-24 rounded-sm bg-border opacity-40" />
            </div>
          </div>

          {/* Stat summary cards */}
          <div className="flex gap-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                key={i}
                className="flex h-14 w-20 shrink-0 flex-col justify-center rounded-[6px] border border-border bg-surface px-3 space-y-1.5 sm:w-24"
              >
                <div className="h-[10px] w-10 rounded-sm bg-border opacity-50" />
                <div className="h-[14px] w-8 rounded-sm bg-border" />
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* CareerArcChart block */}
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-5">
        <div className="mb-4 h-5 w-36 rounded-sm bg-border" />
        <div className="h-72 w-full rounded-[6px] bg-border opacity-40 md:h-96" />
      </div>

      {/* ProjectionPanel block */}
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-5">
        <div className="mb-4 h-5 w-40 rounded-sm bg-border" />
        <div className="space-y-2.5">
          <div className="h-3 w-3/4 rounded-sm bg-border opacity-50" />
          <div className="h-3 w-1/2 rounded-sm bg-border opacity-40" />
          <div className="mt-4 h-24 w-full rounded-[6px] bg-border opacity-30" />
        </div>
      </div>

      {/* SeasonStatsTable block */}
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-5">
        <div className="mb-4 h-5 w-32 rounded-sm bg-border" />
        {/* Header row */}
        <div className="mb-2 flex gap-4 border-b border-border pb-2">
          {[40, 56, 72, 48, 48, 48].map((w, i) => (
            <div key={i} className="h-3 rounded-sm bg-border opacity-50" style={{ width: w }} />
          ))}
        </div>
        {/* Data rows */}
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex gap-4 border-b border-border py-2.5 last:border-0">
            {[40, 56, 72, 48, 48, 48].map((w, j) => (
              <div key={j} className="h-3 rounded-sm bg-border opacity-40" style={{ width: w }} />
            ))}
          </div>
        ))}
      </div>

      {/* ComparablePlayersRow block */}
      <div className="rounded-[8px] border border-border bg-elevated px-6 py-5">
        <div className="mb-4 h-5 w-44 rounded-sm bg-border" />
        <div className="flex gap-3 overflow-hidden">
          {Array.from({ length: 4 }).map((_, i) => (
            <div
              key={i}
              className="flex h-20 w-36 shrink-0 items-center gap-3 rounded-[6px] border border-border bg-surface px-3"
            >
              <div className="h-10 w-10 shrink-0 rounded-full bg-border" />
              <div className="space-y-1.5">
                <div className="h-3 w-16 rounded-sm bg-border" />
                <div className="h-2.5 w-10 rounded-sm bg-border opacity-60" />
              </div>
            </div>
          ))}
        </div>
      </div>

    </div>
  )
}
