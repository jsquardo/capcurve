import HeroSection from '../components/home/HeroSection'
import TrendingSection from '../components/home/TrendingSection'
import StatLeadersSection from '../components/home/StatLeadersSection'

// --- Feed ---
type FeedItem = { type: 'insight' | 'news'; text: string; meta: string }

const feedItems: FeedItem[] = [
  { type: 'insight', text: "Paul Skenes is tracking as one of the fastest arc accelerations for a debut SP since Strasburg in 2010", meta: 'CapCurve · 6h' },
  { type: 'news', text: 'Freddie Freeman agrees to 2-year extension through the 2027 season', meta: 'MLB.com · 4h' },
  { type: 'insight', text: "Bobby Witt Jr.'s arc trajectory now mirrors a young Nolan Arenado at the same career stage", meta: 'CapCurve · 11h' },
  { type: 'news', text: 'AL MVP race heating up: Judge and Witt lead advanced metrics in second half', meta: 'Baseball America · 14h' },
]

export default function HomePage() {
  return (
    <>
      <HeroSection />

      <TrendingSection />

      {/* STAT LEADERS + FEED */}
      <section className="border-b border-border px-4 py-12 sm:px-6 lg:px-10">
        <div className="grid grid-cols-1 gap-12 lg:grid-cols-[1.1fr_1fr]">

          <StatLeadersSection />

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
