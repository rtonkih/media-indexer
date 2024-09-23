package models

import "github.com/google/uuid"

type MediaTag struct {
	MediaID uuid.UUID `gorm:"primaryKey;index"`
	TagID   uuid.UUID `gorm:"primaryKey"`
	TagName string    `gorm:"index"`
}
