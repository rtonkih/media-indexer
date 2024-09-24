package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string    `gorm:"not null"`
	Link string    `gorm:"not null"`
	Tags []Tag     `gorm:"many2many:media_tags;"`
}

func (media *Media) BeforeCreate(_tx *gorm.DB) (err error) {
	media.ID = uuid.New()
	return
}
