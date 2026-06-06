FROM golang:alpine AS builder

WORKDIR /app

RUN apk add build-base sqlite

COPY . .

RUN CGO_ENABLED=1 go build -o fgg-server

ENTRYPOINT ["./fgg-server"]
