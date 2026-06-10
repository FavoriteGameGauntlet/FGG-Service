FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o fgg-server

ENTRYPOINT ["./fgg-server"]
