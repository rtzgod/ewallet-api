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

.PHONY: clean
clean:
	rm coverage.html
	rm cover.out