run docker:
	docker compose up --build todo-app

run local:
	go build -o todo-app ./cmd/main.go
	./todo-app

test:
	go test -v ./...

migrate:
	goose -dir pkg/repository/migrations postgres "postgresql://postgres:postgres@localhost:5432?sslmode=disable" up

swag:
	swag init -g cmd/main.go