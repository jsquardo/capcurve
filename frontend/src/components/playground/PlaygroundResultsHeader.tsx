const SORT_OPTIONS: { label: string; value: string }[] = [
  { label: 'Arc Score ↓',  value: '-value_score'     },
  { label: 'Arc Score ↑',  value: 'value_score'      },
  { label: 'HR ↓',         value: '-home_runs'        },
  { label: 'AVG ↓',        value: '-batting_avg'      },
  { label: 'OBP ↓',        value: '-obp'              },
  { label: 'SLG ↓',        value: '-slg'              },
  { label: 'ERA ↓',        value: '-era'              },
  { label: 'K/9 ↓',        value: '-strikeouts_per_9' },
  { label: 'IP ↓',         value: '-innings_pitched'  },
  { label: 'Year ↑',       value: 'year'              },
  { label: 'Year ↓',       value: '-year'             },
]

interface Props {
  total: number
  activeFilterCount: number // computed by PlaygroundPage, passed in
  sort: string
  onSortChange: (sort: string) => void
}

export default function PlaygroundResultsHeader({ total, activeFilterCount, sort, onSortChange }: Props) {
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
          {SORT_OPTIONS.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      </div>
    </div>
  )
}
