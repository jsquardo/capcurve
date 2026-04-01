import { useEffect, useState } from 'react'
import { NavLink } from 'react-router-dom'

type Theme = 'dark' | 'light'

const NAV_ITEMS = [
  { label: 'Home', to: '/' },
  { label: 'Leaderboards', to: '/leaderboards' },
] as const

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
      <nav className="shell-container flex min-h-[72px] items-center gap-4 py-4">
        <div className="flex min-w-0 flex-1 items-center gap-4 lg:gap-8">
          <NavLink
            to="/"
            onClick={closeMobileMenu}
            className="shrink-0 font-display text-[2rem] leading-none tracking-[0.14em] text-accent"
          >
            CAP<span className="text-text">CURVE</span>
          </NavLink>

          <div className="hidden items-center gap-6 lg:flex">
            {NAV_ITEMS.map((item) => (
              <NavLink
                key={item.label}
                to={item.to}
                className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
              >
                {item.label}
              </NavLink>
            ))}
            {/* TODO: wire to /players when PlayerListPage exists */}
            <span className="text-sm font-medium tracking-[0.02em] text-text-muted/75">Players</span>
            {/* TODO: wire to /playground when PlaygroundPage is built */}
            <span className="text-sm font-medium tracking-[0.02em] text-text-muted/75">Stat Playground</span>
          </div>
        </div>

        <div className="hidden min-w-0 flex-1 items-center justify-end gap-3 md:flex">
          <div className="relative w-full max-w-xs">
            <input
              type="search"
              placeholder="Search players..."
              className="shell-input w-full pr-10"
              aria-label="Search players"
            />
            <span className="pointer-events-none absolute inset-y-0 right-4 flex items-center text-text-subtle">
              ⌕
            </span>
          </div>
          <button
            type="button"
            onClick={toggleTheme}
            className="shell-button h-11 w-11"
            aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
            title={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
          >
            {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
          </button>
        </div>

        <div className="ml-auto flex items-center gap-2 md:hidden">
          <button
            type="button"
            onClick={toggleTheme}
            className="shell-button h-10 w-10"
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
        <div id="mobile-nav-panel" className="border-t border-border/70 md:hidden">
          <div className="shell-container space-y-4 py-4">
            <div className="space-y-3">
              {NAV_ITEMS.map((item) => (
                <NavLink
                  key={item.label}
                  to={item.to}
                  onClick={closeMobileMenu}
                  className={({ isActive }) =>
                    `block text-base font-medium ${isActive ? 'text-text' : 'text-text-muted'}`
                  }
                >
                  {item.label}
                </NavLink>
              ))}
              {/* TODO: wire to /players when PlayerListPage exists */}
              <div className="text-base font-medium text-text-muted/75">Players</div>
              {/* TODO: wire to /playground when PlaygroundPage is built */}
              <div className="text-base font-medium text-text-muted/75">Stat Playground</div>
            </div>
            <input
              type="search"
              placeholder="Search players..."
              className="shell-input w-full"
              aria-label="Search players"
            />
          </div>
        </div>
      ) : null}
    </header>
  )
}
