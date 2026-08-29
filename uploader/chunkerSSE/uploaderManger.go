package chunkersse

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (m *ChunkedUploadManager) InitiateUpload(req *RequestWrapper) (*ChunkedUpload, error) {
	total_no_of_chunks := int(req.TotalSize / req.ChunkSize)
	//like say 105 MB 10 MB so total is 10 chunks
	//5 left so that why at last increased 1 to process that data as well
	if (int)(req.TotalSize%req.ChunkSize) != 0 {
		total_no_of_chunks += 1
	}

	uploadId := generateUploadID()

	uploadSessionchunk := &ChunkedUpload{
		ID:             uploadId,
		FileName:       filepath.Base(req.FileName),
		TotalSize:      req.TotalSize,
		ChunkSize:      req.ChunkSize,
		TotalChunks:    int64(total_no_of_chunks),
		UploadedChunks: make(map[int]bool),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now().Add(24 * time.Hour),
	}
	//create temp directory path for this upload
	uploadpath := filepath.Join(m.uploadDir, "chunks", uploadId)
	err := os.Mkdir(uploadpath, 0755)
	if err != nil {
		return uploadSessionchunk, err
	}
	m.mu.Lock()
	m.uploads[uploadId] = uploadSessionchunk
	m.mu.Unlock()
	return uploadSessionchunk, nil
}

// handles upload of a single chunk
// chunks can be uploaded in any order and retried if they fail
func (m *ChunkedUploadManager) UploadChunk(uploadId string, chunkNumber int, data io.Reader) error {
	//retriveing upload session
	m.mu.Lock()
	uploaded_session, exists := m.uploads[uploadId]
	m.mu.Unlock()
	if !exists {
		return fmt.Errorf("uploaded session not fountd")
	}
	//validate chunk number
	if chunkNumber < 0 || chunkNumber >= int(uploaded_session.TotalChunks) {
		return fmt.Errorf("Invalid chunk number")
	}
	//check already chunk is uploaded
	m.mu.Lock()
	if uploaded_session.UploadedChunks[chunkNumber] == true {
		m.mu.Unlock()
		fmt.Println("chunk already uploaded")
		return nil
	}
	m.mu.TryLock()
	//creqate the new chunk file

	//create the chunk file name
	chunkfilePath := filepath.Join(m.uploadDir, "chunks", uploadId, fmt.Sprintf("%s.chunk.%d", chunkNumber))
	//create a new chunk file and buffer the data intoi it
	nextchunkfile, err := os.Create(chunkfilePath)
	if err != nil {
		return err
	}
	defer nextchunkfile.Close()
	//straming chunk data to file to handler
	//validate chunk size except for lasdt chunk
	expected_size := uploaded_session.ChunkSize
	if chunkNumber == int(uploaded_session.TotalChunks-1) {
		expected_size = uploaded_session.TotalSize - (int64(chunkNumber) * uploaded_session.ChunkSize)
	}

	m.mu.Lock()
	uploaded_session.UploadedChunks[chunkNumber] = true
	uploadedCount := len(uploaded_session.UploadedChunks)
	m.mu.Unlock()

	return nil
}

//completeupload assembles chunks into final file
// this should be called after all chunks have been uploaded

func (m *ChunkedUploadManager) CompleteUpload(UploadId string) (*ChunkedUpload, error) {
	m.mu.Lock()
	uploadedSesssion, exists := m.uploads[UploadId]
	m.mu.Unlock()
	if !exists {
		return nil, fmt.Errorf("uploaded session not fountd")
	}
	//verify all chunks are uploaded
	//total no of entries in map
	m.mu.Lock()
	uploadcount := len(uploadedSesssion.UploadedChunks)
    m.mu.Unlock()
	if uploadcount != int(uploadedSesssion.TotalChunks){
		return nil,fmt.Errorf("not all chunks are uploaded")
	}
	//create and assembling all chunks in final file chunks
	FinalchunkfilePath := filepath.Join(m.uploadDir, uploadedSesssion.FileName)
	//create a new chunk file and buffer the data intoi it
	Finalchunkfile, err := os.Create(FinalchunkfilePath)
	if err != nil {
		return nil,err
	}
	defer Finalchunkfile.Close()
	//assembling chunks in all order
	chunksdir:=filepath.Join(m.uploadDir,"chunks",UploadId)
	for i:=0;i<int(uploadedSesssion.TotalChunks);i++{
		chunkpath:=filepath.join(chunksdir,fmt.Sprintf("chunk_%d", i))

		chunkfile,err:=os.Open(chunkpath)
		if err!=nil{
			return nil,err
		}
	}
}
//return the current state oif chunked upload
//clients determine which chunks need to upload after a failure
func(m *ChunkedUploadManager)GetUploadStatus(UploadId string)(*ChunkedUpload,error){
   m.mu.Lock()
	uploadedSesssion, exists := m.uploads[UploadId]
	m.mu.Unlock()
	if !exists {
		return nil, fmt.Errorf("uploaded session not fountd")
	}
    missingChunks := make([]int,0)
	m.mu.Lock()
	for i:=0 ; i<int(uploadedSesssion.TotalChunks);i++{
		if uploadedSesssion.UploadedChunks[i]==false{
			missingChunks = append(missingChunks, i)
		}
	}
	uploadedcnt := len(uploadedSesssion.UploadedChunks)
	m.mu.Unlock()
	
}

func generateUploadID() string {
	data := make([]byte, 16)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Errorf("Failed to generate upload ID: %w", err))
	}
	return hex.EncodeToString(data)
}
