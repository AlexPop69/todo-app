FROM golang:1.22

WORKDIR /app

ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build -o todo-app ./cmd/main.go

CMD ["./todo-app"]


