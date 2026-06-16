# Stage 1: Build binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/echoid ./cmd/fingerprint

# Stage 2: Minimal runtime
FROM alpine:latest

RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /app/echoid .

# NOTE: yt-dlp is not included; the youtube subcommand will not work
# inside the container unless you install yt-dlp separately.

ENTRYPOINT ["/app/echoid"]
