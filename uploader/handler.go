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

//isse with this is it buffers the file in memory
// large file stored in temporary files avoiding IO overhead
func UploadHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into file uploading")
	//parse our multipart form so that max limit of file is 10 MB
	//bitwise left--10*2^20--2^20--1 MB--10 MB
	r.ParseMultipartForm(10<<20)
	 // FormFile returns the first file for the given key `aniFile"
	// curl -X POST -F "anifile=@c:\Users\Aniket\Downloads\TUN_6417.JPG" http://localhost:8080/
	//curl -X POST http://localhost:8080/upload -H "Content-Type: multipart/form-data" -F "file=@c:\Users\Aniket\Documents\1080_30_8.00_Jun222021(1).mp4"

	file,header,err := r.FormFile("anifile")
	if err!=nil{
		fmt.Println("Error in file retrieving")
	    return 
	}
	defer file.Close()
	fmt.Println(header.Filename)
	fmt.Println(header.Header)
	fmt.Println(header.Size)

	//saving it locally
	if _,err=os.Stat("tmp");os.IsNotExist(err){
		os.Mkdir("tmp",0755)
	}
	dst,err:=os.Create(filepath.Join("tmp",header.Filename))
	if err!=nil{
		fmt.Println("Error in moving the file to target destination")
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _,err:=os.ReadFile(dst.Name());err!=nil{
		fmt.Println("Error in savingthe file")
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
}

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