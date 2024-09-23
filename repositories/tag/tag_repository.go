package tag

import "media-indexer/models"

type TagRepository interface {
	Create(tag *models.Tag) error
	List(offset int, limit int) ([]models.Tag, error)
	FindOrCreateTagByName(name string) (*models.Tag, error)
	FindByName(name string) (*models.Tag, error)
	Count() (int64, error)
}
