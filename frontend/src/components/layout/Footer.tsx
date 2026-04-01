export default function Footer() {
  return (
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
            {/* TODO: wire to /about when About page exists */}
            <span>About</span>
            {/* TODO: wire to /data-sources when that page exists */}
            <span>Data Sources</span>
            {/* TODO: wire to /api when API docs page exists */}
            <span>API</span>
          </div>
          <p>Data: MLB Stats API and Baseball Savant</p>
        </div>
      </div>
    </footer>
  )
}
