package repository

import (
	"video_processing_pipeline/internal/model"

	"gorm.io/gorm"
)

type Jobs_Repository interface {
	create(jobdb *model.Jobs_Database)error
	update(jobdb *model.Jobs_Database)error
	FindByVideoID(videoId string)(*model.Jobs_Database,error)
	FindByStatus(status string)([]*model.Jobs_Database,error)
}

type JobQueueDB struct {
	DB *gorm.DB
}