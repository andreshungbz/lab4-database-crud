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

# ==================================================================================== #
# Rate Limiting Middleware
# ==================================================================================== #

# Rate Limiting
.PHONY: test/rate-limiting-ip
test/rate-limiting-ip:
	for i in {1..6}; do curl http://localhost:4000/v1/healthcheck; done

# ==================================================================================== #
# Guest Model
# ==================================================================================== #

# GET
.PHONY: test/api/guests/get
test/api/guests/get:
	curl -i http://localhost:4000/v1/guests/A1234567

# GET ALL
.PHONY: test/api/guests/get-all
test/api/guests/get-all:
	curl -i http://localhost:4000/v1/guests

# GET ALL (filtered by name)
.PHONY: test/api/guests/get-all-name
test/api/guests/get-all-name:
	curl -i http://localhost:4000/v1/guests?name=mae

# GET ALL (filters, pagination and sorting) (note that quotes are necessary)
.PHONY: test/api/guests/get-all-filters
test/api/guests/get-all-filters:
	curl -i "http://localhost:4000/v1/guests?country=belize&page=1&page_size=3&sort=-passport_number"

# POST
.PHONY: test/api/guests/post
test/api/guests/post:
	curl -i -X POST http://localhost:4000/v1/guests -d @test/guest/01-post.json

# PUT
.PHONY: test/api/guests/put
test/api/guests/put:
	curl -i -X PUT http://localhost:4000/v1/guests/P0000000 -d @test/guest/02-put.json

# PATCH
.PHONY: test/api/guests/patch
test/api/guests/patch:
	curl -i -X PATCH http://localhost:4000/v1/guests/P0000000 -d @test/guest/03-patch.json

# DELETE
.PHONY: test/api/guests/delete
test/api/guests/delete:
	curl -i -X DELETE http://localhost:4000/v1/guests/P0000000

# ==================================================================================== #
# Hotel Model
# ==================================================================================== #

# GET
.PHONY: test/api/hotels/get
test/api/hotels/get:
	curl -i http://localhost:4000/v1/hotels/1

# GET ALL
.PHONY: test/api/hotels/get-all
test/api/hotels/get-all:
	curl -i http://localhost:4000/v1/hotels

# GET ALL (filtered by name)
.PHONY: test/api/hotels/get-all-name
test/api/hotels/get-all-name:
	curl -i http://localhost:4000/v1/hotels?name=grand

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/hotels/get-all-filters
test/api/hotels/get-all-filters:
	curl -i "http://localhost:4000/v1/hotels?city=belize&page=1&page_size=3&sort=-name"

# POST
.PHONY: test/api/hotels/post
test/api/hotels/post:
	curl -i -X POST http://localhost:4000/v1/hotels -d @test/hotel/01-post.json

# PUT
.PHONY: test/api/hotels/put
test/api/hotels/put:
	curl -i -X PUT http://localhost:4000/v1/hotels/3 -d @test/hotel/02-put.json

# PATCH
.PHONY: test/api/hotels/patch
test/api/hotels/patch:
	curl -i -X PATCH http://localhost:4000/v1/hotels/3 -d @test/hotel/03-patch.json

# DELETE
.PHONY: test/api/hotels/delete
test/api/hotels/delete:
	curl -i -X DELETE http://localhost:4000/v1/hotels/3

# ==================================================================================== #
# Department Model
# ==================================================================================== #

# GET
.PHONY: test/api/departments/get
test/api/departments/get:
	curl -i "http://localhost:4000/v1/departments/Hotel%20Operations"

# GET ALL
.PHONY: test/api/departments/get-all
test/api/departments/get-all:
	curl -i http://localhost:4000/v1/departments

# GET ALL (filtered by name)
.PHONY: test/api/departments/get-all-name
test/api/departments/get-all-name:
	curl -i http://localhost:4000/v1/departments?name=housekeeping

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/departments/get-all-filters
test/api/departments/get-all-filters:
	curl -i "http://localhost:4000/v1/departments?page=1&page_size=3&sort=-budget"

