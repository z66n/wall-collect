# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o wally \
    wally.go

# Final stage with distroless
FROM gcr.io/distroless/static-debian11

# Set workdir
WORKDIR /app

# Copy binary and uploads dir with correct ownership for `nonroot` (UID 65532)
COPY --chown=65532:65532 --from=builder /app/wally .
COPY --chown=65532:65532 --from=builder /app/uploads ./uploads

# Set user to nonroot
USER nonroot:nonroot

# Expose port and set environment variables
EXPOSE 8080
ENV UPLOAD_DIR=/app/uploads \
    PORT=8080

# Run the app
CMD ["./wally", "-addr", "0.0.0.0:8080"]
