package tags

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"media-indexer/models"
	tagRepo "media-indexer/repositories/tag"
	"media-indexer/services/tag"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Tag{})
	return db
}

func SetupRouter(tagController *TagController) *gin.Engine {
	router := gin.Default()
	router.POST("/tags", tagController.CreateTag)
	router.GET("/tags", tagController.ListTags)
	return router
}

func SetupTestEnv() (*gorm.DB, *gin.Engine) {
	db := SetupTestDB()
	tagRepo := tagRepo.NewTagRepository(db)
	tagService := tag.NewTagService(tagRepo)
	tagController := NewTagController(tagService)
	router := SetupRouter(tagController)
	return db, router
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestCreateTag(t *testing.T) {
	_, router := SetupTestEnv()

	tag := map[string]string{"name": "Mbappe"}
	tagJSON, _ := json.Marshal(tag)

	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer(tagJSON))

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var responseBody TagResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, "mbappe", responseBody.Name)
}

func TestCreateTagAlreadyExists(t *testing.T) {
	db, router := SetupTestEnv()

	tag := models.Tag{Name: "mbappe"}
	db.Create(&tag)

	tagRequest := map[string]string{"name": "mbappe"}
	tagJSON, _ := json.Marshal(tagRequest)

	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer(tagJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)

	var responseBody TagResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, "mbappe", responseBody.Name)
}

func TestListTags(t *testing.T) {
	db, router := SetupTestEnv()

	tags := []models.Tag{
		{Name: "ronaldo"},
		{Name: "messi"},
		{Name: "neymar"},
		{Name: "mbappe"},
	}
	db.Create(&tags)

	req, _ := http.NewRequest(http.MethodGet, "/tags?page=1&pageSize=2", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var responseBody PaginatedTagsResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Len(t, responseBody.Tags, 2)
	assert.Equal(t, "ronaldo", responseBody.Tags[0].Name)
	assert.Equal(t, "messi", responseBody.Tags[1].Name)
	assert.Equal(t, 1, responseBody.Page)
	assert.Equal(t, 2, responseBody.PageSize)
	assert.Equal(t, int64(4), responseBody.TotalItems)
	assert.Equal(t, 2, responseBody.TotalPages)

	req, _ = http.NewRequest(http.MethodGet, "/tags?page=2&pageSize=2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Len(t, responseBody.Tags, 2)
	assert.Equal(t, "neymar", responseBody.Tags[0].Name)
	assert.Equal(t, "mbappe", responseBody.Tags[1].Name)
	assert.Equal(t, 2, responseBody.Page)
	assert.Equal(t, 2, responseBody.PageSize)
	assert.Equal(t, int64(4), responseBody.TotalItems)
	assert.Equal(t, 2, responseBody.TotalPages)

}
