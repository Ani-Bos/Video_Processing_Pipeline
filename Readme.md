
Raw response: 
```json
{
    "ID":"84e90d3e1e0a7bea4905e4d029a8b949",
    "FileName":"1080_30_8.00_Jun222021(1).mp4",
    "TotalSize":28857413,"ChunkSize":5242880,
    "TotalChunks":6,
    "UploadedChunks":{},
    "CreatedAt":"2026-08-30T15:54:59.370928+05:30",
    "UpdatedAt":"2026-08-31T15:54:59.370928+05:30"
}
Upload ID: 84e90d3e1e0a7bea4905e4d029a8b949
Uploading chunk 0...
{
    "ChunkNumber":0,
    "UploadedChunks":1,
    "TotalChunks":6
}
Uploading chunk 1...
{
    "ChunkNumber":1,
    "UploadedChunks":2,
    "TotalChunks":6
}
Uploading chunk 2...
{
    "ChunkNumber":2,
    "UploadedChunks":3,
    "TotalChunks":6
}
Uploading chunk 3...
{
    "ChunkNumber":3,
    "UploadedChunks":4,
    "TotalChunks":6
}
Uploading chunk 4...
{
    "ChunkNumber":4,
    "UploadedChunks":5,
    "TotalChunks":6
}
Uploading chunk 5...
{
    "ChunkNumber":5,
    "UploadedChunks":6,
    "TotalChunks":6
}
{
    "UploadId":"84e90d3e1e0a7bea4905e4d029a8b949",
    "FileName":"1080_30_8.00_Jun222021(1).mp4",
    "Size":28857413,
    "FilePath":"uploads\\1080_30_8.00_Jun222021(1).mp4"
}
```

http://localhost:8080/upload/status?upload_id=34899b39c4948d91fef7ab047c57d8c7
```json
{
    "UploadId": "34899b39c4948d91fef7ab047c57d8c7",
    "FileName": "1080_30_8.00_Jun222021(1).mp4",
    "UploadedChunks": 6,
    "TotalChunks": 6,
    "MissingChunks": [],
    "IsComplete": true
}
```

