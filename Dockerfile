FROM golang:1.25.4 AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080
CMD ["air"]



FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .



FROM debian:bookworm-slim as prod

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates

RUN mkdir -p /app/uploads
RUN mkdir -p /app/stockage/livraisons
RUN mkdir -p /app/stockage/plannings

EXPOSE 8080
CMD ["./server"]