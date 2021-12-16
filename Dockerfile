FROM golang:1.17.5-alpine3.15 as builder
RUN mkdir /src
ADD src /src
WORKDIR /src
RUN go build -o portscan-prometheus-exporter
FROM alpine:3.15
COPY --from=builder /src/portscan-prometheus-exporter /app/portscan-prometheus-exporter
WORKDIR /app
ENTRYPOINT ["/app/portscan-prometheus-exporter"]


