FROM golang:1.23-alpine AS builder
WORKDIR /app


RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache postgresql-client


COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY .env .
COPY ./sql/schema /app/sql

COPY wait-for-db.sh /app/wait-for-db.sh
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/wait-for-db.sh /app/entrypoint.sh

EXPOSE ${PORT}

CMD ["/app/main"]
