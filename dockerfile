#basically need to compile both worker and upload service and ffmpeg is there
FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/upload ./cmd/uploadService
RUN go build -o /out/worker ./cmd/workerService

#https://github.com/ggml-org/whisper.cpp
#https://github.com/ggml-org/whisper.cpp/releases/tag/v1.9.5

FROM alpine:3.20 AS whisperworker
WORKDIR /whisper
RUN apk add --no-cache git make g++ bash curl cmake ffmpeg ca-certificates libstdc++ libgomp
RUN git clone --depth 1 --branch v1.9.5 https://github.com/ggml-org/whisper.cpp.git .
RUN cmake -B build -DCMAKE_BUILD_TYPE=Release
RUN cmake --build build -j --config Release
RUN bash ./models/download-ggml-model.sh small

FROM alpine:3.20 AS upload
WORKDIR /app
COPY --from=builder /out/upload .
RUN chmod +x ./upload
EXPOSE 8080
CMD ["./upload"]

FROM alpine:3.20 AS worker
WORKDIR /app
RUN apk add --no-cache ffmpeg ca-certificates libstdc++
COPY --from=whisperworker /whisper/build /tmp/whisper-build
RUN find /tmp/whisper-build -name 'lib*.so*' -exec cp -a {} /usr/local/lib/ \; \
    && cp /tmp/whisper-build/bin/whisper-cli /usr/local/bin/whisper-cli \
    && chmod +x /usr/local/bin/whisper-cli \
    && rm -rf /tmp/whisper-build
RUN mkdir -p /app/models
COPY --from=whisperworker /whisper/models/ggml-small.bin /app/models/ggml-small.bin
COPY --from=builder /out/worker .
RUN chmod +x ./worker
ENV LD_LIBRARY_PATH=/usr/local/lib
ENV FFMPEG_BIN=ffmpeg
ENV WHISPER_BIN=whisper-cli
ENV MODEL_LOCATION=/app/models/ggml-small.bin
ENV FFMPEG_BIN=ffmpeg
CMD ["./worker"]