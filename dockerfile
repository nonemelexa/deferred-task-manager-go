# Stage 1: build
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./app

# Stage 2: run
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]