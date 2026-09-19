#basically need to compile both worker and upload service and ffmpeg is there
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/upload ./cmd/uploadService
RUN go build -o /out/worker ./cmd/workerService

FROM alpine:3.20 AS upload
WORKDIR /app
COPY --from=builder /out/upload .
EXPOSE 8080
CMD ["./upload"]

FROM alpine:3.20 AS worker
WORKDIR /app
RUN apk add --no-cache ffmpeg ca-certificates
COPY --from=builder /out/worker .
ENV FFMPEG_BIN=ffmpeg
CMD ["./worker"]