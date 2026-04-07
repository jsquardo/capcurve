import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { usePlayerSearch } from '@/hooks/usePlayerSearch'
import type { PlayerListItem } from '@/types'

// MLB PrimaryPosition.Name values as returned by the ingestion pipeline.
// Unknown positions fall through to the full position string — safer than
// truncating an unmapped name.
const POSITION_ABBR: Record<string, string> = {
  'Pitcher': 'P',
  'Starting Pitcher': 'SP',
  'Relief Pitcher': 'RP',
  'Catcher': 'C',
  'First Base': '1B',
  'First Baseman': '1B',
  'Second Base': '2B',
  'Second Baseman': '2B',
  'Third Base': '3B',
  'Third Baseman': '3B',
  'Shortstop': 'SS',
  'Left Field': 'LF',
  'Left Fielder': 'LF',
  'Center Field': 'CF',
  'Center Fielder': 'CF',
  'Right Field': 'RF',
  'Right Fielder': 'RF',
  'Outfielder': 'OF',
  'Infielder': 'IF',
  'Designated Hitter': 'DH',
  'Two-Way Player': 'TWP',
}

function positionAbbr(position: string): string {
  return POSITION_ABBR[position] ?? position
}

type Props = {
  inputClassName?: string  // width override: "w-[240px]" desktop, "w-full" mobile
  onSelect?: () => void    // called after navigation — mobile uses this to close the panel
}

export default function NavSearch({ inputClassName = 'w-[240px]', onSelect }: Props) {
  const [inputValue, setInputValue] = useState('')
  const [activeIndex, setActiveIndex] = useState(-1)
  const containerRef = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()

  const { results, isLoading } = usePlayerSearch(inputValue)

  // isOpen is fully derived — no separate piece of state.
  // trimmedInput drives the threshold so whitespace-only input never opens
  // the dropdown. All close paths (Escape, click-outside, handleSelect) call
  // setInputValue(''), making trimmedInput.length < 2 and closing the panel.
  // The results.length guard is intentionally removed so the "No players
  // found" empty state can render when a valid query returns zero results.
  const trimmedInput = inputValue.trim()
  const isOpen = trimmedInput.length >= 2

  // Reset keyboard highlight when results change mid-type
  useEffect(() => {
    setActiveIndex(-1)
  }, [results])

  // Click-outside closes by clearing the input
  useEffect(() => {
    function handleMouseDown(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setInputValue('')
        setActiveIndex(-1)
      }
    }
    document.addEventListener('mousedown', handleMouseDown)
    return () => document.removeEventListener('mousedown', handleMouseDown)
  }, [])

  function handleSelect(player: PlayerListItem) {
    navigate(`/players/${player.id}`)
    setInputValue('')
    setActiveIndex(-1)
    onSelect?.()
  }

  function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setActiveIndex((i) => Math.min(i + 1, results.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setActiveIndex((i) => Math.max(i - 1, -1))
    } else if (e.key === 'Enter') {
      if (activeIndex >= 0 && results[activeIndex]) {
        handleSelect(results[activeIndex])
      }
    } else if (e.key === 'Escape') {
      setInputValue('')
      setActiveIndex(-1)
    }
  }

  return (
    <div ref={containerRef} className={`relative ${inputClassName}`}>
      <span className="pointer-events-none absolute inset-y-0 left-3 flex items-center text-[13px] text-text-subtle">
        ⌕
      </span>
      <input
        type="search"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Search any player..."
        aria-label="Search players"
        aria-expanded={isOpen}
        aria-haspopup="listbox"
        autoComplete="off"
        className="w-full rounded-[8px] border border-border bg-elevated py-2 pl-9 pr-4 text-[13px] text-text outline-none placeholder:text-text-subtle transition focus:border-accent"
      />

      {isOpen && (
        <div className="absolute left-0 top-full z-[60] mt-1 w-full min-w-[300px] rounded-[8px] border border-border bg-overlay py-1 shadow-lg">
          {isLoading ? (
            <div className="px-4 py-2.5 text-[13px] text-text-muted">Searching…</div>
          ) : results.length === 0 ? (
            <div className="px-4 py-2.5 text-[13px] text-text-muted">No players found</div>
          ) : (
            <ul role="listbox" aria-label="Player search results">
              {results.map((player, i) => (
                <SearchResultItem
                  key={player.id}
                  player={player}
                  isActive={i === activeIndex}
                  onMouseEnter={() => setActiveIndex(i)}
                  onSelect={() => handleSelect(player)}
                />
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}

// Extracted to keep the parent JSX under 60 lines
function SearchResultItem({
  player,
  isActive,
  onMouseEnter,
  onSelect,
}: {
  player: PlayerListItem
  isActive: boolean
  onMouseEnter: () => void
  onSelect: () => void
}) {
  const abbr = positionAbbr(player.position)
  const team = player.latest_season?.team_name ?? null

  return (
    <li
      role="option"
      aria-selected={isActive}
      onMouseEnter={onMouseEnter}
      onClick={onSelect}
      className={`flex cursor-pointer flex-col gap-0.5 px-4 py-2.5 transition-colors ${
        isActive ? 'bg-accent/10' : 'hover:bg-elevated'
      }`}
    >
      <div className="flex items-center gap-2">
        <span className="text-[13px] font-medium text-text">{player.full_name}</span>
        {!player.active && (
          <span className="rounded border border-border px-1 py-px text-[10px] leading-tight text-text-subtle">
            Retired
          </span>
        )}
      </div>
      <div className="text-[12px] text-text-muted">
        {abbr}{team ? ` · ${team}` : ''}
      </div>
    </li>
  )
}
