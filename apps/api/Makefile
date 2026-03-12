# Define the connection string once as a variable
DB_URL="postgresql://hackersquare-admin:password@localhost:5432/hackersquare-dev-db?sslmode=disable"

dev:
	@air

build:
	@go build -o bin/api ./cmd/api

start:
	@./bin/api

format:
	@go fmt ./...

db-it:
	@docker exec -it postgres psql -U hackersquare-admin -d hackersquare-dev-db

# 1. Applies all pending migrations
db-migrate:
	@migrate -path ./migrations -database $(DB_URL) up

# 2. Rolls back the very last migration (undo button)
db-rollback:
	@migrate -path ./migrations -database $(DB_URL) down 1
	
# Usage: make db-force version=2
db-force:
	@migrate -path ./migrations -database $(DB_URL) force $(version)

# 3. Generates new blank up/down migration files
# Usage: make db-create-migration name=add_users_table
db-create-migration:
	@migrate create -ext sql -dir ./migrations -seq $(name)

