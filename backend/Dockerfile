# stage 1: build
FROM golang:1.23.4-alpine AS build
WORKDIR /app

# Копируем ОБА файла зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем (путь правильный!)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/todo-service ./cmd/todo-service/main.go

# stage 2: run
FROM alpine:latest
WORKDIR /app
COPY --from=build /bin/todo-service .
EXPOSE 8080
CMD [ "./todo-service" ]
