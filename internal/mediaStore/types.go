package mediastore

import (
	"context"
)

type MediaStoreManager interface {
	PresignedGet(ctx context.Context,key string )(string,error)
	DownloadObject(ctx context.Context,key string,dest string)error
	PutDirObjects(ctx context.Context,key string, pref string)(error)
	DownloadHLSObject(ctx context.Context,key string, pref string)(error)
	PutFileObject(ctx context.Context,key string,dest string)(error)
}