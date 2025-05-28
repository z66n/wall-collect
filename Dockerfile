# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# First copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Create empty uploads directory for build
RUN mkdir -p /app/uploads

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o wally \
    wally.go

# Final stage
FROM gcr.io/distroless/static-debian11
WORKDIR /app

# Copy binary
COPY --from=builder /app/wally .

# Create uploads directory (empty by default)
RUN mkdir -p /app/uploads

# Set permissions
USER nonroot:nonroot

EXPOSE 8080
ENV UPLOAD_DIR=/app/uploads \
    PORT=8080

CMD ["./wally", "-addr", ":$PORT"]