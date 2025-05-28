# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./


# Build the binary (remove inline comments and fix line continuations)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o wally \
    wally.go

RUN mkdir -p /app/uploads && chown -R 65532:65532 /app/uploads

# Final stage
FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /app/wally .
COPY --from=builder /app/uploads ./uploads
USER nonroot:nonroot
EXPOSE 8080
ENV UPLOAD_DIR=/app/uploads \
    PORT=8080
CMD ["./wally", "-addr", "0.0.0.0:8080"]