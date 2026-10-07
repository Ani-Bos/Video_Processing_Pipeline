package mediastore

import (
	"context"
	"fmt"
	"io"
)

type MediaStore struct {
	Media MediaStoreManager
}

func (M *MediaStore)Put(ctx context.Context, bucket, key string, r io.Reader, contentType string) error {
  fmt.Println("Entering into putting the key into the minio bucket")
  return nil
}

func(M *MediaStore)Get(ctx context.Context,bucket,key string)(*io.Reader,error){
  fmt.Println("Entering into getting the key from bucket may be chunk uploaded")
  return nil,nil
}