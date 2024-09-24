package tags

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"media-indexer/services/tag"
	"media-indexer/utils"
)

type TagController struct {
	TagService tag.TagService
}

func NewTagController(tagService tag.TagService) *TagController {
	return &TagController{TagService: tagService}
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedTagsResponse struct {
	Tags       []TagResponse `json:"tags"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalItems int64         `json:"totalItems"`
	TotalPages int           `json:"totalPages"`
}

type CreateTagRequest struct {
	Name string `form:"name" binding:"required"`
}

// @BasePath /api/v1
// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag with the given name
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body object true "Tag name"
// @Success 201 {object} TagResponse "Created tag"
// @Success 409 {object} TagResponse "Tag already exists"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Failed to create tag"
// @Router /tags [post]
func (tc *TagController) CreateTag(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	normalizedTagName := utils.NormalizeTag(input.Name)

	existingTag, err := tc.TagService.FindByName(normalizedTagName)
	if err == nil && existingTag != nil {
		c.JSON(http.StatusConflict, TagResponse{
			ID:   existingTag.ID,
			Name: existingTag.Name,
		})
		return
	}

	tag, err := tc.TagService.CreateTag(normalizedTagName)
	if err != nil {
		log.Printf("Error creating tag: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	})
}

// @BasePath /api/v1
// ListTags godoc
// @Summary List all tags
// @Description Retrieve a list of all tags with pagination
// @Tags tags
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of tags per page"
// @Success 200 {object} PaginatedTagsResponse "List of tags with pagination"
// @Failure 500 {object} ErrorResponse "Failed to retrieve tags"
// @Router /tags [get]
func (tc *TagController) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page <= 0 || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number or page size"})
		return
	}

	tags, totalItems, err := tc.TagService.ListTags(page, pageSize)
	if err != nil {
		log.Printf("Error retrieving tags: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve tags"})
		return
	}

	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	var tagResponses []TagResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	c.JSON(http.StatusOK, PaginatedTagsResponse{
		Tags:       tagResponses,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}
