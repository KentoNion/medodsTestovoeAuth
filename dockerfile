FROM golang:1.23.4 AS builder

RUN go version
ENV GOPATH=/

COPY ./authApp ./

RUN go mod download
RUN go build -o app ./cmd/main.go

cmd ["./app"]