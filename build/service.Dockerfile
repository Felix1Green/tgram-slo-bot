FROM golang:1.17-alpine
WORKDIR /service
COPY .. .
RUN go build -o service ./cmd/service/main.go
ENTRYPOINT ["./service"]