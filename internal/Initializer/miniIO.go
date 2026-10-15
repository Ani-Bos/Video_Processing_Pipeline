package initializer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMiniIOAndCreateBucket(ctx context.Context, bucket string)(*minio.Client,error){
	fmt.Println("Entering into connect minio IO package")
	host:="minio:9000"
	// host:=os.Getenv("MINIO_HOST")
	// fmt.Println("gettig host details", host)
	access_key:=os.Getenv("S3_ACCESS_KEY")
	secret_key:=os.Getenv("S3_SECRET_KEY")
	client,err:=minio.New(host,&minio.Options{
		Creds: credentials.NewStaticV4(access_key,secret_key,""),
		Secure: false,
	})
	if err!=nil{
		fmt.Errorf("Failed to connect and authneticate with minio chainguard client",err)
		return nil,err
	}
	fmt.Println("successfully initialized mminiIO client chainguard")
	var exists bool
	for i:=0;i<5;i++{
      exists,err=client.BucketExists(ctx,bucket)
	  if err==nil{
		break
	  }
	  fmt.Println("Adding retries in minio startup sleeping for 2 sec",i)
	  time.Sleep(2*time.Second)
	}
	if err!=nil{
		return nil,fmt.Errorf("Failed to get the bucket list",err)
	}
	if !exists{
		fmt.Println("Entering into creating buckets")
     err:=client.MakeBucket(ctx,bucket,minio.MakeBucketOptions{})
	 if err!=nil{
		return nil,fmt.Errorf("Failed to create bucket",err)
	 }
	 fmt.Println("Bucket creation done")
	}else{
      fmt.Printf("Bucket name already exist no need to create")
	}
	return client,nil
}