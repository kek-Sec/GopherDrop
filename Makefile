.PHONY: build up down logs dbshell start-db stop-db test ui-build ui-up ui-down

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

dbshell:
	docker-compose exec db psql -U $(DB_USER) -d $(DB_NAME)

start-db:
	docker-compose up -d db

stop-db:
	docker-compose stop db

test:
	go test ./... -v

ui-build:
	@echo "Building UI..."
	cd ui && npm install && npm run build

ui-up:
	cd ui && npm run dev
