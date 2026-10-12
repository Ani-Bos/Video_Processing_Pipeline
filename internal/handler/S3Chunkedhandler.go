package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"video_processing_pipeline/internal/model"
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"
	"video_processing_pipeline/internal/service"
	"video_processing_pipeline/internal/uploader/chunkersse"
)

type HandlerStructForS3 struct {
	manager chunkersse.MiniIOManager
	srvc service.InterfaceInjectRepoJob
	repo repository.ChunkRepo
	asynqsrvc queue.AsyncPublishersrvc
}

func NewHandlerStructForS3(mgr chunkersse.MiniIOManager, srvce service.InterfaceInjectRepoJob, repostry repository.ChunkRepo, asysrvc queue.AsyncPublishersrvc)*HandlerStructForS3{
	return &HandlerStructForS3{
       manager: mgr,
	   srvc: srvce,
	   repo: repostry,
	   asynqsrvc: asysrvc,
	}
}

func(h *HandlerStructForS3)HandleStartUpload(w http.ResponseWriter, r *http.Request){
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
	fmt.Println("upload session is", uploadsession)
	//returning uploaded session details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsession)
}



func(h *HandlerStructForS3)HandleUploadChunks(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into uploading chunks")
	uploadID := r.URL.Query().Get("upload_id")
	fmt.Println("upload Id in fxn handle upload chunks is ", uploadID)
	chunknum:=0
	chunkstr := r.URL.Query().Get("chunk")
	chunknum, err := strconv.Atoi(chunkstr)
	fmt.Println("chunk number is ",chunknum)
	if err!=nil{
		http.Error(w, "Invalid chunk number", http.StatusBadRequest)
        return
	}
	upload_ack,err:=h.manager.UploadChunk(uploadID,chunknum,r.Body)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("upload ack is", upload_ack)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(upload_ack)
}


func(h *HandlerStructForS3) HandleCompleteUpload(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into Handle complete Upload")
	uploadID := r.URL.Query().Get("upload_id")
	uploadResp,err:=h.manager.CompleteUpload(uploadID)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    uploadSession,err:=h.repo.FindChunkByUploadID(r.Context(),uploadID)
	//all update in db as well
	jobsdb := model.Jobs_Database{
		VideoId: uploadResp.UploadId,
		FileName: uploadResp.FileName,
		RawPath: uploadResp.FilePath,
		Status: "Uploaded",
		Metadata: uploadSession.Metadata, 
	}
	err=h.srvc.Insert(r.Context(),&jobsdb)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//also upload to worker queue i.e redis queue
	//using async for publishing queue for transcoding as xadd redis internal implements later
	JobqueueData := &queue.JobQueue{
        VideoId: uploadResp.UploadId,
	    RawPath: uploadResp.FilePath,
	    FileName: uploadResp.FileName,
	}
	err=h.asynqsrvc.Publish(r.Context(),queue.TypeTranscode,JobqueueData)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"upload_id": uploadID,
		"status": "queued",
	})
}

func(h *HandlerStructForS3) HandleGetStatusUpload(w http.ResponseWriter, r *http.Request){
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