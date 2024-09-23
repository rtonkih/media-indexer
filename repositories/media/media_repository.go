package media

import (
	"github.com/google/uuid"

	"media-indexer/models"
)

type MediaRepository interface {
	Create(media *models.Media) error
	FindByTagNames(tagNames []string, page int, pageSize int) ([]models.Media, int64, error)
	AssociateMediaWithTag(mediaID uuid.UUID, tagID uuid.UUID, tagName string) error
}
