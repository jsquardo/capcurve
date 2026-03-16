# CapCurve ⚾

MLB career arc visualizer and contract value tracker.

## Stack

- **Backend**: Go + Echo + GORM + PostgreSQL
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Dev**: Docker Compose + VS Code Dev Containers

## Getting Started

### Prerequisites
- Docker Desktop
- VS Code with Dev Containers extension (optional but recommended)

### Option A: Dev Containers (recommended)
1. Open the project in VS Code
2. When prompted, click **"Reopen in Container"**
3. VS Code will build and start all containers automatically
4. The Go API and React frontend will be running with hot reload

### Option B: Docker Compose directly
```bash
cp .env.example .env
make up
make migrate
```

## Services

| Service | URL |
|---|---|
| React Frontend | http://localhost:5173 |
| Go API | http://localhost:8080 |
| API Health | http://localhost:8080/health |
| PostgreSQL | localhost:5432 |

## Common Commands

```bash
make up            # Start all containers
make down          # Stop all containers
make logs          # Follow all logs
make migrate       # Run pending migrations
make migrate-down  # Roll back one migration
make shell-backend # Shell into Go container
make shell-db      # psql into PostgreSQL
make test          # Run Go tests
```

## Project Structure

```
capcurve/
├── backend/
│   ├── cmd/server/         # Entry point
│   ├── internal/
│   │   ├── config/         # Environment config
│   │   ├── database/       # DB connection
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Echo middleware
│   │   ├── models/         # GORM models
│   │   ├── services/       # Business logic
│   │   └── jobs/           # Weekly sync workers
│   └── db/migrations/      # SQL migrations
├── frontend/
│   └── src/
│       ├── api/            # Typed API client
│       ├── components/     # Shared UI components
│       ├── pages/          # Route-level pages
│       ├── hooks/          # Custom React hooks
│       └── types/          # TypeScript types
├── .devcontainer/          # VS Code Dev Container config
├── docker-compose.yml
└── Makefile
```

## Database Schema

- `players`: Core player identity and MLB metadata
- `season_stats`: One active row per player-season for persisted hitting/pitching stats, team metadata, and computed `value_score`
- `contracts`: Contract headers such as total value, years, and signing team
- `contract_seasons`: Per-season breakdown rows for each contract
- `career_arcs`: Precomputed career-arc summaries used by the visualizer

`season_stats` active-row uniqueness is enforced in PostgreSQL, not just in application code. The database uses a partial unique index on `(player_id, year) WHERE deleted_at IS NULL`, which allows soft-deleted tombstones to remain for ingestion history without permitting duplicate active rows for the same player-season.
