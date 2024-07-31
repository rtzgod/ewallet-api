test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format short-verbose

build:
	docker compose build

run:
	docker compose up

coverage:
	go test ./... -coverprofile=cover.out | go tool cover -html=cover.out -o coverage.html | firefox coverage.html

migrateUp:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' up

migrateDown:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' down

clean:
	rm coverage.html
	rm cover.out