package uploadservice

import (
	"fmt"
	"net/http"
	"log"
	"video_processing_pipeline/internal/uploader"
	"video_processing_pipeline/internal/uploader/chunkersse"
)

func main() {
	manager:=chunkersse.NewChunkedUploadManager("uploads")
	handler:=uploader.NewHandlerStruct(manager)
	http.HandleFunc("/upload/init",handler.HandleStartUpload)
	http.HandleFunc("/upload/chunk",handler.HandleUploadChunks)
	http.HandleFunc("/upload/complete",handler.HandleCompleteUpload)
	http.HandleFunc("/upload/status",handler.HandleGetStatusUpload)
	http.HandleFunc("/",uploader.UploadHandler)
	http.HandleFunc("/upload",uploader.UploadStreamingHandler)
	// http.HandleFunc("/handleUpload",uploader.UploadChunkedHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}