# POST
.PHONY: test/api/departments/post
test/api/departments/post:
	curl -i -X POST http://localhost:4000/v1/departments -d @test/department/01-post.json

# PUT
.PHONY: test/api/departments/put
test/api/departments/put:
	curl -i -X PUT "http://localhost:4000/v1/departments/Restaurant%20Operations" -d @test/department/02-put.json

# PATCH
.PHONY: test/api/departments/patch
test/api/departments/patch:
	curl -i -X PATCH "http://localhost:4000/v1/departments/Restaurant%20Operations" -d @test/department/03-patch.json

# DELETE
.PHONY: test/api/departments/delete
test/api/departments/delete:
	curl -i -X DELETE "http://localhost:4000/v1/departments/Restaurant%20Operations"

# ==================================================================================== #
# Employee Model
# ==================================================================================== #

# GET
.PHONY: test/api/employees/get
test/api/employees/get:
	curl -i http://localhost:4000/v1/employees/1

# GET ALL
.PHONY: test/api/employees/get-all
test/api/employees/get-all:
	curl -i http://localhost:4000/v1/employees

# GET ALL (filtered by role)
.PHONY: test/api/employees/get-all-role
test/api/employees/get-all-role:
	curl -i "http://localhost:4000/v1/employees?role=operations_manager"

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/employees/get-all-filters
test/api/employees/get-all-filters:
	curl -i "http://localhost:4000/v1/employees?page=1&page_size=2&sort=-id"

# POST Operations Manager
.PHONY: test/api/employees/post-om
test/api/employees/post-om:
	curl -i -X POST http://localhost:4000/v1/employees -d @test/employee/01-post-om.json

# POST Front Desk
.PHONY: test/api/employees/post-fd
test/api/employees/post-fd:
	curl -i -X POST http://localhost:4000/v1/employees -d @test/employee/02-post-fd.json

# POST Housekeeper
.PHONY: test/api/employees/post-hk
test/api/employees/post-hk:
	curl -i -X POST http://localhost:4000/v1/employees -d @test/employee/03-post-hk.json

# PUT
.PHONY: test/api/employees/put
test/api/employees/put:
	curl -i -X PUT http://localhost:4000/v1/employees/14 -d @test/employee/04-put.json

# PATCH
.PHONY: test/api/employees/patch
test/api/employees/patch:
	curl -i -X PATCH http://localhost:4000/v1/employees/14 -d @test/employee/05-patch.json

# DELETE
.PHONY: test/api/employees/delete
test/api/employees/delete:
	curl -i -X DELETE http://localhost:4000/v1/employees/14

# ==================================================================================== #
# Room Model
# ==================================================================================== #

# GET
.PHONY: test/api/rooms/get
test/api/rooms/get:
	curl -i http://localhost:4000/v1/hotels/1/rooms/101

# GET ALL
.PHONY: test/api/rooms/get-all
test/api/rooms/get-all:
	curl -i http://localhost:4000/v1/hotels/1/rooms

# GET ALL (filters, pagination and sorting) (note that quotes are necessary)
.PHONY: test/api/rooms/get-all-filters
test/api/rooms/get-all-filters:
	curl -i "http://localhost:4000/v1/hotels/1/rooms?page=1&page_size=3&sort=-number"

# POST
.PHONY: test/api/rooms/post
test/api/rooms/post:
	curl -i -X POST http://localhost:4000/v1/hotels/1/rooms -d @test/room/01-post.json

# PUT
.PHONY: test/api/rooms/put
test/api/rooms/put:
	curl -i -X PUT http://localhost:4000/v1/hotels/1/rooms/501 -d @test/room/02-put.json

# PATCH
.PHONY: test/api/rooms/patch
test/api/rooms/patch:
	curl -i -X PATCH http://localhost:4000/v1/hotels/1/rooms/501 -d @test/room/03-patch.json

