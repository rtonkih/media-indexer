package tag

import (
	"media-indexer/models"
	"media-indexer/repositories/tag"
)

type TagServiceImpl struct {
	TagRepo tag.TagRepository
}

func NewTagService(tagRepo tag.TagRepository) TagService {
	return &TagServiceImpl{TagRepo: tagRepo}
}

func (s *TagServiceImpl) CreateTag(name string) (*models.Tag, error) {
	tag := &models.Tag{Name: name}
	if err := s.TagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagServiceImpl) ListTags(page int, pageSize int) ([]models.Tag, int64, error) {
	offset := (page - 1) * pageSize
	tags, err := s.TagRepo.List(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	totalItems, err := s.TagRepo.Count()
	if err != nil {
		return nil, 0, err
	}
	return tags, totalItems, nil
}

func (s *TagServiceImpl) FindByName(name string) (*models.Tag, error) {
	return s.TagRepo.FindByName(name)
}
