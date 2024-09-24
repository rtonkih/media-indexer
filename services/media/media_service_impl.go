package media

import (
	"github.com/google/uuid"

	"media-indexer/models"
	"media-indexer/repositories/media"
	"media-indexer/repositories/tag"
	"media-indexer/utils"
)

type MediaServiceImpl struct {
	MediaRepo media.MediaRepository
	TagRepo   tag.TagRepository
}

func NewMediaService(mediaRepo media.MediaRepository, tagRepo tag.TagRepository) MediaService {
	return &MediaServiceImpl{MediaRepo: mediaRepo, TagRepo: tagRepo}
}

func (s *MediaServiceImpl) CreateMedia(name string, fileName string, tagNames []string) (*models.Media, []models.Tag, error) {
	normalizedTagNames := make([]string, len(tagNames))
	for i, tagName := range tagNames {
		normalizedTagNames[i] = utils.NormalizeTag(tagName)
	}

	media := models.Media{
		Name: name,
		Link: fileName,
	}

	if err := s.MediaRepo.Create(&media); err != nil {
		return nil, nil, err
	}

	tags, err := s.FetchOrCreateTagsAndAssociate(media.ID, normalizedTagNames)
	if err != nil {
		return nil, nil, err
	}

	return &media, tags, nil
}

func (s *MediaServiceImpl) SearchMediaByTags(tagNames []string, page int, pageSize int) ([]models.Media, int64, error) {
	normalizedTagNames := make([]string, len(tagNames))
	for i, tagName := range tagNames {
		normalizedTagNames[i] = utils.NormalizeTag(tagName)
	}
	return s.MediaRepo.FindByTagNames(normalizedTagNames, page, pageSize)
}

func (s *MediaServiceImpl) FetchOrCreateTagsAndAssociate(mediaID uuid.UUID, tagNames []string) ([]models.Tag, error) {
	var tags []models.Tag
	for _, tagName := range tagNames {
		tag, err := s.TagRepo.FindOrCreateTagByName(tagName)
		if err != nil {
			return nil, err
		}

		if err := s.MediaRepo.AssociateMediaWithTag(mediaID, tag.ID, tag.Name); err != nil {
			return nil, err
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}
