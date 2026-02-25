FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make tar gzip

# Copy project files
COPY . .

# Build the binary
RUN go mod download
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache tar gzip mariadb-client dcron

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/flux-backup /app/flux-backup

# Create backup directory
RUN mkdir -p /backups

# Create a cron job
RUN echo "0 2 * * * /app/flux-backup -fivem /mnt/fivem -db-name fivem -db-user backup -db-pass \$DB_PASSWORD -backup-dir /backups -retention 30" > /etc/crontabs/root

# Entrypoint
ENTRYPOINT ["crond", "-f"]
