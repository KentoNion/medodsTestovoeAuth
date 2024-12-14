FROM golang:1.23.4-alpine AS builder

WORKDIR /usr/local/src

# Копируем зависимости
COPY authApp/go.mod authApp/go.sum ./
RUN go mod download

# Копируем исходники и и билдим
COPY authApp ./
RUN go build -ldflags="-s -w" -o /app cmd/main.go

# Стадия выполнения
FROM scratch AS runner

WORKDIR /root/

# Добавляем бинарный файл и миграции с конфигом
COPY --from=builder /app ./app
COPY authApp/config.yaml ./
COPY authApp/migrations ./migrations

CMD ["./app"]

#можно было бы в теории и тесты из докер контейнера вырезать с применением .dockerignore, но пусть лежат