package uploader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)


func UploadStreamingHandler(w http.ResponseWriter, r* http.Request){
	fmt.Println("Entering into file uploading using streaming instead of directly putting the file in memory")
	const MAX_UPLOAD_SIZE = 1<<30 // 1GB max size(1*2^30)
	const MAX_BUFFER_SIZE = 64<<10 //64 KB (64*2^10)
	//limit the request body to max upload size
	r.Body = http.MaxBytesReader(w,r.Body,MAX_UPLOAD_SIZE)
	contentType := r.Header.Get("Content-Type")
	if contentType==""{
        http.Error(w,"content-type is not set in header",http.StatusBadRequest)
		return
	}
	reader,err:=r.MultipartReader()
    if err!=nil{
		http.Error(w,"unable to create multi part reader",http.StatusBadRequest)
		return
	}
	//process each prt of request
	// #avoiding entire request buffer in memory
	for{
		next_part,err:=reader.NextPart()
		//no parts left for processing
		if err==io.EOF{
			break
		}
		if err!=nil{
           http.Error(w,"error reading part",http.StatusBadRequest)
		   return 
		return
		}
		if next_part.FileName()==""{
			continue
		}
		//processing file part
		if err:=processPart(next_part);err!=nil{
            http.Error(w, err.Error(), http.StatusInternalServerError)
		   return 
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("Upload completed successfully")
}

func processPart(p *multipart.Part)error{
	const MAX_BUFFER_SIZE = 64<<10//64 KB (64*2^10)

	filename:=filepath.Base(p.FileName())
	if filename=="." || filename=="/"{
      return fmt.Errorf("invalid filename")
	}
	// if _,err=os.Stat("tmp");os.IsNotExist(err){
	// 	os.Mkdir("tmp",0755)
	// }
	dst,err:=os.Create(filepath.Join("tmp",p.FileName()))
	if err!=nil{
		fmt.Println("Error in moving the file to target destination")
		return fmt.Errorf("Error in moving the file to target destination")
	}
	defer dst.Close()
	//creating a hash writer to check checksum while streaming

	hashr := sha256.New()
	// using multiwriter to write to both hash and file
	multiwrtr:=io.MultiWriter(dst,hashr)
	//streaming file using fixed size buffer whatever be the file size i.e. constant memory usage
	buffr:=make([]byte,MAX_BUFFER_SIZE)
	written,err:=io.CopyBuffer(multiwrtr,p,buffr)
	if err!=nil{
		return fmt.Errorf("error writing file: %w", err)
	}
	//computing final checksum
	final_checksum:=hex.EncodeToString(hashr.Sum(nil))
	fmt.Printf("File %s uploaded: %d bytes, checksum: %s\n",filename, written, final_checksum)
	return nil
}