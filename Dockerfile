FROM golang:1.21 AS builder

WORKDIR /app

COPY ./ ./
RUN make clean && make build

FROM alpine AS cacerts
RUN apk update && apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /app/dopark /dopark
COPY --from=cacerts /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt 
ENTRYPOINT ["/dopark"]
