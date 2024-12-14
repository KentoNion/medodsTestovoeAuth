FROM golang:1.23.4-alpine AS builder

WORKDIR /usr/local/src

COPY ["authApp/go.mod", "authApp/go.sum", "./"]
RUN go mod download

#build
COPY authApp ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app ./
COPY authApp/config.yaml ./
COPY authApp/migrations ./migrations

CMD ["/app"]