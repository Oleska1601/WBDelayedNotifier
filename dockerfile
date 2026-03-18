
FROM golang:alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./notifier ./cmd/app/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/notifier ./notifier
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/internal/repo/postgres/migrations/ ./internal/repo/postgres/migrations/*.sql

COPY --from=builder /app/web ./web

EXPOSE 8081

CMD ["./notifier"]
