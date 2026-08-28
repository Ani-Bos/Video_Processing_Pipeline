package uploader
import (

	"fmt"
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

