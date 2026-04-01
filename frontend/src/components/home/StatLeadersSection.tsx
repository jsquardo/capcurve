import { useState } from 'react'
import LeaderRow from './LeaderRow'

type LeaderEntry = { n: string; t: string; v: number | string }
type LeaderCategory = 'hr' | 'avg' | 'era' | 'k9' | 'ops'

// Mock data — replace with API responses when leaderboard endpoint is wired to the frontend
// TODO: wire to GET /api/v1/leaderboards when frontend data-fetching is added in Phase 3
const leadersData: Record<LeaderCategory, LeaderEntry[]> = {
  hr: [
    { n: 'Aaron Judge', t: 'NYY', v: 58 },
    { n: 'Shohei Ohtani', t: 'LAD', v: 54 },
    { n: 'Kyle Schwarber', t: 'PHI', v: 47 },
    { n: 'Pete Alonso', t: 'NYM', v: 46 },
    { n: 'Yordan Alvarez', t: 'HOU', v: 45 },
  ],
  avg: [
    { n: 'Luis Arraez', t: 'MIA', v: '.354' },
    { n: 'Freddie Freeman', t: 'LAD', v: '.331' },
    { n: 'Paul Goldschmidt', t: 'STL', v: '.317' },
    { n: 'Rafael Devers', t: 'BOS', v: '.312' },
    { n: 'Trea Turner', t: 'PHI', v: '.308' },
  ],
  era: [
    { n: 'Spencer Strider', t: 'ATL', v: '2.11' },
    { n: 'Zack Wheeler', t: 'PHI', v: '2.34' },
    { n: 'Gerrit Cole', t: 'NYY', v: '2.56' },
    { n: 'Sandy Alcantara', t: 'MIA', v: '2.70' },
    { n: 'Kevin Gausman', t: 'SF', v: '2.81' },
  ],
  k9: [
    { n: 'Spencer Strider', t: 'ATL', v: '13.7' },
    { n: 'Corbin Burnes', t: 'BAL', v: '11.9' },
    { n: 'Dylan Cease', t: 'SD', v: '11.4' },
    { n: 'Julio Urías', t: 'LAD', v: '10.8' },
    { n: 'Shane Bieber', t: 'CLE', v: '10.3' },
  ],
  ops: [
    { n: 'Shohei Ohtani', t: 'LAD', v: '1.038' },
    { n: 'Aaron Judge', t: 'NYY', v: '.999' },
    { n: 'Yordan Alvarez', t: 'HOU', v: '.987' },
    { n: 'Ronald Acuña Jr.', t: 'ATL', v: '.961' },
    { n: 'Luis Robert Jr.', t: 'CWS', v: '.943' },
  ],
}

const leaderCategories: { key: LeaderCategory; label: string }[] = [
  { key: 'hr', label: 'HR' },
  { key: 'avg', label: 'AVG' },
  { key: 'era', label: 'ERA' },
  { key: 'k9', label: 'K/9' },
  { key: 'ops', label: 'OPS' },
]

// Bar width for numeric stats: proportional to leader. For string stats (AVG, ERA):
// rank-based so the top entry is widest without needing to parse formatted strings.
function computeBarPct(index: number, value: number | string, leaders: LeaderEntry[]): number {
  if (typeof leaders[0].v === 'number') {
    return Math.round((Number(value) / Number(leaders[0].v)) * 100)
  }
  return Math.round(((leaders.length - index) / leaders.length) * 85) + 15
}

export default function StatLeadersSection() {
  const [leaderCat, setLeaderCat] = useState<LeaderCategory>('hr')
  const leaders = leadersData[leaderCat]

  return (
    <div>
      <div className="mb-7 flex items-end justify-between">
        <div>
          <div className="mb-[6px] text-[10px] font-semibold uppercase tracking-[2px] text-accent">
            2024 Season
          </div>
          <div className="font-display text-[32px] leading-none tracking-[1px]">Stat Leaders</div>
        </div>
        <a href="/leaderboards" className="border-b border-link/30 pb-[2px] text-[13px] text-link">
          Full leaderboards →
        </a>
      </div>

      <div className="mb-5 flex flex-wrap gap-2">
        {leaderCategories.map(({ key, label }) => (
          <button
            key={key}
            type="button"
            onClick={() => setLeaderCat(key)}
            className={`rounded-full border px-[14px] py-[5px] text-[12px] font-medium transition-all ${
              leaderCat === key
                ? 'border-accent bg-accent text-[#0a0d12]'
                : 'border-border text-text-subtle hover:border-border-strong hover:text-text-muted'
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      <div className="flex flex-col gap-[6px]">
        {leaders.map((entry, i) => (
          <LeaderRow
            key={entry.n}
            rank={i + 1}
            name={entry.n}
            team={entry.t}
            value={entry.v}
            barPct={computeBarPct(i, entry.v, leaders)}
          />
        ))}
      </div>
    </div>
  )
}
