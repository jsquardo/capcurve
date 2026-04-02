interface PlayerListSkeletonProps {
  rows?: number
}

export default function PlayerListSkeleton({ rows = 10 }: PlayerListSkeletonProps) {
  return (
    <div className="space-y-1.5 animate-pulse">
      {Array.from({ length: rows }).map((_, i) => (
        <div
          key={i}
          className="flex items-center gap-3 rounded-[8px] border border-border bg-elevated px-[14px] py-[10px]"
        >
          {/* Avatar circle */}
          <div className="h-9 w-9 shrink-0 rounded-full bg-border" />

          {/* Name + team */}
          <div className="flex-1 space-y-1.5">
            <div className="h-[13px] w-2/5 rounded-sm bg-border" />
            <div className="h-[10px] w-1/4 rounded-sm bg-border opacity-60" />
          </div>

          {/* Badge */}
          <div className="h-[20px] w-14 shrink-0 rounded-full bg-border" />

          {/* Value */}
          <div className="w-[52px] space-y-1.5 text-right">
            <div className="ml-auto h-[14px] w-8 rounded-sm bg-border" />
            <div className="ml-auto h-[10px] w-6 rounded-sm bg-border opacity-60" />
          </div>
        </div>
      ))}
    </div>
  )
}
