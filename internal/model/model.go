package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Jobs_Database struct {
	gorm.Model
	FileName string
	RawPath string
	Status string
	Metadata datatypes.JSON
}