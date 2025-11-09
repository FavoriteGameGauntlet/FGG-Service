FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .
RUN apk add build-base sqlite
RUN CGO_ENABLED=1 go build
RUN sqlite3 ./FGG.db < ./database/FGG.sql

ENTRYPOINT ["./FGG-Service"]