# DELETE
.PHONY: test/api/rooms/delete
test/api/rooms/delete:
	curl -i -X DELETE http://localhost:4000/v1/hotels/1/rooms/501

# ==================================================================================== #
# RoomType Model
# ==================================================================================== #

# GET
.PHONY: test/api/room_types/get
test/api/room_types/get:
	curl -i http://localhost:4000/v1/room_types/1

# GET ALL
.PHONY: test/api/room_types/get-all
test/api/room_types/get-all:
	curl -i http://localhost:4000/v1/room_types

# GET ALL (filtered by title)
.PHONY: test/api/room_types/get-all-title
test/api/room_types/get-all-title:
	curl -i http://localhost:4000/v1/room_types?title=suite

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/room_types/get-all-filters
test/api/room_types/get-all-filters:
	curl -i "http://localhost:4000/v1/room_types?page=1&page_size=3&sort=-base_rate"

# POST
.PHONY: test/api/room_types/post
test/api/room_types/post:
	curl -i -X POST http://localhost:4000/v1/room_types -d @test/room_type/01-post.json

# PUT
.PHONY: test/api/room_types/put
test/api/room_types/put:
	curl -i -X PUT http://localhost:4000/v1/room_types/4 -d @test/room_type/02-put.json

# PATCH
.PHONY: test/api/room_types/patch
test/api/room_types/patch:
	curl -i -X PATCH http://localhost:4000/v1/room_types/4 -d @test/room_type/03-patch.json

# DELETE
.PHONY: test/api/room_types/delete
test/api/room_types/delete:
	curl -i -X DELETE http://localhost:4000/v1/room_types/4

# ==================================================================================== #
# HousekeepingTask Model
# ==================================================================================== #

# GET
.PHONY: test/api/housekeeping_tasks/get
test/api/housekeeping_tasks/get:
	curl -i http://localhost:4000/v1/housekeeping_tasks/1

# GET ALL (tasks for a room)
.PHONY: test/api/housekeeping_tasks/get-all
test/api/housekeeping_tasks/get-all:
	curl -i http://localhost:4000/v1/hotels/1/rooms/101/housekeeping_tasks

# GET ALL (filtered by housekeeper)
.PHONY: test/api/housekeeping_tasks/get-all-housekeeper
test/api/housekeeping_tasks/get-all-housekeeper:
	curl -i "http://localhost:4000/v1/hotels/1/rooms/101/housekeeping_tasks?housekeeper_id=3"

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/housekeeping_tasks/get-all-filters
test/api/housekeeping_tasks/get-all-filters:
	curl -i "http://localhost:4000/v1/hotels/1/rooms/101/housekeeping_tasks?page=1&page_size=1&sort=-created_at"

# POST
.PHONY: test/api/housekeeping_tasks/post
test/api/housekeeping_tasks/post:
	curl -i -X POST http://localhost:4000/v1/hotels/1/rooms/301/housekeeping_tasks -d @test/housekeeping_task/01-post.json

# PUT
.PHONY: test/api/housekeeping_tasks/put
test/api/housekeeping_tasks/put:
	curl -i -X PUT http://localhost:4000/v1/housekeeping_tasks/6 -d @test/housekeeping_task/02-put.json

# PATCH
.PHONY: test/api/housekeeping_tasks/patch
test/api/housekeeping_tasks/patch:
	curl -i -X PATCH http://localhost:4000/v1/housekeeping_tasks/6 -d @test/housekeeping_task/03-patch.json

# DELETE
.PHONY: test/api/housekeeping_tasks/delete
test/api/housekeeping_tasks/delete:
	curl -i -X DELETE http://localhost:4000/v1/housekeeping_tasks/6

# ==================================================================================== #
# MaintenanceReport Model
# ==================================================================================== #

