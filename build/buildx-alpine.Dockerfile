# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine

# Metadata as defined in OCI image spec annotations
LABEL org.opencontainers.image.vendor="truvami"
LABEL org.opencontainers.image.title="decoder"
LABEL org.opencontainers.image.authors="hey@truvami.com"

ENV GOROOT /usr/local/go

# Allow to download a more recent version of Go.
# https://go.dev/doc/toolchain
# GOTOOLCHAIN=auto is shorthand for GOTOOLCHAIN=local+auto
ENV GOTOOLCHAIN auto

# gcc is required to support cgo;
# git and mercurial are needed most times for go get`, etc.
# See https://github.com/docker-library/golang/issues/80
RUN apk --no-cache add gcc musl-dev git mercurial

# Set all directories as safe
RUN git config --global --add safe.directory '*'

COPY decoder /usr/bin/
CMD ["decoder"]