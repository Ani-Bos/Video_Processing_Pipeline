package mediastore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

type MediaStore struct {
	Media MediaStoreManager
	Client *minio.Client
	Bucket string
}

func (M *MediaStore)PresignedGet(ctx context.Context, key string )(string,error){
  fmt.Println("Entering into getting presigned url")
  url,err:=M.Client.PresignedGetObject(ctx,M.Bucket,key,24*time.Hour,nil)
  if err!=nil{
	return "",err
  }
  return url.String(),nil
}

func(M *MediaStore)DownloadObject(ctx context.Context,key string,dest string)error{
	fmt.Println("Entering into getting object from bucket")
	err:=os.MkdirAll(filepath.Dir(dest),0755)
	if err!=nil{
		return err
	}
	err=M.Client.FGetObject(ctx,M.Bucket,key,dest,minio.GetObjectOptions{})
	if err!=nil{
		return err
	}
	return nil
}