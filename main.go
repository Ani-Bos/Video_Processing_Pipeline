package main
import (
	"fmt"
	"net/http"
	"log"
	"video_processing_pipeline/uploader"
)

func main() {
	http.HandleFunc("/",uploader.UploadHandler)
	http.HandleFunc("/upload",uploader.UploadStreamingHandler)
	http.HandleFunc("/handleUpload",uploader.UploadChunkedHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}