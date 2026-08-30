#!/bin/bash

VIDEO="c:\\Users\\Aniket\\Documents\\1080_30_8.00_Jun222021(1).mp4"
CHUNK_SIZE=5242880

split -d -b "$CHUNK_SIZE" "$VIDEO" /tmp/chunk_

resp=$(curl -s -X POST http://localhost:8080/upload/init \
  -H "Content-Type: application/json" \
  -d '{"filename":"c:\\Users\\Aniket\\Documents\\1080_30_8.00_Jun222021(1).mp4","total_size":28857413,"chunk_size":5242880}')

echo "Raw response: $resp"

ID=$(echo "$resp" | grep -o '"ID":"[^"]*' | grep -o '[^"]*$')
echo "Upload ID: $ID"

if [ -z "$ID" ]; then
    echo "ERROR: Failed to get upload ID."
    exit 1
fi

for i in {0..5}; do
    CHUNK_FILE=$(printf "/tmp/chunk_%02d" $i)
    if [ -f "$CHUNK_FILE" ]; then
        echo "Uploading chunk $i..."
        curl -s -X POST "http://localhost:8080/upload/chunk?upload_id=$ID&chunk=$i" \
          -H "Content-Type: application/octet-stream" \
          --data-binary @"$CHUNK_FILE"
    else
        echo "Missing chunk file: $CHUNK_FILE"
    fi
done

curl -X POST "http://localhost:8080/upload/complete?upload_id=$ID"

# rm -f /tmp/chunk_*