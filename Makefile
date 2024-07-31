test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format short-verbose

coverage:
	go test ./... -coverprofile=cover.out | go tool cover -html=cover.out -o coverage.html | firefox coverage.html

postgres:
	sudo docker run --name=ewallet -e POSTGRES_PASSWORD=1 -p 5436:5432 -d --rm postgres

migrateUp:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' up

migrateDown:
	migrate -path ./schema -database 'postgres://postgres:1@localhost:5436/postgres?sslmode=disable' down

dockerMigrateUp:
	migrate -path ./schema -database 'postgres://postgres:1@db:5432/postgres?sslmode=disable' up

dockerMigrateDown:
	migrate -path ./schema -database 'postgres://postgres:1@db:5432/postgres?sslmode=disable' down

clean:
	rm coverage.html
	rm cover.out