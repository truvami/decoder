FROM --platform=linux/amd64 golang:1.22.4 AS base

# Metadata as defined in OCI image spec annotations
LABEL org.opencontainers.image.vendor="truvami"
LABEL org.opencontainers.image.title="decoder"
LABEL org.opencontainers.image.authors="hey@truvami.com"

FROM alpine:3.19 AS certs
RUN apk --update add ca-certificates && rm -rf /var/cache/apk/*

# Build the application
FROM base AS builder

# ENV GOPROXY https://artifactory.swisscom.com/artifactory/api/go/proxy-golang-go-virtual

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/decoder ./main.go

# Create a minimal image
FROM --platform=linux/amd64 scratch AS runner

ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

COPY --from=base /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=base /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2

# Create a user
USER 1000

# Copy the binary from the builder image
COPY --from=builder /app/decoder /app/decoder

ENTRYPOINT ["/app/decoder"]