package service

import (
	"context"
	"fmt"
	"video_processing_pipeline/internal/model"
	"video_processing_pipeline/internal/repository"
)

type InterfaceInjectRepoJob struct{
	Repo repository.Jobs_Repository
}

func(srvc *InterfaceInjectRepoJob)Insert(ctx context.Context,jb *model.Jobs_Database)(error){
  fmt.Println("Enterinto into service layer for inserting into jobs data")
  return srvc.Repo.Create(ctx,jb)
}