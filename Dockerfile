FROM golang:1.13-alpine3.11 AS builder

WORKDIR /go/src/github.com/mjpitz/highlander-proxy

COPY bin bin

RUN OS=$(go env GOOS) && \
    ARCH=$(go env GOARCH) && \
    cp bin/highlander-proxy_${OS}_${ARCH} highlander-proxy

RUN echo "highlander:x:65534:65534:highlander:/:" > /etc_passwd

FROM alpine:3.11

COPY --from=builder /go/src/github.com/mjpitz/highlander-proxy/highlander-proxy /usr/bin/highlander-proxy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc_passwd /etc/passwd

USER highlander

ENTRYPOINT [ "/usr/bin/highlander-proxy" ]
