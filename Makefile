.PHONY: test
test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format short-verbose
.PHONY: build
build:
	docker compose build
.PHONY: run
run:
	docker compose up
.PHONY: coverage
coverage:
	go test ./... -coverprofile=cover.out | go tool cover -html=cover.out -o coverage.html | firefox coverage.html

.PHONY: migrateUp
migrateUp:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' up
.PHONY: migrateDown
migrateDown:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' down

.PHONY: clean
clean:
	rm coverage.html
	rm cover.out