# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wallplayer ./cmd

# Runtime stage
FROM alpine:latest

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Create app user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/wallplayer .

# Create directories for videos and data
RUN mkdir -p /app/videos /app/data/thumbnails /app/data/subtitles

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 9999

# Set environment variables
ENV PORT=9999
ENV VIDEOS_DIR=/videos

# Run the binary
CMD ["./wallplayer"]