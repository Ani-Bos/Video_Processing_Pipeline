package model

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Chunk_Session struct {
	gorm.Model
	UploadID       string  `gorm:"uniqueIndex"` 
	FileName       string
	TotalSize      int64
	ChunkSize      int64
	TotalChunks    int64
	Metadata   datatypes.JSON `gorm:"type:jsonb"`
}

type Chunk struct{
	gorm.Model
	UploadID string `gorm:"uniqueIndex:idx_upload_chunk;not null"`
	Index int `gorm:"column:chunk_index;uniqueIndex:idx_upload_chunk;not null"`
}