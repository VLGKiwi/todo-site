# stage 1: build
FROM golang:1.23.4-alpine AS build

# Работаем в backend директории
WORKDIR /app/backend

# Копируем go.mod из backend директории
COPY backend/go.mod .
RUN go mod download

# Копируем весь backend код
COPY backend/ .

# Собираем
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/todo-service ./cmd/todo-service/main.go

# stage 2: run
FROM alpine:latest
WORKDIR /app
COPY --from=build /bin/todo-service .
EXPOSE 8080
CMD [ "./todo-service" ]
