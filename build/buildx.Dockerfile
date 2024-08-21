# syntax=docker/dockerfile:1.4
FROM golang:1.23

# Metadata as defined in OCI image spec annotations
LABEL org.opencontainers.image.vendor="truvami"
LABEL org.opencontainers.image.title="decoder"
LABEL org.opencontainers.image.authors="hey@truvami.com"

ENV GOROOT /usr/local/go

# Allow to download a more recent version of Go.
# https://go.dev/doc/toolchain
# GOTOOLCHAIN=auto is shorthand for GOTOOLCHAIN=local+auto
ENV GOTOOLCHAIN auto

COPY decoder /usr/bin/
CMD ["decoder"]