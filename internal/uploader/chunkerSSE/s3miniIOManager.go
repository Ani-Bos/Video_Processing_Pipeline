package chunkersse

import (
	"context"
	"fmt"
	"time"
	"video_processing_pipeline/internal/model"
	"video_processing_pipeline/internal/repository"
    "io"
	"github.com/minio/minio-go/v7"
)

type MiniIOManager struct {
	repo repository.ChunkRepo
	client minio.Client
	bucketName string
}

func NewMiniIOManager(Repo repository.ChunkRepo, clnt minio.Client,bucket string)*MiniIOManager{
	return &MiniIOManager{
      repo: Repo,
	  client: clnt,
	  bucketName: bucket,
	}
}

func(m *MiniIOManager)InitiateUpload(req *RequestWrapper)(*ChunkedUpload,error){
	fmt.Println("Entering into initiate upload for MiniIO s3 equivalent")
	ctx:=context.Background()
   total_no_of_chunks := int(req.TotalSize / req.ChunkSize)
	//like say 105 MB 10 MB so total is 10 chunks
	//5 left so that why at last increased 1 to process that data as well
	if (int)(req.TotalSize%req.ChunkSize) != 0 {
		total_no_of_chunks += 1
	}
    uploadID := generateUploadID()
	// uploadPath := filepath.Join(d.uploadDir, "chunks", uploadID)
	// if err := os.MkdirAll(uploadPath, 0755); err != nil {
	// 	return nil, err
	// }

	uploadsession := &model.Chunk_Session{
		UploadID:    uploadID,
		FileName:    req.FileName,
		TotalSize:   req.TotalSize,
		ChunkSize:   req.ChunkSize,
		TotalChunks: int64(total_no_of_chunks),
	}
	err:=m.repo.Create(ctx,uploadsession)
	if err!=nil{
		return nil,err
	}

	return &ChunkedUpload{
		ID:             uploadsession.UploadID,
		FileName:       uploadsession.FileName,
		TotalSize:      uploadsession.TotalSize,
		ChunkSize:      uploadsession.ChunkSize,
		TotalChunks:    uploadsession.TotalChunks,
		UploadedChunks: make(map[int]bool),
		CreatedAt:     uploadsession.CreatedAt,
		UpdatedAt:      uploadsession.CreatedAt.Add(24 * time.Hour),
	}, nil
}

func(m *MiniIOManager)UploadChunk(uploadId string, chunkNumber int, data io.Reader)(*ChunkAcknowledgemnt,error){
  fmt.Println("entering into uploading chunks for s3 buckets")
  ctx:=context.Background()
   //find session /chunk by it exist by uploadID
   chunk_ssn,err:=m.repo.FindChunkByUploadID(ctx,uploadId)
   if err!=nil{
	 return nil,err
   }
   if chunkNumber < 0 || chunkNumber >= int(chunk_ssn.TotalChunks) {
		return nil,fmt.Errorf("Invalid chunk number")
   }
   //check if chunk is already uploaded
   uploaded_list,err1:=m.repo.FindUploadedChunks(ctx,uploadId)
   if err1!=nil{
	return nil,err1
   }
   if contains(uploaded_list,chunkNumber){
     return &ChunkAcknowledgemnt{
		ChunkNumber:    chunkNumber,
		UploadedChunks: len(uploaded_list),
		TotalChunks:    chunk_ssn.TotalChunks,
	 },nil
   }
   //create the chunk file name
	// chunkfilePath := filepath.Join(d.uploadDir, "chunks", uploadId, fmt.Sprintf("chunk_%d", chunkNumber))
	// //create a new chunk file and buffer the data intoi it
	// nextchunkfile, err := os.Create(chunkfilePath)
	// if err != nil {
	// 	return nil,err
	// }
	// defer nextchunkfile.Close()
	// //straming chunk data to file to handle
	// written,err:= io.Copy(nextchunkfile,data)
	// if err!=nil{
	// 	os.Remove(chunkfilePath)
	// 	return nil,err
	// }
	//validate chunk size except for lasdt chunk
	expected_size := chunk_ssn.ChunkSize
	if chunkNumber == int(chunk_ssn.TotalChunks-1) {
		expected_size = chunk_ssn.TotalSize - (int64(chunkNumber) * chunk_ssn.ChunkSize)
	}
	chunkobj:=fmt.Sprintf("chunks/%s/chunk_%d", uploadId, chunkNumber)
	//now pushing chunks into s3 equivalent mini.Io
	objinfo,err:=m.client.PutObject(ctx,m.bucketName,chunkobj,data,expected_size,minio.PutObjectOptions{
		ContentType:  "application/octet-stream",
	})
	if err!=nil{
		return nil,err
	}
    if objinfo.Size!=expected_size{
		err:=m.client.RemoveObject(ctx,m.bucketName,chunkobj,minio.RemoveObjectOptions{
		})
		return nil, err
	}
    err3:=m.repo.MarkChunkUploaded(ctx,uploadId,chunkNumber)
	if err3!=nil{
		return nil,err3
	}
	uploadcnt,err4:=m.repo.GetUploadedCount(ctx,uploadId)
	if err4!=nil{
		return nil,err4
	}
	return &ChunkAcknowledgemnt{
              ChunkNumber: chunkNumber,
              UploadedChunks: int(uploadcnt),
              TotalChunks:chunk_ssn.TotalChunks,
	},nil
}

