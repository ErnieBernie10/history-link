# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate

RUN go build -o main cmd/api/main.go

# Production stage
FROM alpine:3.20.1 AS prod
WORKDIR /app

# Install curl
RUN apk add --no-cache curl

RUN curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
RUN chmod +x /usr/local/bin/dbmate

# Copy built binary and migration script
COPY --from=build /app/main /app/main
COPY migrate-and-run.sh /app/migrate-and-run.sh

# Ensure the script is executable
RUN chmod +x /app/migrate-and-run.sh

# Set the entrypoint to run migrations first, then start the app
ENTRYPOINT ["/app/migrate-and-run.sh"]
