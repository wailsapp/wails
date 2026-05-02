# Wails Server Mode Dockerfile
# Multi-stage build for minimal image size
#
# BUILD FROM v3 ROOT (to use local code before server mode is published):
#   docker build -t server-example -f examples/server/Dockerfile .
#
# BUILD FROM EXAMPLE DIR (after server mode is published):
#   cd examples/server && docker build -t server-example .

# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy source code
COPY . .

# If building from example dir, remove replace directive
# If building from v3 root, this is a no-op
RUN sed -i '/^replace/d' go.mod 2>/dev/null || true

# Download dependencies
RUN go mod tidy

# Build the server binary and prepare assets
# When building from v3 root, change to example dir
RUN if [ -d "examples/server" ]; then \
      cd examples/server && go build -tags server -ldflags="-s -w" -o /app/server . && \
      cp -r frontend/dist /app/frontend_dist; \
    else \
      go build -tags server -ldflags="-s -w" -o /app/server . && \
      cp -r frontend/dist /app/frontend_dist; \
    fi

# Runtime stage - minimal image
FROM gcr.io/distroless/static-debian12

# Copy the binary
COPY --from=builder /app/server /server

# Copy frontend assets
COPY --from=builder /app/frontend_dist /frontend/dist

# Expose the default port
EXPOSE 8080

# Bind to all interfaces (required for Docker)
# Can be overridden at runtime with -e WAILS_SERVER_HOST=...
ENV WAILS_SERVER_HOST=0.0.0.0

# Run the server
ENTRYPOINT ["/server"]
