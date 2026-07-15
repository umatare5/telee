# Dockerfile for telee

FROM scratch

# Copy ca-certificates for HTTPS/SSH to network devices
COPY --from=alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY telee /telee

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/telee"]

# Default command shows help
CMD ["--help"]
