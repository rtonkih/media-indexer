package media

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"media-indexer/services/media"
	"media-indexer/storage"
)

type MediaController struct {
	MediaService media.MediaService
	Storage      storage.StorageProvider
}

func NewMediaController(mediaService media.MediaService) *MediaController {
	storageProvider := storage.StorageFactory()

	return &MediaController{MediaService: mediaService, Storage: storageProvider}
}

type MediaResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Link string    `json:"link"`
	Tags []string  `json:"tags"`
}

type SearchMediaResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Tags    []string  `json:"tags"`
	FileURL string    `json:"fileUrl"`
}

type PaginatedMediaResponse struct {
	Media      []SearchMediaResponse `json:"media"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalItems int64                 `json:"totalItems"`
	TotalPages int                   `json:"totalPages"`
}

// CreateMedia godoc
// @Summary Create media
// @Description Create a new media item with associated tags
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name of the media"
// @Param tags formData array true "Tags associated with the media"
// @Param file formData file true "File to upload"
// @Success 201 {object} MediaResponse "Created media"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 409 {object} gin.H "Conflict"
// @Router /media [post]
func (mc *MediaController) CreateMedia(c *gin.Context) {
	var input struct {
		Name string                `form:"name" binding:"required"`
		Tags []string              `form:"tags" binding:"required"`
		File *multipart.FileHeader `form:"file" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := filepath.Base(input.File.Filename)
	fileURL, err := mc.Storage.UploadFile(input.File, fileName)
	if err != nil {
		log.Printf("failed to upload file to storage: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file to storage"})
		return
	}

	media, tags, err := mc.MediaService.CreateMedia(input.Name, fileURL, input.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create media"})
		return
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	c.JSON(http.StatusCreated, MediaResponse{
		ID:   media.ID,
		Name: media.Name,
		Link: media.Link,
		Tags: tagNames,
	})
}

// SearchMediaByTag godoc
// @Summary Search media by tag
// @Description Search for media items by tag name
// @Tags media
// @Accept json
// @Produce json
// @Param tag query array true "Tag name(s) to search for"
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of media items per page"
// @Success 200 {object} PaginatedMediaResponse "Search results"
// @Failure 400 {object} gin.H "Bad Request"
// @Router /media [get]
func (mc *MediaController) SearchMediaByTag(c *gin.Context) {
	tagNames := c.QueryArray("tag")
	if len(tagNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag names are required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page <= 0 || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number or page size"})
		return
	}

	media, totalItems, err := mc.MediaService.SearchMediaByTags(tagNames, page, pageSize)
	if err != nil {
		log.Printf("Error searching media by tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search media"})
		return
	}

	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	var mediaResponses []SearchMediaResponse
	for _, m := range media {
		var tags []string
		for _, tag := range m.Tags {
			tags = append(tags, tag.Name)
		}
		mediaResponses = append(mediaResponses, SearchMediaResponse{
			ID:      m.ID,
			Name:    m.Name,
			Tags:    tags,
			FileURL: m.Link,
		})
	}

	c.JSON(http.StatusOK, PaginatedMediaResponse{
		Media:      mediaResponses,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}
