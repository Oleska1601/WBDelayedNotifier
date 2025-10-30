# Используем официальный образ Go
FROM golang:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY ./go.mod ./go.sum ./
# Загружаем зависимости
RUN go mod download

# Копируем весь код приложения
COPY . .

RUN go build -o ./blogs ./cmd/app/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/blogs ./blogs
COPY --from=builder /app/config/config.yml ./config/config.yml
COPY --from=builder /app/internal/database/repo/migrations/ ./internal/database/repo/sql/*.sql

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./blogs"]
