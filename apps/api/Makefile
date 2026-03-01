dev:
	@air

build:
	@go build -o bin/api ./cmd/api

start:
	@./bin/api

db-it:
	@docker exec -it postgres psql -U hackersquare-admin -d hackersquare-dev-db

db-migrate:
	@migrate -path ./migrations -database "postgresql://hackersquare-admin:password@localhost:5432/hackersquare-dev-db?sslmode=disable" up 