import { useEffect, useState } from 'react'
import { Outlet, NavLink } from 'react-router-dom'
import { motion } from 'framer-motion'

const tickerItems = [
  { name: 'Elly De La Cruz', val: 'Arc +9.1', dir: '▲' },
  { name: 'Paul Skenes', val: 'Arc +7.4', dir: '▲' },
  { name: 'Aaron Judge', val: '58 HR', dir: '' },
  { name: 'Bobby Witt Jr.', val: 'Arc +5.9', dir: '▲' },
  { name: 'Shohei Ohtani', val: 'OPS 1.038', dir: '' },
  { name: 'Jackson Chourio', val: 'Arc +6.8', dir: '▲' },
  { name: 'Spencer Strider', val: 'ERA 2.11', dir: '' },
  { name: 'Gunnar Henderson', val: 'Arc +4.2', dir: '▲' },
]

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

export default function Layout() {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window === 'undefined') {
      return 'dark'
    }

    const storedTheme = window.localStorage.getItem('capcurve-theme')
    return storedTheme === 'light' ? 'light' : 'dark'
  })

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    window.localStorage.setItem('capcurve-theme', theme)
  }, [theme])

  function toggleTheme() {
    setTheme((currentTheme) => (currentTheme === 'dark' ? 'light' : 'dark'))
  }

  function closeMobileMenu() {
    setMobileOpen(false)
  }

  return (
    <div className="min-h-screen bg-transparent text-text">
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
              <span className="text-sm font-medium tracking-[0.02em] text-text-muted/75">Players</span>
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
                    className={({ isActive }) => `block text-base font-medium ${isActive ? 'text-text' : 'text-text-muted'}`}
                  >
                    {item.label}
                  </NavLink>
                ))}
                <div className="text-base font-medium text-text-muted/75">Players</div>
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

      {/* Ticker lives here — outside shell-container so it's naturally full-width */}
      <div className="overflow-hidden bg-accent py-2">
        <motion.div
          initial={{ x: '0%' }}
          animate={{ x: '-50%' }}
          transition={{ duration: 32, repeat: Infinity, ease: 'linear' }}
          className="flex w-max whitespace-nowrap"
        >
          {[...tickerItems, ...tickerItems].map((item, index) => (
            <div
              key={index}
              className="inline-flex items-center gap-[10px] border-r border-black/15 px-8 text-[12px] font-semibold tracking-[0.3px] text-[#0a0d12]"
            >
              <span>{item.name}</span>
              <span className="font-mono text-[12px]">{item.val}</span>
              {item.dir ? <span className="text-[10px] font-bold opacity-70">{item.dir}</span> : null}
            </div>
          ))}
        </motion.div>
      </div>

      <main className="flex-1">
        <div className="shell-container py-8 sm:py-10 lg:py-12">
          <Outlet />
        </div>
      </main>

      <footer className="border-t border-border/70 bg-app/80 backdrop-blur-sm">
        <div className="shell-container flex flex-col gap-4 py-6 text-sm text-text-muted sm:flex-row sm:items-center sm:justify-between">
          <div>
            <div className="font-display text-xl tracking-[0.12em] text-text">
              CAP<span className="text-accent">CURVE</span>
            </div>
            <p className="mt-1 text-xs uppercase tracking-[0.18em] text-text-subtle">
              MLB career intelligence
            </p>
          </div>
          <div className="flex flex-col gap-2 text-xs sm:items-end">
            <div className="flex gap-4">
              <span>About</span>
              <span>Data Sources</span>
              <span>API</span>
            </div>
            <p>Data: MLB Stats API and Baseball Savant</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
