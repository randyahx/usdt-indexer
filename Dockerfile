# Build stage
FROM golang:1.22 AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Disable CGO and build statically
RUN CGO_ENABLED=0 GOOS=linux go build -o usdt-indexer .

# Use a minimal base image for the runtime
FROM gcr.io/distroless/static

# Set the working directory
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/usdt-indexer /app/

COPY config.json /app/config.json

# Expose the port (if needed)
EXPOSE 2112

# Run the Go binary
CMD ["/app/usdt-indexer", "--config-path", "/app/config.json"]
