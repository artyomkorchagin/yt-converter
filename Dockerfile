FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/yt-converter ./cmd/main.go

FROM alpine:3.18
WORKDIR /app

RUN apk add --no-cache ffmpeg

COPY --from=builder /app/yt-converter /app/yt-converter
COPY ./config.env /app/config.env
COPY ./privkey.pem /app/privkey.pem
COPY ./cert.pem /app/cert.pem
COPY ./web /app/web
EXPOSE 3000
CMD ["/app/yt-converter"]