```json
upload-1    | Entering into starting upload session for chunking
upload-1    | { 0 0}
upload-1    | &{0x14a211626048 <nil> <nil> false true {{} {0 0}} false false false 0x6dbc20}
upload-1    | Decode error: <nil>
upload-1    | Decoded req: {FileName:c:\Users\Aniket\Documents\1080_30_8.00_Jun222021(1).mp4 TotalSize:28857413 ChunkSize:5242880}
upload-1    | upload session is &{355f881081fb5914a218f898f6177a9c c:\Users\Aniket\Documents\1080_30_8.00_Jun222021(1).mp4 28857413 5242880 6 map[] 2026-09-19 20:33:42.921684171 +0000 UTC 2026-09-20 20:33:42.921684171 +0000 UTC}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  0
upload-1    | upload ack is &{0 1 6}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  1
upload-1    | upload ack is &{1 2 6}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  2
upload-1    | upload ack is &{2 3 6}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  3
upload-1    | upload ack is &{3 4 6}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  4
upload-1    | upload ack is &{4 5 6}
upload-1    | Entering into uploading chunks
upload-1    | upload Id in fxn handle upload chunks is  355f881081fb5914a218f898f6177a9c
upload-1    | chunk number is  5
upload-1    | upload ack is &{5 6 6}
upload-1    | Entering into Handle complete Upload
upload-1    | Enterinto into service layer for inserting into jobs data
upload-1    | Entering into publish events
upload-1    | task queue created &{transcode [123 34 86 105 100 101 111 73 100 34 58 34 51 53 53 102 56 56 49 48 56 49 102 98 53 57 49 52 97 50 49 56 102 56 57 56 102 54 49 55 55 97 57 99 34 44 34 82 97 119 80 97 116 104 34 58 34 47 100 97 116 97 47 117 112 108 111 97 100 115 47 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 70 105 108 101 78 97 109 101 34 58 34 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 79 117 116 80 117 116 68 105 114 34 58 34 34 125] map[] [3 critical] <nil>}
worker-1    | Entering into transcioding handler
worker-1    | Entering into transcoding using ffmpeg and hls format
worker-1    | Entering into publish events
worker-1    | task queue created &{thumbnail [123 34 86 105 100 101 111 73 100 34 58 34 51 53 53 102 56 56 49 48 56 49 102 98 53 57 49 52 97 50 49 56 102 56 57 56 102 54 49 55 55 97 57 99 34 44 34 82 97 119 80 97 116 104 34 58 34 47 100 97 116 97 47 117 112 108 111 97 100 115 47 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 70 105 108 101 78 97 109 101 34 58 34 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 79 117 116 80 117 116 68 105 114 34 58 34 47 100 97 116 97 47 116 114 97 110 115 99 111 100 101 47 51 53 53 102 56 56 49 48 56 49 102 98 53 57 49 52 97 50 49 56 102 56 57 56 102 54 49 55 55 97 57 99 34 125] map[] [3 critical] <nil>}
worker-1    | Enter into thubnail Handler
worker-1    | Entering into generating thubnail from FFMPEG output raw hls segments
worker-1    | ffmpeg version 6.1.1 Copyright (c) 2000-2023 the FFmpeg developers
worker-1    |   built with gcc 13.2.1 (Alpine 13.2.1_git20240309) 20240309
worker-1    |   configuration: --prefix=/usr --disable-librtmp --disable-lzma --disable-static --disable-stripping --enable-avfilter --enable-gpl --enable-ladspa --enable-libaom --enable-libass --enable-libbluray --enable-libdav1d --enable-libdrm --enable-libfontconfig --enable-libfreetype --enable-libfribidi --enable-libharfbuzz --enable-libmp3lame --enable-libopenmpt --enable-libopus --enable-libplacebo --enable-libpulse --enable-librav1e --enable-librist --enable-libsoxr --enable-libsrt --enable-libssh --enable-libtheora --enable-libv4l2 --enable-libvidstab --enable-libvorbis --enable-libvpx --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxcb --enable-libxml2 --enable-libxvid --enable-libzimg --enable-libzmq --enable-lto=auto --enable-lv2 --enable-openssl --enable-pic --enable-postproc --enable-pthreads --enable-shared --enable-vaapi --enable-vdpau --enable-version3 --enable-vulkan --optflags=-O3 --enable-libjxl --enable-libsvtav1 --enable-libvpl
worker-1    |   libavutil      58. 29.100 / 58. 29.100
worker-1    |   libavcodec     60. 31.102 / 60. 31.102
worker-1    |   libavformat    60. 16.100 / 60. 16.100
worker-1    |   libavdevice    60.  3.100 / 60.  3.100
worker-1    |   libavfilter     9. 12.100 /  9. 12.100
worker-1    |   libswscale      7.  5.100 /  7.  5.100
worker-1    |   libswresample   4. 12.100 /  4. 12.100
worker-1    |   libpostproc    57.  3.100 / 57.  3.100
worker-1    | [hls @ 0x70d9f84ff900] Skip ('#EXT-X-VERSION:3')
worker-1    | [hls @ 0x70d9f84ff900] Opening '/data/transcode/355f881081fb5914a218f898f6177a9c/index0.ts' for reading
worker-1    | Input #0, hls, from '/data/transcode/355f881081fb5914a218f898f6177a9c/index.m3u8':
worker-1    |   Duration: 00:00:25.20, start: 1.443444, bitrate: 0 kb/s
worker-1    |   Program 0 
worker-1    |     Metadata:
worker-1    |       variant_bitrate : 0
worker-1    |   Stream #0:0: Video: h264 (High) ([27][0][0][0] / 0x001B), yuv420p(tv, smpte170m/bt470bg/smpte170m), 1080x1080, 30 fps, 30 tbr, 90k tbn
worker-1    |     Metadata:
worker-1    |       variant_bitrate : 0
worker-1    |   Stream #0:1(eng): Audio: aac (LC) ([15][0][0][0] / 0x000F), 44100 Hz, stereo, fltp
worker-1    |     Metadata:
worker-1    |       variant_bitrate : 0
worker-1    | Stream mapping:
worker-1    |   Stream #0:0 -> #0:0 (h264 (native) -> mjpeg (native))
worker-1    | Press [q] to stop, [?] for help
worker-1    | [hls @ 0x70d9f84ff900] Opening '/data/transcode/355f881081fb5914a218f898f6177a9c/index0.ts' for reading
worker-1    | [swscaler @ 0x70d9f1dec040] deprecated pixel format used, make sure you did set range correctly
worker-1    |     Last message repeated 3 times
worker-1    | Output #0, image2, to '/data/transcode/355f881081fb5914a218f898f6177a9c/ani.jpg':
worker-1    |   Metadata:
worker-1    |     encoder         : Lavf60.16.100
worker-1    |   Stream #0:0: Video: mjpeg, yuvj420p(pc, smpte170m/bt470bg/smpte170m, progressive), 1080x1080, q=2-31, 200 kb/s, 30 fps, 30 tbn
worker-1    |     Metadata:
worker-1    |       variant_bitrate : 0
worker-1    |       encoder         : Lavc60.31.102 mjpeg
worker-1    |     Side data:
worker-1    |       cpb: bitrate max/min/avg: 0/0/200000 buffer size: 0 vbv_delay: N/A
[image2 @ 0x70d9f7a22340] The specified filename '/data/transcode/355f881081fb5914a218f898f6177a9c/ani.jpg' does not contain an image sequence pattern or a pattern is invalid.
worker-1    | [image2 @ 0x70d9f7a22340] Use a pattern such as %03d for an image sequence or use the -update option (with -frames:v 1 if needed) to write a single image.
worker-1    | [out#0/image2 @ 0x70d9f61d12c0] video:7kB audio:0kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: unknown
worker-1    | frame=    1 fps=0.0 q=3.0 Lsize=N/A time=00:00:00.00 bitrate=N/A speed=   0x    
worker-1    | Entering into publish events
worker-1    | task queue created &{notify [123 34 86 105 100 101 111 73 100 34 58 34 51 53 53 102 56 56 49 48 56 49 102 98 53 57 49 52 97 50 49 56 102 56 57 56 102 54 49 55 55 97 57 99 34 44 34 82 97 119 80 97 116 104 34 58 34 47 100 97 116 97 47 117 112 108 111 97 100 115 47 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 70 105 108 101 78 97 109 101 34 58 34 99 58 92 92 85 115 101 114 115 92 92 65 110 105 107 101 116 92 92 68 111 99 117 109 101 110 116 115 92 92 49 48 56 48 95 51 48 95 56 46 48 48 95 74 117 110 50 50 50 48 50 49 40 49 41 46 109 112 52 34 44 34 79 117 116 80 117 116 68 105 114 34 58 34 47 100 97 116 97 47 116 114 97 110 115 99 111 100 101 47 51 53 53 102 56 56 49 48 56 49 102 98 53 57 49 52 97 50 49 56 102 56 57 56 102 54 49 55 55 97 57 99 34 125] map[] [3 critical] <nil>}
worker-1    | asynq: pid=1 2026/09/19 20:35:49.268275 WARN: Retry exhausted for task id=609fa42c-d674-4cba-bb3e-ebce980a0805

```

docker compose build --no-cache
docker compose up

docker exec -it video_processing_pipeline-postgres-1 psql -U postgres -d video_db

docker cp video_processing_pipeline-worker-1:/data/transcode/355f881081fb5914a218f898f6177a9c/ani.jpg C:\Users\Aniket\Desktop\LLD\Video_Processing_Pipeline\internal\uploader\tmp\ani.jpg

```