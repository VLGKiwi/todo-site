# stage 1: build
FROM golang:1.23.4-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/todo-service ./cmd/todo-service/main.go

# stage 2: run
FROM alpine:latest
WORKDIR /app

# Устанавливаем curl для health check
RUN apk --no-cache add curl

# Копируем бинарник
COPY --from=build /bin/todo-service .

# Создаем директорию для данных (важно для Render!)
RUN mkdir -p /data

# Создаем простой health check скрипт
RUN echo '#!/bin/sh' > /app/health-check.sh && \
    echo 'curl -f http://localhost:8080/health || exit 1' >> /app/health-check.sh && \
    chmod +x /app/health-check.sh

EXPOSE 8080

# Убедись, что переменные окружения используются в коде
ENV PORT=8080
ENV DATA_PATH=/data/todo.db

CMD [ "./todo-service" ]
