FROM golang:1.19-alpine as builder
WORKDIR /builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o cloudflare-ddns cmd/cloudflare-ddns/main.go

FROM redhat/ubi8-minimal
WORKDIR /opt/app
COPY --from=builder /builder/cloudflare-ddns .
RUN chgrp -R 0 /opt/app && chmod -R g=u /opt/app
USER 1000

ENTRYPOINT ["/opt/app/cloudflare-ddns"]
