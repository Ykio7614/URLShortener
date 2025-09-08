FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o shortener ./cmd/shortener

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/shortener .
COPY ./db/migrations ./db/migrations

CMD ["./shortener"]