import { useEffect, useRef, useState } from 'react'
import { Link, NavLink } from 'react-router-dom'

export type DropdownItem = {
  label: string
  to?: string // absent → disabled/coming-soon item
}

type Props = {
  label: string
  to: string
  items: DropdownItem[]
}

export default function NavDropdown({ label, to, items }: Props) {
  const [open, setOpen] = useState(false)
  const closeTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  function scheduleClose() {
    closeTimer.current = setTimeout(() => setOpen(false), 120)
  }

  function cancelClose() {
    if (closeTimer.current !== null) {
      clearTimeout(closeTimer.current)
      closeTimer.current = null
    }
  }

  // Close on Escape while panel is open
  useEffect(() => {
    if (!open) return
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open])

  return (
    <div
      className="relative"
      onMouseEnter={() => { cancelClose(); setOpen(true) }}
      onMouseLeave={scheduleClose}
    >
      {/* Trigger — NavLink navigates on click; chevron rotates when panel is open.
          The entire outer div (NavLink + chevron) is the hover zone. */}
      <div className="flex items-center gap-1">
        <NavLink
          to={to}
          onClick={() => setOpen(false)}
          className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
        >
          {label}
        </NavLink>
        <span
          aria-hidden="true"
          className={`text-[10px] text-text-subtle transition-transform duration-150 ${open ? 'rotate-180' : ''}`}
        >
          ▾
        </span>
      </div>

      {/* Dropdown panel — z-10 is sufficient because it inherits the navbar's z-50
          stacking context and therefore paints above all page content. */}
      {open && (
        <div
          role="menu"
          aria-label={`${label} submenu`}
          className="absolute left-0 top-full z-[60] mt-2 min-w-[180px] rounded-[8px] border border-border bg-panel py-1.5 shadow-lg"
        >
          {items.map((item) =>
            item.to ? (
              <Link
                key={item.label}
                to={item.to}
                role="menuitem"
                onClick={() => setOpen(false)}
                className="block rounded-[6px] px-4 py-2 text-[13px] font-medium text-text-muted transition-colors hover:bg-elevated hover:text-text"
              >
                {item.label}
              </Link>
            ) : (
              // TODO: wire these items when their destination pages/routes are built
              <span
                key={item.label}
                role="menuitem"
                aria-disabled="true"
                className="block cursor-default select-none px-4 py-2 text-[13px] font-medium text-text-subtle opacity-50"
              >
                {item.label}
              </span>
            )
          )}
        </div>
      )}
    </div>
  )
}
