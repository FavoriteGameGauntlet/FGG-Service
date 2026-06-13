FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o fgg-server

ENTRYPOINT ["./fgg-server"]
