# --- Этап 1: Сборка ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем goose (утилиту для миграций)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Собираем статически слинкованный бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# --- Этап 2: Запуск ---
FROM alpine:3.20

# Устанавливаем корневые сертификаты (для HTTPS запросов) и curl (для healthcheck)
RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Копируем бинарник из первого этапа
COPY --from=builder /app/server .

# Копируем миграции и goose
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Открываем порт
EXPOSE 8080

# Запускаем сервис
CMD ["./server"]