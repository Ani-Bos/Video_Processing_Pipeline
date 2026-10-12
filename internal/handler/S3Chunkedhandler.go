package handler

import (
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"
	"video_processing_pipeline/internal/service"
	"video_processing_pipeline/internal/uploader/chunkersse"
)

type HandlerStructForS3 struct {
	manager chunkersse.MiniIOManager
	srvc service.InterfaceInjectRepoJob
	repo repository.ChunkRepo
	asynqsrvc queue.AsyncPublishersrvc
}

func NewHandlerStructForS3(mgr chunkersse.MiniIOManager, srvce service.InterfaceInjectRepoJob, repostry repository.ChunkRepo, asysrvc queue.AsyncPublishersrvc)*HandlerStructForS3{
	return &HandlerStructForS3{
       manager: mgr,
	   srvc: srvce,
	   repo: repostry,
	   asynqsrvc: asysrvc,
	}
}