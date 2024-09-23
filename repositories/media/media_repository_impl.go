package media

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-indexer/models"
)

type MediaRepositoryImpl struct {
	DB *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepositoryImpl {
	return &MediaRepositoryImpl{DB: db}
}

func (r *MediaRepositoryImpl) Create(media *models.Media) error {
	return r.DB.Create(media).Error
}

func (r *MediaRepositoryImpl) FindByTagNames(tagNames []string, page int, pageSize int) ([]models.Media, int64, error) {
	var mediaList []models.Media
	offset := (page - 1) * pageSize

	var totalItems int64
	countQuery := r.DB.Model(&models.Media{}).
		Joins("JOIN media_tags ON media_tags.media_id::uuid = media.id::uuid").
		Where("media_tags.tag_name::text IN ?", tagNames).
		Group("media.id").
		Count(&totalItems)

	if countQuery.Error != nil {
		return nil, 0, countQuery.Error
	}

	// If no items found, return empty result
	if totalItems == 0 {
		return []models.Media{}, 0, nil
	}

	err := r.DB.Model(&models.Media{}).
		Joins("JOIN media_tags ON media_tags.media_id::uuid = media.id::uuid").
		Where("media_tags.tag_name::text IN ?", tagNames).
		Group("media.id").
		Offset(offset).
		Limit(pageSize).
		Preload("Tags").
		Find(&mediaList).Error

	if err != nil {
		return nil, 0, err
	}

	return mediaList, totalItems, nil
}

func (r *MediaRepositoryImpl) AssociateMediaWithTag(mediaID uuid.UUID, tagID uuid.UUID, tagName string) error {
	mediaTag := models.MediaTag{
		MediaID: mediaID,
		TagID:   tagID,
		TagName: tagName,
	}
	return r.DB.Create(&mediaTag).Error
}
