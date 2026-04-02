interface LeaderboardPaginationProps {
  page: number
  totalPages: number
  onPrev: () => void
  onNext: () => void
}

export default function LeaderboardPagination({ page, totalPages, onPrev, onNext }: LeaderboardPaginationProps) {
  if (totalPages <= 1) return null

  return (
    <div className="flex items-center justify-center gap-4">
      <button
        onClick={onPrev}
        disabled={page <= 1}
        className="rounded-[6px] border border-border bg-elevated px-4 py-1.5 text-[12px] font-medium text-text-muted transition-colors hover:border-border-strong hover:text-text disabled:cursor-not-allowed disabled:opacity-40"
      >
        ← Prev
      </button>
      <span className="font-mono text-[12px] text-text-subtle">
        Page {page} of {totalPages}
      </span>
      <button
        onClick={onNext}
        disabled={page >= totalPages}
        className="rounded-[6px] border border-border bg-elevated px-4 py-1.5 text-[12px] font-medium text-text-muted transition-colors hover:border-border-strong hover:text-text disabled:cursor-not-allowed disabled:opacity-40"
      >
        Next →
      </button>
    </div>
  )
}
