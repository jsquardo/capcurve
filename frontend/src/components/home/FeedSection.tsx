import FeedItem from './FeedItem'

// Mock data — replace with API responses when RSS + insights endpoints are ready
// TODO: wire to RSS aggregation + auto-generated insights endpoints (Phase 3 backend feature)
const feedItems = [
  { type: 'insight' as const, text: "Paul Skenes is tracking as one of the fastest arc accelerations for a debut SP since Strasburg in 2010", meta: 'CapCurve · 6h' },
  { type: 'news' as const, text: 'Freddie Freeman agrees to 2-year extension through the 2027 season', meta: 'MLB.com · 4h' },
  { type: 'insight' as const, text: "Bobby Witt Jr.'s arc trajectory now mirrors a young Nolan Arenado at the same career stage", meta: 'CapCurve · 11h' },
  { type: 'news' as const, text: 'AL MVP race heating up: Judge and Witt lead advanced metrics in second half', meta: 'Baseball America · 14h' },
]

export default function FeedSection() {
  return (
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

      {/* Featured insight — hardcoded until auto-generated insights endpoint exists */}
      <div className="mb-[10px] cursor-pointer rounded-xl border border-border bg-elevated p-6 transition-colors hover:border-link">
        <div className="mb-3 inline-flex items-center gap-[5px] rounded-[5px] border border-accent/20 bg-accent/[0.12] px-[9px] py-[3px] text-[10px] font-semibold uppercase tracking-[1px] text-accent">
          CapCurve Insight
        </div>
        <div className="mb-2 text-[16px] font-semibold leading-[1.5]">
          Juan Soto&apos;s arc score climbed 11.2 points in 14 days — steepest monthly rise since his 2022 breakout
        </div>
        <div className="text-[11px] text-text-subtle">2 hours ago · CapCurve Analysis</div>
      </div>

      <div className="flex flex-col gap-[6px]">
        {feedItems.map((item, i) => (
          <FeedItem key={i} type={item.type} text={item.text} meta={item.meta} />
        ))}
      </div>
    </div>
  )
}
