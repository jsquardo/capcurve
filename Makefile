.PHONY: up up-d down build logs logs-backend migrate migrate-down migrate-status \
        shell-backend shell-db reset test tidy go-build \
        migrate-local migrate-down-local db-local

# ─── Docker detection ───────────────────────────────────────────────────────
# IN_CONTAINER is true when running inside the backend Docker container.
# Used to switch between `docker compose exec` (outside) and direct commands (inside).
IN_CONTAINER := $(shell [ -f /.dockerenv ] && echo true || echo false)

# ─── Compose targets (always run from Mac/host) ─────────────────────────────
up:
	docker compose up

up-d:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

logs-backend:
	docker compose logs -f backend

shell-backend:
	docker compose exec backend sh

shell-db:
	docker compose exec db psql -U capcurve -d capcurve_development

reset:
	docker compose down -v
	docker compose up -d

# ─── Go targets (work from both inside and outside the container) ────────────
test:
ifeq ($(IN_CONTAINER),true)
	cd backend && GOCACHE=/tmp/capcurve-gocache go test -mod=mod ./...
else
	docker compose exec backend sh -c "cd /app && GOCACHE=/tmp/capcurve-gocache go test -mod=mod ./..."
endif

tidy:
ifeq ($(IN_CONTAINER),true)
	cd backend && go mod tidy
else
	docker compose exec backend go mod tidy
endif

go-build:
ifeq ($(IN_CONTAINER),true)
	cd backend && go build ./...
else
	docker compose exec backend go build ./...
endif

go-build-local:
	cd backend && go build ./...

# ─── Migration targets (work from both inside and outside the container) ─────
migrate:
ifeq ($(IN_CONTAINER),true)
	migrate -path /app/db/migrations -database "$(DATABASE_URL)" up
else
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" up
endif

migrate-down:
ifeq ($(IN_CONTAINER),true)
	migrate -path /app/db/migrations -database "$(DATABASE_URL)" down 1
else
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" down 1
endif

migrate-status:
ifeq ($(IN_CONTAINER),true)
	migrate -path /app/db/migrations -database "$(DATABASE_URL)" version
else
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" version
endif

# ─── DB shell (inside container only) ───────────────────────────────────────
db-local:
	psql $(DATABASE_URL)

# ─── Legacy aliases (kept for backwards compatibility) ───────────────────────
migrate-local: migrate
migrate-down-local: migrate-down