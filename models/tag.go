package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string    `gorm:"uniqueIndex;size:255"`
}

func (tag *Tag) BeforeCreate(_tx *gorm.DB) (err error) {
	tag.ID = uuid.New()
	return
}
