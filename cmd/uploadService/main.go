package uploadservice

import (
	"fmt"
	"net/http"
	"log"
	"video_processing_pipeline/internal/handler"
	"video_processing_pipeline/internal/uploader/chunkersse"
)

func main() {
	manager:=chunkersse.NewChunkedUploadManager("uploads")
	handler1:=handler.NewHandlerStruct(manager)
	http.HandleFunc("/upload/init",handler1.HandleStartUpload)
	http.HandleFunc("/upload/chunk",handler1.HandleUploadChunks)
	http.HandleFunc("/upload/complete",handler1.HandleCompleteUpload)
	http.HandleFunc("/upload/status",handler1.HandleGetStatusUpload)
	http.HandleFunc("/",handler.UploadHandler)
	http.HandleFunc("/upload",handler.UploadStreamingHandler)
	// http.HandleFunc("/handleUpload",uploader.UploadChunkedHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}