import { Outlet, NavLink } from 'react-router-dom'

export default function Layout() {
  return (
    <div className="min-h-screen bg-surface flex flex-col">
      <header className="border-b border-white/5 px-6 py-4">
        <nav className="max-w-7xl mx-auto flex items-center justify-between">
          <NavLink to="/" className="font-display text-3xl tracking-wide text-white">
            CAP<span className="text-brand">CURVE</span>
          </NavLink>
          <div className="flex items-center gap-6 text-sm font-medium text-neutral">
            <NavLink
              to="/"
              className={({ isActive }) => isActive ? 'text-white' : 'hover:text-white transition-colors'}
            >
              Players
            </NavLink>
            <NavLink
              to="/leaderboards"
              className={({ isActive }) => isActive ? 'text-white' : 'hover:text-white transition-colors'}
            >
              Leaderboards
            </NavLink>
          </div>
        </nav>
      </header>
      <main className="flex-1 max-w-7xl mx-auto w-full px-6 py-8">
        <Outlet />
      </main>
      <footer className="border-t border-white/5 px-6 py-4 text-center text-xs text-neutral">
        CapCurve · MLB stats data via MLB Stats API
      </footer>
    </div>
  )
}
