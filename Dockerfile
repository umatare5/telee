# The context is assembled by goreleaser from the built binaries and the files named in
# dockers_v2.extra_files, so this copies rather than builds.
FROM scratch

# dockers_v2 lays the build context out as linux/<arch>/<binary>, so the binary is not at
# the context root.
ARG TARGETPLATFORM

# The trust store comes from a pinned image rather than the build context, which carries
# only the binary and the license notices.
COPY --from=alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY LICENSE NOTICE /
COPY $TARGETPLATFORM/telee /telee

USER 65534:65534

ENTRYPOINT ["/telee"]
CMD ["--help"]

LABEL org.opencontainers.image.title="telee"
LABEL org.opencontainers.image.description="CLI for executing commands on network devices"
LABEL org.opencontainers.image.vendor="umatare5"
LABEL org.opencontainers.image.source="https://github.com/umatare5/telee"
LABEL org.opencontainers.image.licenses="MIT"