func(m *MiniIOManager)CompleteUpload(UploadId string)(*UploadResponse,error){
	fmt.Println("Entering into completing miniIO upload all chunks")
	ctx:=context.Background()
   //find session /chunk by it exist by uploadID
   chunk_ssn,err:=m.repo.FindChunkByUploadID(ctx,UploadId)
   if err!=nil{
	 return nil,err
   } 
   upldcnt,err1:=m.repo.GetUploadedCount(ctx,UploadId)
   if err1!=nil{
	return nil,err1
   }
   if upldcnt!=int64(chunk_ssn.TotalChunks){
     	return nil,fmt.Errorf("incomplete upload: %d/%d chunks",  upldcnt, chunk_ssn.TotalChunks)
   }
//    //create and assembling all chunks in final file chunks
// 	FinalchunkfilePath := filepath.Join(d.uploadDir, chunk_ssn.FileName)
// 	//create a new chunk file and buffer the data intoi it
// 	Finalchunkfile, err := os.Create(FinalchunkfilePath)
// 	if err != nil {
// 		return nil,err
// 	}
// 	defer Finalchunkfile.Close()
// 	//assembling chunks in all order
// 	chunksdir:=filepath.Join(d.uploadDir,"chunks",UploadId)
    var cpyfilearr []minio.CopySrcOptions
	for i:=0;i<int(chunk_ssn.TotalChunks);i++{
		chunkobj:=fmt.Sprintf("chunks/%s/chunk_%d",UploadId,i)
        cpyfilearr = append(cpyfilearr, minio.CopySrcOptions{
			Bucket: m.bucketName,
			Object: chunkobj,
		})
	}
	finalObj:=chunk_ssn.FileName
	_,err=m.client.ComposeObject(ctx,minio.CopyDestOptions{
		Bucket: m.bucketName,
		Object: finalObj,
	},cpyfilearr...)
	if err!=nil{
		return nil,fmt.Errorf("Failed to assemble multiple temporary chunk into final one",err)
	}
	//cleanup upload session and remove all chunks by async way
	// os.RemoveAll(chunksdir)
	// m.mu.Lock()
	// delete(m.uploads,UploadId)
	// m.mu.Unlock()
    go func(){
      clnupctx:=context.Background()
	  for i:=0;i<int(chunk_ssn.TotalChunks);i++{
		chunkobjkey:=fmt.Sprintf("chunks/%s/chunk_%d",UploadId,i)
		_=m.client.RemoveObject(clnupctx,m.bucketName,chunkobjkey,minio.RemoveObjectOptions{
		})
	  }
	}()
	return &UploadResponse{
		UploadId:UploadId,
		FileName:chunk_ssn.FileName,
		Size:chunk_ssn.TotalSize,
		FilePath:fmt.Sprintf("%s/%s",m.bucketName,finalObj),
	},nil
	}
	
func(m *MiniIOManager)GetUploadStatus(UploadId string)(*UploadStatus,error){
  
	ctx:=context.Background()
   //find session /chunk by it exist by uploadID
   chunk_ssn,err:=m.repo.FindChunkByUploadID(ctx,UploadId)
   if err!=nil{
	 return nil,err
   } 
   upload_chunk_list,err1:=m.repo.FindUploadedChunks(ctx,UploadId)
   if err1!=nil{
	return nil,err1
   }
	exist := make(map[int]bool, len(upload_chunk_list))
	for _, i := range upload_chunk_list{
		exist[i] = true
	}
	 missingChunks := make([]int,0)
	for i:=0 ; i<int(chunk_ssn.TotalChunks);i++{
		if exist[i]==false {
			missingChunks = append(missingChunks, i)
		}
	}
	return &UploadStatus{
		UploadId:UploadId,
		FileName:chunk_ssn.FileName,
		UploadedChunks:len(upload_chunk_list),
		TotalChunks:chunk_ssn.TotalChunks,
		MissingChunks:missingChunks,
		IsComplete:len(missingChunks)==0,
	},nil
}
