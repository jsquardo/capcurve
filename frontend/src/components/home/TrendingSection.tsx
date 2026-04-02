import { useState } from 'react'
import TrendingCard from './TrendingCard'

// Mock data — replace with API responses when arc-delta and page-view endpoints are ready
// TODO: wire to /api/v1/players/trending?signal=arc_delta&days=14 (Phase 3 backend feature)
const hotTrendData = [
  { rank: 1, name: 'Elly De La Cruz', team: 'CIN · SS', delta: '+9.1 arc pts', bars: [30, 38, 42, 55, 60, 62, 58, 70, 82, 91] },
  { rank: 2, name: 'Paul Skenes', team: 'PIT · SP', delta: '+7.4 arc pts', bars: [40, 48, 52, 60, 65, 70, 68, 75, 85, 88] },
  { rank: 3, name: 'Jackson Chourio', team: 'MIL · OF', delta: '+6.8 arc pts', bars: [20, 25, 30, 38, 44, 50, 54, 60, 70, 78] },
  { rank: 4, name: 'Bobby Witt Jr.', team: 'KC · SS', delta: '+5.9 arc pts', bars: [25, 32, 40, 50, 58, 62, 66, 70, 76, 84] },
]

// TODO: wire to /api/v1/players/trending?signal=page_views&days=14 (Phase 3 backend feature)
const viewTrendData = [
  { rank: 1, name: 'Shohei Ohtani', team: 'LAD · DH/SP', delta: '84.2k views', bars: [60, 65, 70, 72, 75, 80, 82, 85, 88, 90], isViews: true },
  { rank: 2, name: 'Aaron Judge', team: 'NYY · RF', delta: '71.5k views', bars: [55, 62, 68, 74, 80, 84, 86, 88, 85, 82], isViews: true },
  { rank: 3, name: 'Mike Trout', team: 'LAA · CF', delta: '58.3k views', bars: [48, 76, 80, 82, 88, 87, 91, 91, 90, 72], isViews: true },
  { rank: 4, name: 'Fernando Tatis Jr.', team: 'SD · RF', delta: '44.1k views', bars: [35, 55, 72, 78, 65, 60, 68, 74, 78, 80], isViews: true },
]

export default function TrendingSection() {
  const [trendTab, setTrendTab] = useState<'hot' | 'view'>('hot')
  const trendData = trendTab === 'hot' ? hotTrendData : viewTrendData

  return (
    <section className="border-b border-border">
      <div className="shell-container py-12">
      <div className="mb-7 flex items-end justify-between">
        <div>
          <div className="mb-[6px] text-[10px] font-semibold uppercase tracking-[2px] text-accent">
            Live Data · Past 14 Days
          </div>
          <div className="font-display text-[32px] leading-none tracking-[1px]">Trending Players</div>
        </div>
        {/* TODO: wire to /players when PlayerListPage exists */}
        <span className="border-b border-link/30 pb-[2px] text-[13px] text-link">
          View all players →
        </span>
      </div>

      <div className="mb-6 flex border-b border-border">
        <button
          type="button"
          onClick={() => setTrendTab('hot')}
          className={`-mb-px border-b-2 px-[22px] py-[10px] text-[13px] font-medium transition-colors ${
            trendTab === 'hot'
              ? 'border-accent text-accent'
              : 'border-transparent text-text-subtle hover:text-text-muted'
          }`}
        >
          🔥 Hottest Arc Movement
        </button>
        <button
          type="button"
          onClick={() => setTrendTab('view')}
          className={`-mb-px border-b-2 px-[22px] py-[10px] text-[13px] font-medium transition-colors ${
            trendTab === 'view'
              ? 'border-accent text-accent'
              : 'border-transparent text-text-subtle hover:text-text-muted'
          }`}
        >
          👁 Most Viewed
        </button>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {trendData.map((player) => (
          <TrendingCard
            key={player.name}
            rank={player.rank}
            name={player.name}
            team={player.team}
            delta={player.delta}
            bars={player.bars}
            isViews={player.isViews}
          />
        ))}
      </div>
      </div>
    </section>
  )
}
