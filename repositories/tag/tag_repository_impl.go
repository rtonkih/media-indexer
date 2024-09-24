package tag

import (
	"gorm.io/gorm"

	"media-indexer/models"
)

type TagRepositoryImpl struct {
	DB *gorm.DB
}

func (r *TagRepositoryImpl) FindOrCreateTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.DB.Where("name = ?", name).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tag = models.Tag{Name: name}
			if err := r.DB.Create(&tag).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &tag, nil
}

func NewTagRepository(db *gorm.DB) *TagRepositoryImpl {
	return &TagRepositoryImpl{DB: db}
}

func (r *TagRepositoryImpl) Create(tag *models.Tag) error {
	return r.DB.Create(tag).Error
}

func (r *TagRepositoryImpl) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryImpl) List(offset int, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.DB.Offset(offset).Limit(limit).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&models.Tag{}).Count(&count).Error
	return count, err
}
