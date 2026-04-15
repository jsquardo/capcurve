type SortOption = { label: string; value: string; group?: 'hitting' | 'pitching' }

const SORT_OPTIONS: SortOption[] = [
  { label: 'Arc Score ↓',  value: '-value_score'     },
  { label: 'Arc Score ↑',  value: 'value_score'      },
  { label: 'HR ↓',         value: '-home_runs',        group: 'hitting' },
  { label: 'AVG ↓',        value: '-batting_avg',      group: 'hitting' },
  { label: 'OBP ↓',        value: '-obp',              group: 'hitting' },
  { label: 'SLG ↓',        value: '-slg',              group: 'hitting' },
  { label: 'ERA ↓',        value: '-era',              group: 'pitching' },
  { label: 'K/9 ↓',        value: '-strikeouts_per_9', group: 'pitching' },
  { label: 'IP ↓',         value: '-innings_pitched',  group: 'pitching' },
  { label: 'Year ↑',       value: 'year'              },
  { label: 'Year ↓',       value: '-year'             },
]

interface Props {
  total: number
  activeFilterCount: number // computed by PlaygroundPage, passed in
  sort: string
  onSortChange: (sort: string) => void
  group?: 'all' | 'hitting' | 'pitching'
}

export default function PlaygroundResultsHeader({ total, activeFilterCount, sort, onSortChange, group }: Props) {
  // Hide sort options that are meaningless for the active group. Universal options
  // (no group tag) are always shown. Hitting-only options are hidden for pitching
  // results and vice versa.
  const visibleOptions = SORT_OPTIONS.filter((opt) => {
    if (!opt.group) return true
    if (group === 'hitting') return opt.group === 'hitting'
    if (group === 'pitching') return opt.group === 'pitching'
    return true // group === 'all' or undefined — show everything
  })

  return (
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center gap-2">
        <span className="text-[14px] font-semibold text-text">
          {total.toLocaleString()} result{total !== 1 ? 's' : ''}
        </span>
        {activeFilterCount > 0 && (
          <span className="rounded-full border border-border bg-elevated px-2 py-0.5 text-[11px] font-medium text-text-subtle">
            {activeFilterCount} filter{activeFilterCount !== 1 ? 's' : ''}
          </span>
        )}
      </div>

      <div className="flex items-center gap-2">
        <span className="text-[11px] text-text-subtle">Sort</span>
        <select
          value={sort}
          onChange={(e) => onSortChange(e.target.value)}
          className="shell-input py-1 pr-7 text-[12px]"
        >
          {visibleOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      </div>
    </div>
  )
}
