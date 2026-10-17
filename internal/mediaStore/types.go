package mediastore

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type MediaStoreManager interface {
	PresignedGet(ctx context.Context,key string )(string,error)
	DownloadObject(ctx context.Context,key string,dest string)error
	PutObjects(ctx context.Context,key string)(*minio.UploadInfo,error)
}