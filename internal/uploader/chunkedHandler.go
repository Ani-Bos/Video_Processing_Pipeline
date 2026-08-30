package uploader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"video_processing_pipeline/internal/uploader/chunkersse"
)

type HandlerStruct struct{
	manager chunkersse.UploadManager
}

func NewHandlerStruct(mgr chunkersse.UploadManager)*HandlerStruct{
	return &HandlerStruct{
      manager: mgr,
	}
}

func(h *HandlerStruct)HandleStartUpload(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into starting upload session for chunking")
	var req chunkersse.RequestWrapper
    //validating request
	fmt.Println(req)
	fmt.Println(r.Body)
    err:=json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("Decode error: %v\n", err)
    fmt.Printf("Decoded req: %+v\n", req)
	if err!=nil{
       http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.TotalSize<=0 || req.ChunkSize<=0{
		http.Error(w, "Invalid size parameters", http.StatusBadRequest)
        return
	}
    uploadsession,err:=h.manager.InitiateUpload(&req)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//returning uploaded session details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsession)
}

func(h *HandlerStruct)HandleUploadChunks(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into uploading chunks")
	uploadID := r.URL.Query().Get("upload_id")
	chunknum:=0
	chunkstr := r.URL.Query().Get("chunk")
	chunknum, err := strconv.Atoi(chunkstr)
	if err!=nil{
		http.Error(w, "Invalid chunk number", http.StatusBadRequest)
        return
	}
	upload_ack,err:=h.manager.UploadChunk(uploadID,chunknum,r.Body)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(upload_ack)
}


func(h *HandlerStruct) HandleCompleteUpload(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into Handle complete Upload")
	uploadID := r.URL.Query().Get("upload_id")
	uploadResp,err:=h.manager.CompleteUpload(uploadID)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadResp)
}

func(h *HandlerStruct) HandleGetStatusUpload(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into getting status of uploading chunks")
	uploadID := r.URL.Query().Get("upload_id")
	uploadsts,err:=h.manager.GetUploadStatus(uploadID)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsts)
}