# GET
.PHONY: test/api/maintenance_reports/get
test/api/maintenance_reports/get:
	curl -i http://localhost:4000/v1/maintenance_reports/1

# GET ALL (for a room)
.PHONY: test/api/maintenance_reports/get-all
test/api/maintenance_reports/get-all:
	curl -i http://localhost:4000/v1/hotels/1/rooms/101/maintenance_reports

# POST
.PHONY: test/api/maintenance_reports/post
test/api/maintenance_reports/post:
	curl -i -X POST http://localhost:4000/v1/hotels/1/rooms/101/maintenance_reports -d @test/maintenance_report/01-post.json

# PUT
.PHONY: test/api/maintenance_reports/put
test/api/maintenance_reports/put:
	curl -i -X PUT http://localhost:4000/v1/maintenance_reports/5 -d @test/maintenance_report/02-put.json

# PATCH
.PHONY: test/api/maintenance_reports/patch
test/api/maintenance_reports/patch:
	curl -i -X PATCH http://localhost:4000/v1/maintenance_reports/5 -d @test/maintenance_report/03-patch.json

# DELETE
.PHONY: test/api/maintenance_reports/delete
test/api/maintenance_reports/delete:
	curl -i -X DELETE http://localhost:4000/v1/maintenance_reports/5

# ==================================================================================== #
# Registration Model
# ==================================================================================== #

# GET
.PHONY: test/api/registrations/get
test/api/registrations/get:
	curl -i http://localhost:4000/v1/registrations/2/1/102

# GET ALL (registrations for a reservation)
.PHONY: test/api/registrations/get-all
test/api/registrations/get-all:
	curl -i http://localhost:4000/v1/registrations/2

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/registrations/get-all-filters
test/api/registrations/get-all-filters:
	curl -i "http://localhost:4000/v1/registrations/2?page=1&page_size=2&sort=-room_number"

# POST
.PHONY: test/api/registrations/post
test/api/registrations/post:
	curl -i -X POST http://localhost:4000/v1/registrations -d @test/registration/01-post.json

# PUT
.PHONY: test/api/registrations/put
test/api/registrations/put:
	curl -i -X PUT http://localhost:4000/v1/registrations/2/1/102 -d @test/registration/02-put.json

# PATCH
.PHONY: test/api/registrations/patch
test/api/registrations/patch:
	curl -i -X PATCH http://localhost:4000/v1/registrations/2/1/402 -d @test/registration/03-patch.json

# DELETE
.PHONY: test/api/registrations/delete
test/api/registrations/delete:
	curl -i -X DELETE http://localhost:4000/v1/registrations/2/1/102

# ==================================================================================== #
# Registration Model
# ==================================================================================== #

# GET
.PHONY: test/api/reservations/get
test/api/reservations/get:
	curl -i http://localhost:4000/v1/reservations/2

# GET ALL (reservations)
.PHONY: test/api/reservations/get-all
test/api/reservations/get-all:
	curl -i http://localhost:4000/v1/reservations

# GET ALL (filters, pagination and sorting)
.PHONY: test/api/reservations/get-all-filters
test/api/reservations/get-all-filters:
	curl -i "http://localhost:4000/v1/reservations?page=1&page_size=2&sort=-id"

# POST
.PHONY: test/api/reservations/post
test/api/reservations/post:
	curl -i -X POST http://localhost:4000/v1/reservations -d @test/reservation/01-post.json

# PUT
.PHONY: test/api/reservations/put
test/api/reservations/put:
	curl -i -X PUT http://localhost:4000/v1/reservations/7 -d @test/reservation/02-put.json

# PATCH
.PHONY: test/api/reservations/patch
test/api/reservations/patch:
	curl -i -X PATCH http://localhost:4000/v1/reservations/7 -d @test/reservation/03-patch.json

# DELETE
.PHONY: test/api/reservations/delete
test/api/reservations/delete:
	curl -i -X DELETE http://localhost:4000/v1/reservations/7
