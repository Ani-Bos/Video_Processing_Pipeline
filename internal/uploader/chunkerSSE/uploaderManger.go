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
	err := os.MkdirAll(uploadpath, 0755)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.uploads[uploadId] = uploadSessionchunk
	m.mu.Unlock()
	return uploadSessionchunk, nil
}

// handles upload of a single chunk
// chunks can be uploaded in any order and retried if they fail
func (m *ChunkedUploadManager) UploadChunk(uploadId string, chunkNumber int, data io.Reader) (*ChunkAcknowledgemnt,error) {
	//retriveing upload session
	m.mu.RLock()
	uploaded_session, exists := m.uploads[uploadId]
	m.mu.RUnlock()
	if !exists {
		return nil,fmt.Errorf("uploaded session not fountd")
	}
	//validate chunk number
	if chunkNumber < 0 || chunkNumber >= int(uploaded_session.TotalChunks) {
		return nil,fmt.Errorf("Invalid chunk number")
	}
	//check already chunk is uploaded
	m.mu.RLock()
	if uploaded_session.UploadedChunks[chunkNumber] == true {
		m.mu.RUnlock()
		fmt.Println("chunk already uploaded")
		return &ChunkAcknowledgemnt{
			  ChunkNumber: chunkNumber,
              UploadedChunks: len(uploaded_session.UploadedChunks),
              TotalChunks: uploaded_session.TotalChunks,
       },nil
	}
	m.mu.RUnlock()
	//creqate the new chunk file

	//create the chunk file name
	chunkfilePath := filepath.Join(m.uploadDir, "chunks", uploadId, fmt.Sprintf("chunk_%d", chunkNumber))
	//create a new chunk file and buffer the data intoi it
	nextchunkfile, err := os.Create(chunkfilePath)
	if err != nil {
		return nil,err
	}
	defer nextchunkfile.Close()
	//straming chunk data to file to handle
	written,err:= io.Copy(nextchunkfile,data)
	if err!=nil{
		os.Remove(chunkfilePath)
		return nil,err
	}
	//validate chunk size except for lasdt chunk
	expected_size := uploaded_session.ChunkSize
	if chunkNumber == int(uploaded_session.TotalChunks-1) {
		expected_size = uploaded_session.TotalSize - (int64(chunkNumber) * uploaded_session.ChunkSize)
	}
    if written!=expected_size{
		os.Remove(chunkfilePath)
		return nil, fmt.Errorf("chunk size mismatch: got %d, expected %d", written, expected_size)
	}
	m.mu.Lock()
	uploaded_session.UploadedChunks[chunkNumber] = true
	uploadedCount := len(uploaded_session.UploadedChunks)
	m.mu.Unlock()

	return &ChunkAcknowledgemnt{
              ChunkNumber: chunkNumber,
              UploadedChunks: uploadedCount,
              TotalChunks: uploaded_session.TotalChunks,
	},nil
}

//completeupload assembles chunks into final file
// this should be called after all chunks have been uploaded

func (m *ChunkedUploadManager) CompleteUpload(UploadId string)(*UploadResponse,error){
	m.mu.RLock()
	uploadedSesssion, exists := m.uploads[UploadId]
	m.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("uploaded session not fountd")
	}
	//verify all chunks are uploaded
	//total no of entries in map
	m.mu.RLock()
	uploadcount := len(uploadedSesssion.UploadedChunks)
    m.mu.RUnlock()
	if uploadcount != int(uploadedSesssion.TotalChunks){
		return nil,fmt.Errorf("incomplete upload: %d/%d chunks", uploadcount, uploadedSesssion.TotalChunks)
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
		chunkpath:=filepath.Join(chunksdir,fmt.Sprintf("chunk_%d", i))
		chunkfile,err:=os.Open(chunkpath)
		if err!=nil{
			return nil,err
		}
		_,err=io.Copy(Finalchunkfile,chunkfile)
		if err!=nil{
			chunkfile.Close()
			return nil,err
		}
		chunkfile.Close()
	}
	//cleanup upload session and remove all chunks
	os.RemoveAll(chunksdir)
	// m.mu.Lock()
	// delete(m.uploads,UploadId)
	// m.mu.Unlock()
	return &UploadResponse{
		UploadId:UploadId,
		FileName:uploadedSesssion.FileName,
		Size:uploadedSesssion.TotalSize,
		FilePath:FinalchunkfilePath,
	},nil
  //	FileName string
	// RawPath string
	// Status string
	// Metadata datatypes.JSON
	//all update in db as well
	//also upload to worker queue i.e redis queue
}
//return the current state oif chunked upload
//clients determine which chunks need to upload after a failure
func(m *ChunkedUploadManager)GetUploadStatus(UploadId string)(*UploadStatus,error){
   m.mu.RLock()
	uploadedSesssion, exists := m.uploads[UploadId]
	m.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("uploaded session not fountd")
	}
	//finding missing chunks
    missingChunks := make([]int,0)
	m.mu.RLock()
	for i:=0 ; i<int(uploadedSesssion.TotalChunks);i++{
		if uploadedSesssion.UploadedChunks[i]==false{
			missingChunks = append(missingChunks, i)
		}
	}
	uploadedcnt := len(uploadedSesssion.UploadedChunks)
	m.mu.RUnlock()
    return &UploadStatus{
		UploadId:UploadId,
		FileName:uploadedSesssion.FileName,
		UploadedChunks:uploadedcnt,
		TotalChunks:uploadedSesssion.TotalChunks,
		MissingChunks:missingChunks,
		IsComplete:len(missingChunks)==0,
	},nil
	
}

func generateUploadID() string {
	data := make([]byte, 16)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Errorf("Failed to generate upload ID: %w", err))
	}
	return hex.EncodeToString(data)
}
