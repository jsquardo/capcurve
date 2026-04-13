import { useState } from 'react'
import type { ReactNode } from 'react'

interface Props {
  label: string
  children: ReactNode
  defaultOpen?: boolean
}

// Collapsible section wrapper used inside PlaygroundFilterPanel.
// All groups open by default — users collapse what they don't need.
export default function PlaygroundFilterGroup({ label, children, defaultOpen = true }: Props) {
  const [open, setOpen] = useState(defaultOpen)

  return (
    <div className="border-b border-border py-4 last:border-b-0 first:pt-0">
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen((prev) => !prev)}
        className="flex w-full items-center justify-between focus:outline-none focus-visible:ring-1 focus-visible:ring-accent rounded-[2px]"
      >
        <span className="text-[10px] font-semibold uppercase tracking-widest text-text-subtle">
          {label}
        </span>
        <span className={`text-[10px] text-text-subtle transition-transform duration-150 ${open ? 'rotate-180' : ''}`}>
          ▾
        </span>
      </button>

      {open && (
        <div className="mt-3">
          {children}
        </div>
      )}
    </div>
  )
}
