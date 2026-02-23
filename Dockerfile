FROM golang:alpine AS builder

WORKDIR /app

RUN apk add build-base sqlite

COPY . .

RUN CGO_ENABLED=1 go build
#RUN mkdir -p data && sqlite3 ./data/FGG.db < ./dbaccess/FGG.sql

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./FGG-Service"]
