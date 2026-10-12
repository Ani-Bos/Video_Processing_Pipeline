package initializer

import (
	"fmt"
	"log"
	"os"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMiniIOAndCreateBucket(){
	fmt.Println("Entering into connect minio IO package")
	host:=os.Getenv("MINIO_HOST")
	access_key:=os.Getenv("MINIO_USER")
	secret_key:=os.Getenv("MINIO_PASSWORD")
	_,err:=minio.New(host,&minio.Options{
		Creds: credentials.NewStaticV4(access_key,secret_key,""),
	})
	if err!=nil{
		log.Fatalf("Failed to connect and authneticate with minio chainguard client",err)
		return 
	}
	fmt.Println("successfully initialized mminiIO client chainguard")
}