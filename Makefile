# Makefile
# Structure adapted from https://lets-go-further.alexedwards.net/ (2025)

# ==================================================================================== #
# ENVIRONMENT & VARIABLES
# ==================================================================================== #

include .envrc

ECHO_PREFIX = [make]

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: Run the cmp/api application
.PHONY: run
run:
	go run ./cmd/api -db-dsn=${HOTEL_DB_DSN}

# ==================================================================================== #
# DATABASE MIGRATIONS
# ==================================================================================== #

## db/psql: Connect to the hotel database using psql as hotel_user
.PHONY: db/psql
db/psql:
	psql ${HOTEL_DB_DSN}

## db/migrations/new name=$1: Create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: Apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${HOTEL_DB_DSN} up

## db/migrations/down: Apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Reverting all migrations...'
	migrate -path ./migrations -database ${HOTEL_DB_DSN} down

## db/migrations/goto version=$1: Go to the specified migration version
.PHONY: db/migrations/goto
db/migrations/goto:
	@echo 'Going to schema migration version ${version}...'
	migrate -path ./migrations -database ${HOTEL_DB_DSN} goto ${version}

## db/migrations/fix version=$1: Force the schema_migrations table version
.PHONY: db/migrations/fix
db/migrations/fix:
	@echo 'Forcing schema migrations version to ${version}...'
	migrate -path ./migrations -database ${HOTEL_DB_DSN} force ${version}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: Tidy module dependencies and format all .go files
.PHONY: tidy
tidy:
	@echo '${ECHO_PREFIX} Tidying module dependencies...'
	go mod tidy
	@echo '${ECHO_PREFIX} Verifying and vendoring module dependencies...'
	go mod verify
# 	go mod vendor
	@echo '${ECHO_PREFIX} Formatting .go files...'
	go fmt ./...

## audit: Run quality control checks and tests
.PHONY: audit
audit:
	@echo '${ECHO_PREFIX} Checking module dependencies...'
	go mod tidy -diff
	go mod verify
	@echo '${ECHO_PREFIX} Vetting code...'
	go vet ./...
# 	go tool staticcheck ./...
	@echo '${ECHO_PREFIX} Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: Build the cmd/api application
.PHONY: build/api
build/api:
	@echo '${ECHO_PREFIX} Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# TESTS
# ==================================================================================== #

# GET
.PHONY: test/api/get
test/api/get:
	curl -i http://localhost:4000/v1/guests/A1234567

# POST
.PHONY: test/api/post
test/api/post:
	curl -i -X POST http://localhost:4000/v1/guests -d @test/01-post.json

# PUT
.PHONY: test/api/put
test/api/put:
	curl -i -X PUT http://localhost:4000/v1/guests/P0000000 -d @test/02-put.json

# PATCH
.PHONY: test/api/patch
test/api/patch:
	curl -i -X PATCH http://localhost:4000/v1/guests/P0000000 -d @test/03-patch.json

# DELETE
.PHONY: test/api/delete
test/api/delete:
	curl -i -X DELETE http://localhost:4000/v1/guests/P0000000