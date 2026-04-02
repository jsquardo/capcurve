import type { LeaderboardCategory } from '@/types'

interface CategoryConfig {
  id: LeaderboardCategory
  label: string
  description?: string
}

const CATEGORIES: CategoryConfig[] = [
  {
    id: 'peak_arc',
    label: 'Peak Arc',
    description: 'CapCurve composite career peak score',
  },
  { id: 'hr',  label: 'HR'  },
  { id: 'avg', label: 'AVG' },
  { id: 'era', label: 'ERA' },
  { id: 'k9',  label: 'K/9' },
]

interface LeaderboardCategoryTabsProps {
  activeCategory: LeaderboardCategory
  onSelect: (category: LeaderboardCategory) => void
}

export default function LeaderboardCategoryTabs({ activeCategory, onSelect }: LeaderboardCategoryTabsProps) {
  const active = CATEGORIES.find(c => c.id === activeCategory)

  return (
    <div>
      <div className="flex flex-wrap gap-2">
        {CATEGORIES.map(cat => (
          <button
            key={cat.id}
            onClick={() => onSelect(cat.id)}
            className={
              cat.id === activeCategory
                ? 'rounded-full border border-accent bg-accent/10 px-4 py-1.5 text-[12px] font-medium text-accent transition-colors'
                : 'rounded-full border border-border bg-elevated px-4 py-1.5 text-[12px] font-medium text-text-muted transition-colors hover:border-border-strong hover:text-text'
            }
          >
            {cat.label}
          </button>
        ))}
      </div>
      {active?.description && (
        <p className="mt-2 text-[11px] text-text-subtle">{active.description}</p>
      )}
    </div>
  )
}
