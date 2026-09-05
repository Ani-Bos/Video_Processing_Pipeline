package model

import "gorm.io/gorm"

type Chunk_Session struct {
	gorm.Model
	UploadID       string  `gorm:"uniqueIndex"` 
	FileName       string
	TotalSize      int64
	ChunkSize      int64
	TotalChunks    int64
	UploadedChunks []Chunk `gorm:"foreignKey:UploadID;references:UploadID"`
}

type Chunk struct{
	gorm.Model
	UploadID string `gorm:"index"`
	Index int
}