# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \  # Strip debug symbols
    -o wally \         # Output name matches your file
    wally.go           # Builds your main wally file

# Final stage - ultra slim image
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy only the built binary
COPY --from=builder /app/wally .
COPY --from=builder /app/uploads ./uploads

# Non-root user for security
USER nonroot:nonroot

# Environment variables (customize as needed)
ENV UPLOAD_DIR=/app/uploads \
    PORT=8080

EXPOSE 8080

# Entrypoint matches your binary name
ENTRYPOINT ["./wally", "-addr", ":$PORT"]