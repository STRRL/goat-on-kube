FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goat-exporter ./cmd/goat-exporter

FROM debian:bookworm-slim

LABEL org.opencontainers.image.title="goat-exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for GOAT metrics from Ethereum nodes"
LABEL org.opencontainers.image.source="https://github.com/strrl/goat-on-kube"

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/goat-exporter /usr/local/bin/goat-exporter

EXPOSE 9090

ENTRYPOINT ["/usr/local/bin/goat-exporter"]
