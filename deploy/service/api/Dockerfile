FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -ldflags="-s -w" -o app ./cmd/api

# --------
FROM alpine:3.20

COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt

RUN addgroup -S appgroup && adduser -S appuser -G appgroup -u 10001

COPY --from=builder --chown=appuser:appgroup /src/app .

USER appuser

ENTRYPOINT ["/app"]