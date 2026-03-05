.PHONY: up down build logs migrate migrate-down shell-backend shell-db reset

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

migrate:
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" up

migrate-down:
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" down 1

migrate-status:
	docker compose exec backend migrate -path /app/db/migrations -database "$$DATABASE_URL" version

shell-backend:
	docker compose exec backend sh

shell-db:
	docker compose exec db psql -U capcurve -d capcurve_development

reset:
	docker compose down -v
	docker compose up -d

test:
	docker compose exec backend go test ./...

tidy:
	docker compose exec backend go mod tidy
