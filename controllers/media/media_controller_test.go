package media

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"media-indexer/models"
	"media-indexer/services/media"
	"media-indexer/storage"
)

type MockMediaService struct{}

func (m *MockMediaService) CreateMedia(name string, filePath string, _tagNames []string) (*models.Media, []models.Tag, error) {
	media := &models.Media{Name: name, Link: filePath}
	tags := []models.Tag{{Name: "tag1"}, {Name: "tag2"}}
	return media, tags, nil
}

func (m *MockMediaService) SearchMediaByTags(_tagNames []string, page int, pageSize int) ([]models.Media, int64, error) {
	media := []models.Media{
		{Name: "Arsenal", Link: "https://s3.amazonaws.com/bucket/media_1.jpg", Tags: []models.Tag{{Name: "arsenal-mu"}, {Name: "penalty"}}},
		{Name: "MU", Link: "https://s3.amazonaws.com/bucket/media_2.jpg", Tags: []models.Tag{{Name: "arsenal-mu"}, {Name: "goal"}}},
	}
	totalItems := int64(len(media))
	return media, totalItems, nil
}

func (m *MockMediaService) FetchOrCreateTagsAndAssociate(mediaID uuid.UUID, tagNames []string) ([]models.Tag, error) {
	tags := []models.Tag{{Name: "tag1"}, {Name: "tag2"}}
	return tags, nil
}

type MockStorageProvider struct{}

func (m *MockStorageProvider) UploadFile(fileHeader *multipart.FileHeader, filename string) (string, error) {
	return "http://example.com/" + filename, nil
}

func SetupMediaTestRouter(mediaService media.MediaService, storageProvider storage.StorageProvider) *gin.Engine {
	router := gin.Default()
	mediaController := &MediaController{MediaService: mediaService, Storage: storageProvider}
	router.POST("/media", mediaController.CreateMedia)
	router.GET("/media", mediaController.SearchMediaByTag)
	return router
}

func SetupMockServices() (media.MediaService, storage.StorageProvider) {
	mockMediaService := &MockMediaService{}
	mockStorageProvider := &MockStorageProvider{}
	return mockMediaService, mockStorageProvider
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestCreateMedia(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test media")
	writer.WriteField("tags", "tag1")
	writer.WriteField("tags", "tag2")
	part, _ := writer.CreateFormFile("file", "testfile.txt")
	part.Write([]byte("file content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response MediaResponse
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)

	assert.Equal(t, "test media", response.Name)
	assert.Equal(t, "http://example.com/testfile.txt", response.Link)
	assert.ElementsMatch(t, []string{"tag1", "tag2"}, response.Tags)
}

func TestCreateMedia_NoNameProvided(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("tags", "tag1")
	writer.WriteField("tags", "tag2")
	part, _ := writer.CreateFormFile("file", "testfile.txt")
	part.Write([]byte("file content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)

	assert.Contains(t, response, "error")
}

func TestCreateMedia_NoTagsProvided(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "test media")
	part, _ := writer.CreateFormFile("file", "testfile.txt")
	part.Write([]byte("file content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)

	assert.Contains(t, response, "error")
}

func TestSearchMediaByTag(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	req, _ := http.NewRequest(http.MethodGet, "/media?tag=arsenal&page=1&pageSize=1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseBody PaginatedMediaResponse
	json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.Len(t, responseBody.Media, 2)
	assert.Equal(t, "Arsenal", responseBody.Media[0].Name)
	assert.Equal(t, "https://s3.amazonaws.com/bucket/media_1.jpg", responseBody.Media[0].FileURL)
	assert.ElementsMatch(t, []string{"arsenal-mu", "penalty"}, responseBody.Media[0].Tags)
	assert.Equal(t, 1, responseBody.Page)
	assert.Equal(t, 1, responseBody.PageSize)
	assert.Equal(t, int64(2), responseBody.TotalItems)
	assert.Equal(t, 2, responseBody.TotalPages)

}

func TestSearchMediaByTag_NoTagProvided(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	req, _ := http.NewRequest(http.MethodGet, "/media", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.Equal(t, "Tag names are required", responseBody["error"])
}

func TestSearchMediaByTag_InvalidPageNumber(t *testing.T) {
	mediaService, storageProvider := SetupMockServices()
	router := SetupMediaTestRouter(mediaService, storageProvider)

	req, _ := http.NewRequest(http.MethodGet, "/media?tag=arsenal&page=-1&pageSize=1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.Equal(t, "Invalid page number or page size", responseBody["error"])
}
