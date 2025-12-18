# Build stage
FROM golang:1-alpine AS builder

# Set the working directory
WORKDIR /tmp/build

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -trimpath -o telee ./cmd/main.go

# Final stage
FROM alpine:3.23.2
WORKDIR /app
COPY --from=builder /tmp/build/telee /bin/

ENTRYPOINT [ "/bin/telee" ]
