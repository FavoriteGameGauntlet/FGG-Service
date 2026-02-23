FROM golang:alpine AS builder

WORKDIR /app

RUN apk add build-base sqlite

COPY . .

RUN CGO_ENABLED=1 go build

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./FGG-Service"]
