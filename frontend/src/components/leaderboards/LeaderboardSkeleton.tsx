interface LeaderboardSkeletonProps {
  rows?: number
}

export default function LeaderboardSkeleton({ rows = 5 }: LeaderboardSkeletonProps) {
  return (
    <div className="space-y-1.5 animate-pulse">
      {Array.from({ length: rows }).map((_, i) => (
        <div
          key={i}
          className="flex items-center gap-3 rounded-[8px] border border-border bg-elevated px-[14px] py-[10px]"
        >
          {/* Rank badge */}
          <div className="w-[18px] h-[11px] shrink-0 rounded-sm bg-border" />

          {/* Player name + team */}
          <div className="flex-1 space-y-1.5">
            <div className="h-[13px] w-2/5 rounded-sm bg-border" />
            <div className="h-[10px] w-1/3 rounded-sm bg-border opacity-60" />
          </div>

          {/* Bar */}
          <div className="h-1 w-24 shrink-0 rounded-sm bg-border" />

          {/* Value chip */}
          <div className="w-[52px] h-[14px] shrink-0 rounded-sm bg-border" />
        </div>
      ))}
    </div>
  )
}
