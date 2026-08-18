DB_URL := postgres://tasker:tasker_password@localhost:5432/tasker_db?sslmode=disable
DB_TEST_URL := postgres://tasker:tasker_password@localhost:5432/tasker_db_test?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create migrate-test-db-up run run-test build-dev

log-db:
	docker exec -it tasker-postgres psql -U tasker -d tasker_db

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-test-db-up:
	migrate -path migrations -database "$(DB_TEST_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

run:
	go run ./cmd/server/main.go

run-test:
	go test ./...

build-dev:
	docker build -t tasker-api:dev ./backend
