import { useEffect, useState } from 'react'
import type { ChangeEvent } from 'react'
import type { PlaygroundQueryParams } from '@/types'
import PlaygroundFilterGroup from './PlaygroundFilterGroup'

// Values must match MLB Stats API primaryPosition.name strings exactly — same list
// as PlayerFilters.tsx. Do not use abbreviations like "1B" or "SP" as values.
const POSITIONS: { label: string; value: string }[] = [
  { label: 'Pitcher',   value: 'Pitcher'           },
  { label: 'SP',        value: 'Starting Pitcher'  },
  { label: 'RP',        value: 'Relief Pitcher'    },
  { label: 'C',         value: 'Catcher'           },
  { label: '1B',        value: 'First Base'        },
  { label: '2B',        value: 'Second Base'       },
  { label: '3B',        value: 'Third Base'        },
  { label: 'SS',        value: 'Shortstop'         },
  { label: 'OF',        value: 'Outfielder'        },
  { label: 'LF',        value: 'Left Field'        },
  { label: 'CF',        value: 'Center Field'      },
  { label: 'RF',        value: 'Right Field'       },
  { label: 'DH',        value: 'Designated Hitter' },
  { label: 'Two-Way',   value: 'Two-Way Player'    },
]

const EMPTY: PlaygroundQueryParams = {}

