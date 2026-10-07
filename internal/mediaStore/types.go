package mediastore

import (
	"context"
	"io"
)

type MediaStoreManager interface {
	Put(ctx context.Context, bucket, key string, r io.Reader, contentType string)error
	Get(ctx context.Context,bucket,key string)(*io.Reader,error)
}