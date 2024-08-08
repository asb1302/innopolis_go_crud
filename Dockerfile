FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd

# Начинаем новый этап с нуля (multi-stage build)
FROM alpine:latest

RUN apk update && apk add --no-cache curl

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
