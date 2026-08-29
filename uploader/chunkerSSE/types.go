package chunkersse

import (
	"io"
	"sync"
	"time"
)

//chunked upload represent tjhe state of in progrees chunk
type ChunkedUpload struct {
	ID             string
	FileName       string
	TotalSize      int64
	ChunkSize      int64
	TotalChunks    int64
	UploadedChunks map[int]bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

//upload manager to handle the lifecycle of chunked uplopads

type ChunkedUploadManager struct{
	uploads map[string]*ChunkedUpload
	uploadDir string
	mu sync.RWMutex
}

//New chunked manager created a new manager for chunked uoloads

func NewChunkedUploadManager(dir string) * ChunkedUploadManager{
	return &ChunkedUploadManager{
		uploads: make(map[string]*ChunkedUpload),
		uploadDir: dir,
	}
}

type  RequestWrapper struct{
	FileName string
	TotalSize int64
	ChunkSize int64
}

type UploadStatus struct{
	UploadId string
	FileName string
	UploadedChunks int
    TotalChunks int64
	MissingChunks []int
	IsComplete bool

}

type UploadResponse struct{
	UploadId string
	FileName string
	size int64
	FilePath string
}

type ChunkAcknowledgemnt struct{
	ChunkNumber int
	UploadedChunks int
	TotalChunks int64
}

type UploadManager interface{
	InitiateUpload(req *RequestWrapper)(*ChunkedUpload,error)
	UploadChunk(uploadId string, chunkNumber int, data io.Reader)(*ChunkAcknowledgemnt,error)
	CompleteUpload(UploadId string)(*UploadResponse,error)
	GetUploadStatus(UploadId string)(*UploadStatus,error)
}

