# ==========================================
# Stage 1: Build Static Go Binary
# ==========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Cache go modules
COPY go.mod ./
# If go.sum exists, copy it as well
RUN go mod download

# Copy source code
COPY . .

# Build statically linked, stripped production binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=1.0.0" \
    -o /app/server cmd/server/main.go

# ==========================================
# Stage 2: Final Minimal Runtime Image
# ==========================================
FROM alpine:3.20 AS runner

WORKDIR /app

# Install security certificates and curl for healthchecks
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

# Copy compiled binary from builder
COPY --from=builder /app/server /app/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create storage directory for code workspaces and set ownership
RUN mkdir -p /app/storage/workspaces && \
    chown -R appuser:appgroup /app

# Set ownership to non-root user
USER appuser:appgroup

# Expose default HTTP port
EXPOSE 8080

# Configure Health Check
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Start Go Server
ENTRYPOINT ["/app/server"]
