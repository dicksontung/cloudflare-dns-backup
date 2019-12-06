FROM golang:alpine as builder
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

## Minimize
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY cloudflare-dns-backup /cloudflare-dns-backup
WORKDIR /

CMD  ["./cloudflare-dns-backup"]