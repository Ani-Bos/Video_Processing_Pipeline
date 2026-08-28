package uploader

import (
	"fmt"
	"net/http"
)

func UploadChunkedHandler(w http.ResponseWriter, r *http.Request){
  fmt.Println("Entering into chunked uploads")
}