FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/echoid ./cmd/fingerprint

FROM alpine:latest

RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /app/echoid .

ENTRYPOINT ["/app/echoid"]
