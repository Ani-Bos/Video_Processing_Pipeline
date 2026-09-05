package service

import (
	"context"
	"video_processing_pipeline/internal/model"
	"video_processing_pipeline/internal/repository"
)

type ChunkService struct {
	Repo repository.ChunkUploadManager
}


func(c *ChunkService)Insert(ctx context.Context, mdl *model.Chunk_Session)error{
	return c.Repo.Create(ctx,mdl)
}
