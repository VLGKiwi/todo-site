FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-service ./cmd/todo-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/todo-service .
EXPOSE 8080
CMD ["./todo-service"]
