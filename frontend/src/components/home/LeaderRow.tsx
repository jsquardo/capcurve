interface LeaderRowProps {
  rank: number
  name: string
  team: string
  value: number | string
  barPct: number
}

export default function LeaderRow({ rank, name, team, value, barPct }: LeaderRowProps) {
  return (
    <div className="flex items-center gap-3 rounded-[8px] border border-border bg-elevated px-[14px] py-[10px]">
      <span className="w-[18px] shrink-0 text-right font-mono text-[11px] text-text-subtle">
        {rank}
      </span>
      <div className="flex-1">
        <div className="text-[13px] font-medium">{name}</div>
        <div className="text-[10px] text-text-subtle">{team}</div>
      </div>
      <div className="h-1 w-24 shrink-0 overflow-hidden rounded-sm bg-panel">
        <div
          className="h-full rounded-sm bg-accent transition-all duration-300"
          style={{ width: `${barPct}%` }}
        />
      </div>
      <span className="w-[52px] shrink-0 text-right font-mono text-[14px] font-medium text-accent">
        {value}
      </span>
    </div>
  )
}
