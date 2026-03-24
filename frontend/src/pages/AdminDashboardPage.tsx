import { useQuery } from '@tanstack/react-query'
import { getAdminDashboard } from '@/api'

function formatTimestamp(value: string | null) {
  if (!value) {
    return 'Not yet recorded'
  }

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

export default function AdminDashboardPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['admin-dashboard'],
    queryFn: getAdminDashboard,
    refetchInterval: 60_000,
  })

  if (isLoading) {
    return <div className="text-neutral">Loading sync dashboard...</div>
  }

  if (isError || !data) {
    return <div className="text-overpaid">Unable to load admin dashboard.</div>
  }

  const { sync_status: syncStatus } = data

  return (
    <div className="space-y-8">
      <section className="rounded-3xl border border-white/10 bg-white/[0.03] p-8">
        <p className="text-xs uppercase tracking-[0.3em] text-brand">Admin</p>
        <h1 className="mt-3 font-display text-5xl text-white">Sync Dashboard</h1>
        <p className="mt-3 max-w-2xl text-neutral">
          Read-only visibility into the scheduler so we can confirm active-player refreshes are running on time.
        </p>
      </section>

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-5">
          <p className="text-sm text-neutral">Total players</p>
          <p className="mt-2 text-3xl font-semibold text-white">{data.total_players.toLocaleString()}</p>
        </div>
        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-5">
          <p className="text-sm text-neutral">Active players</p>
          <p className="mt-2 text-3xl font-semibold text-white">{data.active_players.toLocaleString()}</p>
        </div>
        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-5">
          <p className="text-sm text-neutral">Last sync</p>
          <p className="mt-2 text-lg font-semibold text-white">{formatTimestamp(syncStatus.last_sync_completed_at)}</p>
        </div>
        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-5">
          <p className="text-sm text-neutral">Next scheduled sync</p>
          <p className="mt-2 text-lg font-semibold text-white">{formatTimestamp(syncStatus.next_scheduled_sync_at)}</p>
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-[1.5fr_1fr]">
        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-6">
          <h2 className="text-xl font-semibold text-white">Scheduler status</h2>
          <dl className="mt-4 grid gap-4 sm:grid-cols-2">
            <div>
              <dt className="text-sm text-neutral">Enabled</dt>
              <dd className="mt-1 text-white">{syncStatus.enabled ? 'Yes' : 'No'}</dd>
            </div>
            <div>
              <dt className="text-sm text-neutral">Running now</dt>
              <dd className="mt-1 text-white">{syncStatus.running ? 'Yes' : 'No'}</dd>
            </div>
            <div>
              <dt className="text-sm text-neutral">Cadence mode</dt>
              <dd className="mt-1 text-white">{syncStatus.in_season ? 'In-season daily' : 'Off-season weekly'}</dd>
            </div>
            <div>
              <dt className="text-sm text-neutral">Target season year</dt>
              <dd className="mt-1 text-white">{syncStatus.target_season_year || 'Not yet scheduled'}</dd>
            </div>
            <div>
              <dt className="text-sm text-neutral">Last successful sync</dt>
              <dd className="mt-1 text-white">{formatTimestamp(syncStatus.last_successful_sync_at)}</dd>
            </div>
            <div>
              <dt className="text-sm text-neutral">Last run started</dt>
              <dd className="mt-1 text-white">{formatTimestamp(syncStatus.last_sync_started_at)}</dd>
            </div>
          </dl>
        </div>

        <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-6">
          <h2 className="text-xl font-semibold text-white">Ingestion errors</h2>
          <p className="mt-4 text-sm leading-6 text-neutral">
            {syncStatus.last_error || 'No ingestion errors recorded in memory.'}
          </p>
        </div>
      </section>
    </div>
  )
}
