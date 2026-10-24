package mediastore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/minio/minio-go/v7"
)

type MediaStore struct {
	Media MediaStoreManager
	Client *minio.Client
	Bucket string
}

func NewMediaStore(Md MediaStoreManager, clnt *minio.Client, Bckt string)*MediaStore{
	return &MediaStore{
      Media: Md,
	  Client: clnt,
	  Bucket: Bckt,
	}
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
	//jbq.RawPath-->key,localfile-->destination
	err=M.Client.FGetObject(ctx,M.Bucket,key,dest,minio.GetObjectOptions{})
	if err!=nil{
		return err
	}
	return nil
}

func(M *MediaStore)PutDirObjects(ctx context.Context,key string,pref string)(error){
	fmt.Println("entering into put objects wrapper for moving the hls format to put into bucket after transcoding and consume in thumbnai generation to prevent concurrency issue")
	localfiles,err:=os.ReadDir(key)
	fmt.Println("localfiles under this dir created by hls one byt thumbnail;", localfiles)
	if err!=nil{
		return err
	}
	for _, i := range(localfiles){
     if i.IsDir(){
		continue
	 }
	 path:=filepath.Join(key,i.Name())
	 fmt.Println("Each filepath under the directory is ",path)
	 file,err:=os.Open(path)
	 if err!=nil{
		return nil
	 }
	 fileinfo,err:=file.Stat()
	 if err!=nil{
		file.Close()
		return err
	 }
	 prefixFilename:=pref+"/"+i.Name()
	 fmt.Println("Prefixfilename for esch file ",prefixFilename)
     _,err=M.Client.PutObject(ctx,M.Bucket,prefixFilename,file,fileinfo.Size(),minio.PutObjectOptions{
	})
	file.Close()
	if err!=nil{
		return err
	}
	}
	return nil
}

func(M *MediaStore)DownloadHLSObject(ctx context.Context,key string, pref string)(error){
	fmt.Println("Entering into getting all hls output by transcoding in the prefix hls directory")
	err:=os.MkdirAll(key,0755)
	if err!=nil{
		return err
	}
	for i:= range M.Client.ListObjects(ctx,M.Bucket,minio.ListObjectsOptions{Prefix: pref+"/",Recursive: true}){
      if i.Err!=nil{
		return i.Err
	  }
	  filename:=strings.TrimPrefix(i.Key,pref+"/")
	  dest:=filepath.Join(key,filename)
	  fmt.Println("filename in thubnail for downloadhls object is",filename)
	  fmt.Println("dest in thumbnail for downloadhls obj is",dest)
	  err=M.Client.FGetObject(ctx,M.Bucket,i.Key,dest,minio.GetObjectOptions{})
	  if err!=nil{
		return err
	  }
	}
	return nil
}
func(M *MediaStore)PutFileObject(ctx context.Context,key string,dest string)(error){
	fmt.Println("Entering into putting the output thubnail image in the minio bucket")
	file,err:=os.Open(dest)
	fmt.Println("Putting the thumbnail image in the file locn",file)
	if err!=nil{
		return err
	}
	defer file.Close()
	fileinfo,err:=file.Stat()
	if err!=nil{
		return err
	}
	fmt.Println("stat info of the file is",fileinfo)
	_,err=M.Client.PutObject(ctx,M.Bucket,key,file,fileinfo.Size(),minio.PutObjectOptions{})
	if err!=nil{
		return err
	}
	return nil
}