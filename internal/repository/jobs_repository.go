package repository

import (
	"context"
	"errors"
	"video_processing_pipeline/internal/model"

	"gorm.io/gorm"
)

type Jobs_Repository interface {
	Create(ctx context.Context,jobdb *model.Jobs_Database)error
	Update(ctx context.Context,jobdb *model.Jobs_Database)error
	FindByVideoID(ctx context.Context,videoId string)(*model.Jobs_Database,error)
	FindByStatus(ctx context.Context,status string)([]*model.Jobs_Database,error)
}

type JobRepo struct {
	DB *gorm.DB
}

func(j *JobRepo)Create(ctx context.Context,jobdb *model.Jobs_Database)error{
   return j.DB.Create(jobdb).Error
}

func(j *JobRepo)Update(ctx context.Context,jobdb *model.Jobs_Database)error{
  return j.DB.Save(jobdb).Error
}

func(j *JobRepo)FindByVideoID(ctx context.Context,videoId string)(*model.Jobs_Database,error){
  if videoId==""{
    return nil, errors.New("videoID cant be null or empty string")
  }
  //select * from db wherre id=videoid;
  var jdb model.Jobs_Database
  err:=j.DB.Where(&model.Jobs_Database{VideoId: videoId}).First(&jdb).Error
  return &jdb,err
}

func(j *JobRepo)FindByStatus(ctx context.Context,status string)([]*model.Jobs_Database,error){
  if status==""{
    return nil, errors.New("Status cant be empty string")
  }
  var jobs []*model.Jobs_Database
  // select * from db wherre Status=status;
  err:=j.DB.Where(&model.Jobs_Database{Status: status}).Find(&jobs).Error
  return jobs,err
}

