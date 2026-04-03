import { useEffect, useState } from 'react'
import { Link, NavLink } from 'react-router-dom'

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
          {NAV_ITEMS.map((item) => (
            <li key={item.label}>
              <NavLink
                to={item.to}
                className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
              >
                {item.label}
              </NavLink>
            </li>
          ))}
          <li>
            <NavLink
              to="/players"
              className={({ isActive }) => `nav-link ${isActive ? 'nav-link-active' : ''}`}
            >
              Players
            </NavLink>
          </li>
          {/* TODO: wire to /playground when PlaygroundPage is built */}
          <li>
            <span className="nav-link opacity-60 cursor-default">Playground</span>
          </li>
        </ul>

        {/* Desktop right: search + Explore + theme toggle */}
        <div className="hidden items-center gap-3 lg:flex ml-auto">
          <div className="relative">
            <span className="pointer-events-none absolute inset-y-0 left-3 flex items-center text-[13px] text-text-subtle">
              ⌕
            </span>
            <input
              type="search"
              placeholder="Search any player..."
              aria-label="Search players"
              className="w-[240px] rounded-[8px] border border-border bg-elevated py-2 pl-9 pr-4 text-[13px] text-text outline-none placeholder:text-text-subtle transition focus:border-accent"
            />
          </div>
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
              {NAV_ITEMS.map((item) => (
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
              ))}
              <NavLink
                to="/players"
                onClick={closeMobileMenu}
                className={({ isActive }) =>
                  `block text-[13px] font-medium ${isActive ? 'text-text' : 'text-text-muted'}`
                }
              >
                Players
              </NavLink>
              {/* TODO: wire to /playground when PlaygroundPage is built */}
              <div className="text-[13px] font-medium text-text-muted opacity-60">Playground</div>
            </div>
            <input
              type="search"
              placeholder="Search any player..."
              className="w-full rounded-[8px] border border-border bg-elevated py-2 px-4 text-[13px] text-text outline-none placeholder:text-text-subtle transition focus:border-accent"
              aria-label="Search players"
            />
          </div>
        </div>
      ) : null}
    </header>
  )
}
