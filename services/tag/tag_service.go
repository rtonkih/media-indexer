package tag

import "media-indexer/models"

type TagService interface {
	CreateTag(name string) (*models.Tag, error)
	ListTags(page int, pageSize int) ([]models.Tag, int64, error)
	FindByName(name string) (*models.Tag, error)
}
