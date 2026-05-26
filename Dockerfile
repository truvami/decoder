# syntax=docker/dockerfile:1

###############################################################################
# certs stage
#
# This Dockerfile is built by GoReleaser's `dockers_v2`, which supplies the
# prebuilt binary in the build context (per target platform). We intentionally
# do NOT compile here: GoReleaser builds the binaries on CI.
# See https://goreleaser.com/customization/package/dockers_v2/
#
# This Alpine stage exists only to source the handful of files a `scratch` image
# lacks: CA certificates, tzdata, and a passwd/group entry for the non-root user.
###############################################################################
FROM alpine:3.21 AS certs

RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates && \
    addgroup -g 65532 -S nonroot && \
    adduser -u 65532 -S -G nonroot -H -h /nonexistent -s /sbin/nologin nonroot

###############################################################################
# final stage — scratch (truly empty, smallest possible runtime)
###############################################################################
FROM scratch

# Metadata as defined in OCI image spec annotations
LABEL org.opencontainers.image.vendor="truvami"
LABEL org.opencontainers.image.title="decoder"
LABEL org.opencontainers.image.authors="hey@truvami.com"

# Files a networked, time-aware, non-root CLI needs (scratch ships none):
#   - CA certificates: outbound TLS (AWS IoT Wireless, LoRa Cloud, self-update, ...)
#   - tzdata: timezone-aware time handling
#   - passwd/group: so USER nonroot resolves to a real uid/gid
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=certs /etc/passwd /etc/passwd
COPY --from=certs /etc/group /etc/group

# Prebuilt binary supplied by GoReleaser dockers_v2 (CGO_ENABLED=0, statically
# linked). The context places artifacts under a per-platform path (e.g.
# linux/amd64/) so the same Dockerfile produces every architecture.
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/decoder /usr/bin/decoder

USER nonroot:nonroot

# HTTP server (`decoder http`), also serves Prometheus metrics on /metrics.
EXPOSE 8080

# Absolute path: scratch has no PATH, which is unreliable under containerd/k8s.
ENTRYPOINT ["/usr/bin/decoder"]
