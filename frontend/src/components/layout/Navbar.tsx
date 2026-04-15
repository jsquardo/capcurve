import { useEffect, useState } from 'react'
import { Link, NavLink } from 'react-router-dom'
import NavDropdown, { type DropdownItem } from './NavDropdown'
import NavSearch from './NavSearch'

type Theme = 'dark' | 'light'

type NavItem =
  | { label: string; to: string; items?: never }
  | { label: string; to: string; items: DropdownItem[] }

const NAV_ITEMS: NavItem[] = [
  { label: 'Home', to: '/' },
  {
    label: 'Players',
    to: '/players',
    items: [
      { label: 'Top Players', to: '/players' },
      { label: 'Trending Players' },       // disabled — needs arc delta backend (Phase 3)
      { label: 'Hitters' },                // disabled — needs URL position filtering on PlayersPage
      { label: 'Pitchers' },               // disabled — same
      { label: 'Top Stats', to: '/leaderboards' },
    ],
  },
  {
    label: 'Leaderboards',
    to: '/leaderboards',
    items: [
      { label: 'Career Arc Peaks', to: '/leaderboards' },
      { label: 'Stat Leaders', to: '/leaderboards' },
      { label: 'Playground Leaderboards' }, // disabled — requires playground
    ],
  },
]

function SunIcon() {
  return (
    <svg aria-hidden="true" viewBox="0 0 24 24" className="h-4 w-4 fill-none stroke-current stroke-[1.8]">
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2.75v2.5M12 18.75v2.5M21.25 12h-2.5M5.25 12h-2.5M18.54 5.46l-1.77 1.77M7.23 16.77l-1.77 1.77M18.54 18.54l-1.77-1.77M7.23 7.23 5.46 5.46" />
    </svg>
  )
}

function MoonIcon() {
  return (
    <svg aria-hidden="true" viewBox="0 0 24 24" className="h-4 w-4 fill-none stroke-current stroke-[1.8]">
      <path d="M20.5 14.16A8.5 8.5 0 1 1 9.84 3.5a6.75 6.75 0 1 0 10.66 10.66Z" />
    </svg>
  )
}

export default function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window === 'undefined') return 'dark'
    const stored = window.localStorage.getItem('capcurve-theme')
    return stored === 'light' ? 'light' : 'dark'
  })

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    window.localStorage.setItem('capcurve-theme', theme)
  }, [theme])

  function toggleTheme() {
    setTheme((current) => (current === 'dark' ? 'light' : 'dark'))
  }

  function closeMobileMenu() {
    setMobileOpen(false)
  }

  return (
    <header className="shell-panel sticky top-0 z-50">
      <nav className="shell-container flex h-[60px] items-center gap-4">
        {/* Logo */}
        <NavLink
          to="/"
          onClick={closeMobileMenu}
          className="shrink-0 font-display text-[2rem] leading-none tracking-[3px] text-accent"
        >
          CAP<span className="text-text">CURVE</span>
        </NavLink>

        {/* Desktop nav links */}
        <ul className="hidden items-center gap-10 list-none lg:flex ml-8">
          {NAV_ITEMS.map((item) =>
            item.items ? (
              <li key={item.label}>
                <NavDropdown label={item.label} to={item.to} items={item.items} />
              </li>
            ) : (
              <li key={item.label}>
                <NavLink
                  to={item.to}
                  className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
                >
                  {item.label}
                </NavLink>
              </li>
            )
          )}
          <li>
            <NavLink
              to="/playground"
              className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
            >
              Playground
            </NavLink>
          </li>
        </ul>

        {/* Desktop right: search + Explore + theme toggle */}
        <div className="hidden items-center gap-3 lg:flex ml-auto">
          <NavSearch inputClassName="w-[240px]" />
          <Link
            to="/players"
            className="rounded-[7px] bg-accent px-[18px] py-[7px] text-[13px] font-medium text-[#0a0d12] transition hover:bg-accent-strong"
          >
            Explore
          </Link>
          <button
            type="button"
            onClick={toggleTheme}
            className="p-1.5 text-text-subtle transition-colors hover:text-text"
            aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
            title={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
          >
            {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
          </button>
        </div>

        {/* Mobile right: theme toggle + hamburger */}
        <div className="ml-auto flex items-center gap-2 lg:hidden">
          <button
            type="button"
            onClick={toggleTheme}
            className="p-1.5 text-text-subtle hover:text-text"
            aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
          >
            {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
          </button>
          <button
            type="button"
            onClick={() => setMobileOpen((open) => !open)}
            className="shell-button h-10 w-10"
            aria-expanded={mobileOpen}
            aria-controls="mobile-nav-panel"
            aria-label="Toggle navigation menu"
          >
            <span className="text-lg leading-none">{mobileOpen ? '×' : '≡'}</span>
          </button>
        </div>
      </nav>

      {mobileOpen ? (
        <div id="mobile-nav-panel" className="border-t border-border/70 lg:hidden">
          <div className="shell-container space-y-4 py-4">
            <div className="space-y-3">
              {NAV_ITEMS.map((item) =>
                item.items ? (
                  // Dropdown parent: show non-link label then indented sub-items
                  <div key={item.label}>
                    <div className="text-[11px] font-semibold uppercase tracking-widest text-text-subtle pb-1">
                      {item.label}
                    </div>
                    <div className="space-y-2 pl-3">
                      {item.items.map((sub) =>
                        sub.to ? (
                          <NavLink
                            key={sub.label}
                            to={sub.to}
                            onClick={closeMobileMenu}
                            className={({ isActive }) =>
                              `block text-[13px] font-medium ${isActive ? 'text-text' : 'text-text-muted'}`
                            }
                          >
                            {sub.label}
                          </NavLink>
                        ) : (
                          <span
                            key={sub.label}
                            className="block text-[13px] font-medium text-text-subtle opacity-50 cursor-default select-none"
                          >
                            {sub.label}
                          </span>
                        )
                      )}
                    </div>
                  </div>
                ) : (
                  <NavLink
                    key={item.label}
                    to={item.to}
                    onClick={closeMobileMenu}
                    className={({ isActive }) =>
                      `block text-[13px] font-medium ${isActive ? 'text-text' : 'text-text-muted'}`
                    }
                  >
                    {item.label}
                  </NavLink>
                )
              )}
              <NavLink
                to="/playground"
                onClick={closeMobileMenu}
                className={({ isActive }) =>
                  `block text-[13px] font-medium ${isActive ? 'text-text' : 'text-text-muted'}`
                }
              >
                Playground
              </NavLink>
            </div>
            <NavSearch inputClassName="w-full" onSelect={closeMobileMenu} />
          </div>
        </div>
      ) : null}
    </header>
  )
}