interface Props {
  // Pre-filled params from URL (?player=:id resolved to q by PlaygroundPage).
  // Changing initialParams resets the draft — handles initial auto-run on load.
  initialParams?: PlaygroundQueryParams
  onSearch: (params: PlaygroundQueryParams) => void
  onReset: () => void
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function optInt(value: string): number | undefined {
  const n = parseInt(value, 10)
  return isNaN(n) ? undefined : n
}

function optFloat(value: string): number | undefined {
  const n = parseFloat(value)
  return isNaN(n) ? undefined : n
}

function strOrUndefined(value: string): string | undefined {
  return value.trim() === '' ? undefined : value.trim()
}

// ── Component ─────────────────────────────────────────────────────────────────

export default function PlaygroundFilterPanel({ initialParams, onSearch, onReset }: Props) {
  const [draft, setDraft] = useState<PlaygroundQueryParams>(initialParams ?? EMPTY)

  // Sync draft when initialParams change (e.g. page load with ?player=:id auto-run)
  useEffect(() => {
    setDraft(initialParams ?? EMPTY)
  }, [initialParams])

  function set<K extends keyof PlaygroundQueryParams>(key: K, value: PlaygroundQueryParams[K]) {
    setDraft((prev) => ({ ...prev, [key]: value }))
  }

  function handlePositionToggle(value: string) {
    const current = draft.position ? draft.position.split(',').map((s) => s.trim()).filter(Boolean) : []
    const next = current.includes(value)
      ? current.filter((p) => p !== value)
      : [...current, value]
    set('position', next.length > 0 ? next.join(',') : undefined)
  }

  function selectedPositions(): string[] {
    return draft.position ? draft.position.split(',').map((s) => s.trim()).filter(Boolean) : []
  }

  // Switching groups clears incompatible stat filters in one atomic update so
  // the backend's validatePlaygroundGroupFilters can never reject the submission.
  function handleGroupChange(g: 'all' | 'hitting' | 'pitching') {
    const group = g === 'all' ? undefined : g
    if (g === 'hitting') {
      setDraft((prev) => ({
        ...prev,
        group,
        min_ip: undefined, max_ip: undefined,
        min_era: undefined, max_era: undefined,
        min_whip: undefined, max_whip: undefined,
        min_k9: undefined, max_k9: undefined,
      }))
    } else if (g === 'pitching') {
      setDraft((prev) => ({
        ...prev,
        group,
        min_pa: undefined, max_pa: undefined,
        min_hr: undefined, max_hr: undefined,
        min_avg: undefined, max_avg: undefined,
        min_obp: undefined, max_obp: undefined,
        min_slg: undefined, max_slg: undefined,
        min_sb: undefined, max_sb: undefined,
      }))
    } else {
      set('group', undefined)
    }
  }

  function handleReset() {
    setDraft(EMPTY)
    onReset()
  }

  const showHitting = draft.group !== 'pitching'
  const showPitching = draft.group !== 'hitting'

  const groupBtnClass = (active: boolean) =>
    `flex-1 rounded-[6px] py-1.5 text-[12px] font-medium transition-colors ${
      active
        ? 'bg-elevated text-text border border-border-strong'
        : 'bg-transparent text-text-subtle hover:text-text'
    }`

  const activeBoolClass = (active: boolean) =>
    `flex-1 rounded-[6px] py-1.5 text-[12px] font-medium transition-colors ${
      active
        ? 'bg-elevated text-text border border-border-strong'
        : 'bg-transparent text-text-subtle hover:text-text'
    }`

  return (
    <div className="flex flex-col gap-0">

      {/* ── Search ── */}
      <PlaygroundFilterGroup label="Search">
        <input
          type="text"
          value={draft.q ?? ''}
          placeholder="Player name…"
          onChange={(e: ChangeEvent<HTMLInputElement>) => set('q', strOrUndefined(e.target.value))}
          className="shell-input w-full"
        />

        {/* Group toggle */}
        <div className="mt-3">
          <p className="mb-1.5 text-[10px] text-text-subtle">Type</p>
          <div className="flex gap-1 rounded-[8px] border border-border bg-panel p-0.5">
            {(['all', 'hitting', 'pitching'] as const).map((g) => (
              <button
                key={g}
                type="button"
                onClick={() => handleGroupChange(g)}
                className={groupBtnClass((draft.group ?? 'all') === g)}
              >
                {g.charAt(0).toUpperCase() + g.slice(1)}
              </button>
            ))}
          </div>
        </div>

        {/* Active toggle */}
        <div className="mt-3">
          <p className="mb-1.5 text-[10px] text-text-subtle">Status</p>
          <div className="flex gap-1 rounded-[8px] border border-border bg-panel p-0.5">
            <button type="button" onClick={() => set('active', undefined)} className={activeBoolClass(draft.active === undefined)}>All</button>
            <button type="button" onClick={() => set('active', true)}      className={activeBoolClass(draft.active === true)}>Active</button>
            <button type="button" onClick={() => set('active', false)}     className={activeBoolClass(draft.active === false)}>Retired</button>
          </div>
        </div>
      </PlaygroundFilterGroup>

      {/* ── Player ── */}
      <PlaygroundFilterGroup label="Player" defaultOpen={false}>
        <div className="mb-3">
          <p className="mb-1.5 text-[10px] text-text-subtle">Team</p>
          <input
            type="text"
            value={draft.team ?? ''}
            placeholder="e.g. Yankees"
            onChange={(e: ChangeEvent<HTMLInputElement>) => set('team', strOrUndefined(e.target.value))}
            className="shell-input w-full"
          />
        </div>
        <p className="mb-1.5 text-[10px] text-text-subtle">Position</p>
        <div className="grid grid-cols-2 gap-x-3 gap-y-1.5">
          {POSITIONS.map(({ label, value }) => (
            <label key={value} className="flex cursor-pointer items-center gap-2 text-[12px] text-text-muted">
              <input
                type="checkbox"
                checked={selectedPositions().includes(value)}
                onChange={() => handlePositionToggle(value)}
                className="accent-accent"
              />
              {label}
            </label>
          ))}
        </div>
      </PlaygroundFilterGroup>

      {/* ── Season / Era ── */}
      <PlaygroundFilterGroup label="Season / Era" defaultOpen={false}>
        <div className="space-y-2">
          {/* Setting a season clears era range (and vice versa) — the backend
              rejects requests that combine season with era_start/era_end. */}
          {/* key changes when the opposing field group is cleared, forcing a
              remount so browsers reliably reflect the empty controlled value.
              type="number" inputs can silently ignore programmatic value=""
              updates without this. */}
          <RangeRow
            key={`season-${draft.era_start !== undefined || draft.era_end !== undefined ? 1 : 0}`}
            label="Season"
            minVal={draft.season?.toString() ?? ''}
            maxVal={''}
            onMinChange={(v) =>
              setDraft((prev) => ({
                ...prev,
                season: optInt(v),
                era_start: undefined,
                era_end: undefined,
              }))
            }
            onMaxChange={() => {}}
            minPlaceholder="e.g. 2023"
            maxPlaceholder=""
            singleValue
          />
          <RangeRow
            key={`era-${draft.season !== undefined ? 1 : 0}`}
            label="Era"
            minVal={draft.era_start?.toString() ?? ''}
            maxVal={draft.era_end?.toString() ?? ''}
            onMinChange={(v) =>
              setDraft((prev) => ({ ...prev, era_start: optInt(v), season: undefined }))
            }
            onMaxChange={(v) =>
              setDraft((prev) => ({ ...prev, era_end: optInt(v), season: undefined }))
            }
            minPlaceholder="From"
            maxPlaceholder="To"
          />
        </div>
      </PlaygroundFilterGroup>

      {/* ── Age ── */}
      <PlaygroundFilterGroup label="Age" defaultOpen={false}>
        <RangeRow
          label=""
          minVal={draft.age_min?.toString() ?? ''}
          maxVal={draft.age_max?.toString() ?? ''}
          onMinChange={(v) => set('age_min', optInt(v))}
          onMaxChange={(v) => set('age_max', optInt(v))}
          minPlaceholder="Min"
          maxPlaceholder="Max"
        />
      </PlaygroundFilterGroup>

      {/* ── Hitting ── */}
      {showHitting && (
        <PlaygroundFilterGroup label="Hitting" defaultOpen={false}>
          <div className="space-y-2">
            <RangeRow label="PA"  minVal={draft.min_pa?.toString() ?? ''}  maxVal={draft.max_pa?.toString() ?? ''}  onMinChange={(v) => set('min_pa',  optInt(v))}   onMaxChange={(v) => set('max_pa',  optInt(v))} />
            <RangeRow label="HR"  minVal={draft.min_hr?.toString() ?? ''}  maxVal={draft.max_hr?.toString() ?? ''}  onMinChange={(v) => set('min_hr',  optInt(v))}   onMaxChange={(v) => set('max_hr',  optInt(v))} />
            <RangeRow label="AVG" minVal={draft.min_avg?.toString() ?? ''} maxVal={draft.max_avg?.toString() ?? ''} onMinChange={(v) => set('min_avg', optFloat(v))} onMaxChange={(v) => set('max_avg', optFloat(v))} minPlaceholder=".200" maxPlaceholder=".400" />
            <RangeRow label="OBP" minVal={draft.min_obp?.toString() ?? ''} maxVal={draft.max_obp?.toString() ?? ''} onMinChange={(v) => set('min_obp', optFloat(v))} onMaxChange={(v) => set('max_obp', optFloat(v))} minPlaceholder=".300" maxPlaceholder=".500" />
            <RangeRow label="SLG" minVal={draft.min_slg?.toString() ?? ''} maxVal={draft.max_slg?.toString() ?? ''} onMinChange={(v) => set('min_slg', optFloat(v))} onMaxChange={(v) => set('max_slg', optFloat(v))} minPlaceholder=".300" maxPlaceholder=".700" />
            <RangeRow label="SB"  minVal={draft.min_sb?.toString() ?? ''}  maxVal={draft.max_sb?.toString() ?? ''}  onMinChange={(v) => set('min_sb',  optInt(v))}   onMaxChange={(v) => set('max_sb',  optInt(v))} />
          </div>
        </PlaygroundFilterGroup>
      )}

      {/* ── Pitching ── */}
      {showPitching && (
        <PlaygroundFilterGroup label="Pitching" defaultOpen={false}>
          <div className="space-y-2">
            <RangeRow label="IP"   minVal={draft.min_ip?.toString() ?? ''}   maxVal={draft.max_ip?.toString() ?? ''}   onMinChange={(v) => set('min_ip',   optFloat(v))} onMaxChange={(v) => set('max_ip',   optFloat(v))} />
            <RangeRow label="ERA"  minVal={draft.min_era?.toString() ?? ''}  maxVal={draft.max_era?.toString() ?? ''}  onMinChange={(v) => set('min_era',  optFloat(v))} onMaxChange={(v) => set('max_era',  optFloat(v))} minPlaceholder="1.00" maxPlaceholder="6.00" />
            <RangeRow label="WHIP" minVal={draft.min_whip?.toString() ?? ''} maxVal={draft.max_whip?.toString() ?? ''} onMinChange={(v) => set('min_whip', optFloat(v))} onMaxChange={(v) => set('max_whip', optFloat(v))} minPlaceholder="0.80" maxPlaceholder="2.00" />
            <RangeRow label="K/9"  minVal={draft.min_k9?.toString() ?? ''}  maxVal={draft.max_k9?.toString() ?? ''}  onMinChange={(v) => set('min_k9',  optFloat(v))} onMaxChange={(v) => set('max_k9',  optFloat(v))} />
          </div>
        </PlaygroundFilterGroup>
      )}

      {/* ── Arc Score ── */}
      <PlaygroundFilterGroup label="Arc Score" defaultOpen={false}>
        <RangeRow
          label=""
          minVal={draft.min_value_score?.toString() ?? ''}
          maxVal={draft.max_value_score?.toString() ?? ''}
          onMinChange={(v) => set('min_value_score', optFloat(v))}
          onMaxChange={(v) => set('max_value_score', optFloat(v))}
          minPlaceholder="Min"
          maxPlaceholder="Max"
        />
      </PlaygroundFilterGroup>

      {/* ── Actions ── */}
      <div className="flex gap-2 pt-4">
        <button
          type="button"
          onClick={() => onSearch(draft)}
          className="shell-button shell-button-accent flex-1 px-4 py-2.5 text-[13px] font-medium"
        >
          Run Query
        </button>
        <button
          type="button"
          onClick={handleReset}
          className="shell-button px-4 py-2.5 text-[13px] font-medium"
        >
          Reset
        </button>
      </div>
    </div>
  )
}

// ── RangeRow ─────────────────────────────────────────────────────────────────
// Reusable min/max number input pair used throughout the filter panel.

interface RangeRowProps {
  label: string
  minVal: string
  maxVal: string
  onMinChange: (v: string) => void
  onMaxChange: (v: string) => void
  minPlaceholder?: string
  maxPlaceholder?: string
  singleValue?: boolean // renders only the min input (for season single-year filter)
}

function RangeRow({
  label,
  minVal,
  maxVal,
  onMinChange,
  onMaxChange,
  minPlaceholder = 'Min',
  maxPlaceholder = 'Max',
  singleValue = false,
}: RangeRowProps) {
  return (
    <div className="flex items-center gap-2">
      {label && (
        <span className="w-9 shrink-0 text-[11px] font-medium text-text-subtle">{label}</span>
      )}
      <input
        type="number"
        value={minVal}
        placeholder={minPlaceholder}
        onChange={(e: ChangeEvent<HTMLInputElement>) => onMinChange(e.target.value)}
        className="shell-input min-w-0 flex-1 text-[12px]"
      />
      {!singleValue && (
        <>
          <span className="text-[11px] text-text-subtle">–</span>
          <input
            type="number"
            value={maxVal}
            placeholder={maxPlaceholder}
            onChange={(e: ChangeEvent<HTMLInputElement>) => onMaxChange(e.target.value)}
            className="shell-input min-w-0 flex-1 text-[12px]"
          />
        </>
      )}
    </div>
  )
}
