.PHONY: build build-debug up up-debug down logs dbshell start-db stop-db test ui-build ui-up ui-down

# Default build (production)
build:
	docker compose build

# Build in debug mode
build-debug:
	docker compose build --build-arg DEBUG=true --build-arg GIN_MODE=debug

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f app

dbshell:
	docker compose exec db psql -U $(DB_USER) -d $(DB_NAME)

start-db:
	docker compose up -d db

stop-db:
	docker compose stop db

test:
	go test ./... -v

ui-build:
	@echo "Building UI..."
	cd ui && npm install && npm run build

ui-up:
	cd ui && npm run dev
