FROM golang:1.23 AS builder
WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/config /app/config

CMD ["/app/server"]