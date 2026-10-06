package model

import (
	"time"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Jobs_Database struct {
	gorm.Model
	VideoId    string
	FileName   string
	RawPath    string
	Status     string
	Stage      string
	Error      string
	NotifyEmail string
	NotifiedAt *time.Time
	Metadata   datatypes.JSON `gorm:"type:jsonb"`
}