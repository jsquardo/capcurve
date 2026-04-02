const POSITIONS = ['SP', 'RP', 'C', '1B', '2B', '3B', 'SS', 'LF', 'CF', 'RF', 'DH']

const SORT_OPTIONS = [
  { value: 'name',         label: 'Name A–Z'      },
  { value: '-name',        label: 'Name Z–A'      },
  { value: '-value_score', label: 'Highest Value' },
  { value: '-recent_year', label: 'Most Recent'   },
]

const ACTIVE_TABS: { label: string; value: boolean | undefined }[] = [
  { label: 'All',     value: undefined },
  { label: 'Active',  value: true      },
  { label: 'Retired', value: false     },
]

const SELECT_CLASS =
  'rounded-[6px] border border-border bg-elevated px-3 py-1.5 text-[12px] text-text-muted focus:border-border-strong focus:outline-none'

interface PlayerFiltersProps {
  q: string
  active: boolean | undefined
  position: string
  sort: string
  onQChange: (q: string) => void
  onActiveChange: (active: boolean | undefined) => void
  onPositionChange: (position: string) => void
  onSortChange: (sort: string) => void
}

export default function PlayerFilters({
  q,
  active,
  position,
  sort,
  onQChange,
  onActiveChange,
  onPositionChange,
  onSortChange,
}: PlayerFiltersProps) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      {/* Text search */}
      <input
        type="search"
        value={q}
        onChange={e => onQChange(e.target.value)}
        placeholder="Search players…"
        className="min-w-[200px] flex-1 rounded-[6px] border border-border bg-elevated px-3 py-1.5 text-[12px] text-text placeholder:text-text-subtle focus:border-border-strong focus:outline-none"
      />

      {/* Active / Retired pills */}
      <div className="flex items-center rounded-[6px] border border-border bg-elevated">
        {ACTIVE_TABS.map(tab => (
          <button
            key={String(tab.value)}
            onClick={() => onActiveChange(tab.value)}
            className={`px-3 py-1.5 text-[12px] font-medium transition-colors first:rounded-l-[5px] last:rounded-r-[5px] ${
              active === tab.value
                ? 'bg-accent/10 text-accent'
                : 'text-text-muted hover:text-text'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Position */}
      <select
        value={position}
        onChange={e => onPositionChange(e.target.value)}
        className={SELECT_CLASS}
      >
        <option value="">All Positions</option>
        {POSITIONS.map(pos => (
          <option key={pos} value={pos}>{pos}</option>
        ))}
      </select>

      {/* Sort */}
      <select
        value={sort}
        onChange={e => onSortChange(e.target.value)}
        className={SELECT_CLASS}
      >
        {SORT_OPTIONS.map(opt => (
          <option key={opt.value} value={opt.value}>{opt.label}</option>
        ))}
      </select>
    </div>
  )
}
