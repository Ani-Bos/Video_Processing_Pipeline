package main	

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Entering into file uploading")
	//parse our multipart form so that max limit of file is 10 MB
	r.ParseMultipartForm(10<<20)
	 // FormFile returns the first file for the given key `aniFile"
	// curl -X POST -F "anifile=@C:\Users\Aniket\Downloads\deep learning interviews.pdf" http://localhost:8080/

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