# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/api .

# Set default environment variables
ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

# Run the executable
CMD ["./api"]
