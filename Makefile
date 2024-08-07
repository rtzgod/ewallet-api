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
.PHONY: up
up:
	migrate -path ./path/for/your/migrations/folder -database 'postgres://POSTGRES_USER:POSTGRES_PASSWORD@HOST:PORT/DB_NAME?sslmode=disable' up
.PHONY: down
down:
	migrate -path ./path/for/your/migrations/folder -database 'postgres://POSTGRES_USER:POSTGRES_PASSWORD@HOST:PORT/DB_NAME?sslmode=disable' down
