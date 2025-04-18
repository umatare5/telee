# Build stage
FROM golang:1-alpine AS builder

WORKDIR /app
COPY telee /bin/

ENTRYPOINT [ "/bin/telee" ]
