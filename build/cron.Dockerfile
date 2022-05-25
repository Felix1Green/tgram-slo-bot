FROM golang:1.17-alpine
WORKDIR /service
COPY .. .
RUN go build -o cron ./cmd/slo_tracker/main.go
ENTRYPOINT ["./cron"]