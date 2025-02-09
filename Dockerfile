FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server/main.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Install PostgreSQL client for goose migrations
RUN apk add --no-cache postgresql-client

# Copy the built binary and other files
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /app/goose
COPY .env .
COPY ./sql/schema /app/sql

# Copy scripts and make them executable
COPY wait-for-db.sh /app/wait-for-db.sh
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/wait-for-db.sh /app/entrypoint.sh /app/goose

# Expose the port
EXPOSE ${PORT}

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]