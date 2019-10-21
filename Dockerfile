FROM golang:1.12.10-alpine AS builder
RUN apk add --no-cache git gcc libc-dev
RUN mkdir -p /go/src/github.com/fufuhu
RUN go get -v github.com/rubenv/sql-migrate/...
COPY src/github.com/fufuhu src/github.com/fufuhu
RUN go build -o sql-web-migrate src/github.com/fufuhu/sql-web-migrate/main.go

FROM alpine:3
COPY --from=builder /go/sql-web-migrate /usr/local/bin/sql-web-migrate
RUN mkdir -p /etc/migrate
COPY conf.d /etc/migrate