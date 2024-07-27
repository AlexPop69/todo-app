run:
	docker compose up --build todo-app

test:
	go test -v ./...

migrate:
	goose -dir pkg/repository/migrations postgres "postgresql://postgres:postgres@localhost:5432?sslmode=disable" up

