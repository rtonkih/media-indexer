package media

import (
	"github.com/google/uuid"

	"media-indexer/models"
)

type MediaService interface {
	CreateMedia(name string, filePath string, tagNames []string) (*models.Media, []models.Tag, error)
	SearchMediaByTags(tagNames []string, page int, pageSize int) ([]models.Media, int64, error)
	FetchOrCreateTagsAndAssociate(mediaID uuid.UUID, tagNames []string) ([]models.Tag, error)
}
