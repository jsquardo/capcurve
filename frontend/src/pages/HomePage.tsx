import { useState } from 'react'
import HeroSection from '../components/home/HeroSection'
import TrendingSection from '../components/home/TrendingSection'

// --- Stat Leaders ---
type LeaderEntry = { n: string; t: string; v: number | string }
type LeaderCategory = 'hr' | 'avg' | 'era' | 'k9' | 'ops'

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

// --- Feed ---
type FeedItem = { type: 'insight' | 'news'; text: string; meta: string }

const feedItems: FeedItem[] = [
  { type: 'insight', text: "Paul Skenes is tracking as one of the fastest arc accelerations for a debut SP since Strasburg in 2010", meta: 'CapCurve · 6h' },
  { type: 'news', text: 'Freddie Freeman agrees to 2-year extension through the 2027 season', meta: 'MLB.com · 4h' },
  { type: 'insight', text: "Bobby Witt Jr.'s arc trajectory now mirrors a young Nolan Arenado at the same career stage", meta: 'CapCurve · 11h' },
  { type: 'news', text: 'AL MVP race heating up: Judge and Witt lead advanced metrics in second half', meta: 'Baseball America · 14h' },
]

export default function HomePage() {
  const [leaderCat, setLeaderCat] = useState<LeaderCategory>('hr')

  const leaders = leadersData[leaderCat]
  const isNumericLeader = typeof leaders[0].v === 'number'

  function leaderBarPct(i: number, v: number | string): number {
    if (isNumericLeader) {
      return Math.round((Number(v) / Number(leaders[0].v)) * 100)
    }
    // Rank-based width for string stat values (AVG, ERA, etc.)
    return Math.round(((leaders.length - i) / leaders.length) * 85) + 15
  }

  return (
    <>
      <HeroSection />

      <TrendingSection />

      {/* STAT LEADERS + FEED */}
      <section className="border-b border-border px-4 py-12 sm:px-6 lg:px-10">
        <div className="grid gap-12" style={{ gridTemplateColumns: '1.1fr 1fr' }}>

          {/* STAT LEADERS */}
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

            {/* Category pills */}
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

            {/* Leaderboard rows */}
            <div className="flex flex-col gap-[6px]">
              {leaders.map((entry, i) => (
                <div
                  key={entry.n}
                  className="flex cursor-pointer items-center gap-3 rounded-[8px] border border-border bg-elevated px-[14px] py-[10px] transition-colors hover:border-border-strong"
                >
                  <span className="w-[18px] shrink-0 text-right font-mono text-[11px] text-text-subtle">
                    {i + 1}
                  </span>
                  <div className="flex-1">
                    <div className="text-[13px] font-medium">{entry.n}</div>
                    <div className="text-[10px] text-text-subtle">{entry.t}</div>
                  </div>
                  <div className="h-1 w-24 shrink-0 overflow-hidden rounded-sm bg-panel">
                    <div
                      className="h-full rounded-sm bg-accent transition-all duration-300"
                      style={{ width: `${leaderBarPct(i, entry.v)}%` }}
                    />
                  </div>
                  <span className="w-[52px] shrink-0 text-right font-mono text-[14px] font-medium text-accent">
                    {entry.v}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* FEED */}
          <div>
            <div className="mb-7 flex items-end justify-between">
              <div>
                <div className="mb-[6px] text-[10px] font-semibold uppercase tracking-[2px] text-accent">
                  Insights &amp; News
                </div>
                <div className="font-display text-[32px] leading-none tracking-[1px]">Latest</div>
              </div>
              {/* TODO: wire to feed/insights page when it exists */}
              <span className="border-b border-link/30 pb-[2px] text-[13px] text-link">
                More →
              </span>
            </div>

            {/* Featured insight */}
            <div className="mb-[10px] cursor-pointer rounded-xl border border-border bg-elevated p-6 transition-colors hover:border-link">
              <div className="mb-3 inline-flex items-center gap-[5px] rounded-[5px] border border-accent/20 bg-accent/[0.12] px-[9px] py-[3px] text-[10px] font-semibold uppercase tracking-[1px] text-accent">
                CapCurve Insight
              </div>
              <div className="mb-2 text-[16px] font-semibold leading-[1.5]">
                Juan Soto&apos;s arc score climbed 11.2 points in 14 days — steepest monthly rise since his 2022 breakout
              </div>
              <div className="text-[11px] text-text-subtle">2 hours ago · CapCurve Analysis</div>
            </div>

            {/* Feed items */}
            <div className="flex flex-col gap-[6px]">
              {feedItems.map((item, i) => (
                <div
                  key={i}
                  className="flex cursor-pointer items-start gap-[10px] rounded-[8px] border border-border bg-elevated px-3 py-[10px] transition-colors hover:border-border-strong"
                >
                  <span
                    className={`mt-[2px] shrink-0 rounded-[4px] px-[7px] py-[3px] text-[9px] font-bold uppercase tracking-[1px] ${
                      item.type === 'insight'
                        ? 'border border-accent/20 bg-accent/[0.12] text-accent'
                        : 'border border-link/20 bg-link/[0.10] text-link'
                    }`}
                  >
                    {item.type}
                  </span>
                  <div className="flex-1">
                    <div className="text-[13px] font-medium leading-[1.4] text-text">{item.text}</div>
                    <div className="mt-[3px] text-[10px] text-text-subtle">{item.meta}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

        </div>
      </section>
    </>
  )
}
