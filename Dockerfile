# Stage 1: Build the Go binary
FROM golang:1.24.2-alpine AS builder

# Install build dependencies (make, git, and bash for Makefile)
RUN apk add --no-cache make git bash

WORKDIR /build

# Copy all source files (including Makefile)
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Build using Makefile
RUN make build VERSION=${VERSION} COMMIT=${COMMIT} BUILD_DATE=${BUILD_DATE}

# Stage 2: Runtime image (distroless for security)
FROM gcr.io/distroless/static:nonroot

# Copy the binary from builder
COPY --from=builder /build/bin/mcp-server /bin/mcp-server

# Distroless runs as nonroot user (uid=65532) by default
USER nonroot:nonroot

# Expose HTTP port
EXPOSE 8080

# Set entrypoint and default command
# The MCP server will run in HTTP mode, suitable for Kubernetes
ENTRYPOINT ["/bin/mcp-server"]
CMD ["serve", "--transport", "http", "--host", "0.0.0.0", "--port", "8080"]
