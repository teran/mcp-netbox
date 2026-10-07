# Uses a pre-built mcp-netbox binary from goreleaser.
# Minimal nonroot scratch image; TLS CA bundle and /etc/passwd for the
# unprivileged user (UID 65534) are copied in from the alpine builder stage.
# Usage:
#   goreleaser build --snapshot --clean
#   cp dist/mcp-netbox_linux_amd64_v1/mcp-netbox mcp-netbox-linux-amd64
#   cp dist/mcp-netbox_linux_arm64_v8.0/mcp-netbox mcp-netbox-linux-arm64
#   docker buildx build --platform linux/amd64,linux/arm64 -t image:tag .

FROM alpine:3.24 AS base
RUN apk add --no-cache ca-certificates && \
    echo 'nobody:x:65534:65534:nobody:/:/sbin/nologin' > /etc/passwd-minimal

FROM scratch
ARG TARGETARCH
COPY --from=base /etc/passwd-minimal /etc/passwd
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY mcp-netbox-linux-${TARGETARCH} /mcp-netbox
USER 65534:65534
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/mcp-netbox", "-mode", "http"]
LABEL org.opencontainers.image.source="https://github.com/teran/mcp-netbox"
LABEL org.opencontainers.image.description="Remote MCP server for NetBox"
LABEL org.opencontainers.image.licenses="Apache-2.0"
