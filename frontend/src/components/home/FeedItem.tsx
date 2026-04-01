interface FeedItemProps {
  type: 'insight' | 'news'
  text: string
  meta: string
}

export default function FeedItem({ type, text, meta }: FeedItemProps) {
  return (
    <div className="flex cursor-pointer items-start gap-[10px] rounded-[8px] border border-border bg-elevated px-3 py-[10px] transition-colors hover:border-border-strong">
      <span
        className={`mt-[2px] shrink-0 rounded-[4px] px-[7px] py-[3px] text-[9px] font-bold uppercase tracking-[1px] ${
          type === 'insight'
            ? 'border border-accent/20 bg-accent/[0.12] text-accent'
            : 'border border-link/20 bg-link/[0.10] text-link'
        }`}
      >
        {type}
      </span>
      <div className="flex-1">
        <div className="text-[13px] font-medium leading-[1.4] text-text">{text}</div>
        <div className="mt-[3px] text-[10px] text-text-subtle">{meta}</div>
      </div>
    </div>
  )
}
