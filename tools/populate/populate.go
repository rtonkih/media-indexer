package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"

	"media-indexer/config"
	"media-indexer/models"
)

const (
	NumTags   = 1000000
	BatchSize = 10000
	NumMedia  = 1000000
)

func init() {
	config.ConnectDB()
}

func batchSaveTags(db *gorm.DB, tags []string, batchSize int) error {
	for i := 0; i < len(tags); i += batchSize {
		end := i + batchSize
		if end > len(tags) {
			end = len(tags)
		}
		batch := tags[i:end]
		tagModels := make([]models.Tag, len(batch))
		for j, tagName := range batch {
			tagModels[j] = models.Tag{Name: tagName}
		}
		if err := db.Create(&tagModels).Error; err != nil {
			return err
		}
	}
	return nil
}

func batchSaveMedia(db *gorm.DB, mediaItems []models.Media, batchSize int) error {
	for i := 0; i < len(mediaItems); i += batchSize {
		end := i + batchSize
		if end > len(mediaItems) {
			end = len(mediaItems)
		}
		batch := mediaItems[i:end]
		if err := db.Create(&batch).Error; err != nil {
			return err
		}

		log.Printf("Processed batch %d to %d\n", i, end)
	}
	return nil
}

func generateRandomTags(min, max int) []string {
	gofakeit.Seed(0)
	numTags := gofakeit.Number(min, max)
	tags := make([]string, numTags)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < numTags; i++ {
		tags[i] = fmt.Sprintf("%s%d", gofakeit.Word(), timestamp)
		timestamp++
	}
	return tags
}

func main() {
	tags := generateRandomTags(NumTags, NumTags)

	if err := batchSaveTags(config.DB, tags, BatchSize); err != nil {
		log.Fatalf("Failed to batch save tags: %v", err)
	}

	var mediaItems []models.Media
	for i := 0; i < NumMedia; i++ {
		mediaItems = append(mediaItems, models.Media{
			Name: fmt.Sprintf("Media %d", i+1),
			Link: fmt.Sprintf("https://example.com/media%d.jpg", i+1),
			Tags: []models.Tag{{Name: tags[i%len(tags)]}},
		})
	}

	if err := batchSaveMedia(config.DB, mediaItems, BatchSize); err != nil {
		log.Fatalf("Failed to batch save media items: %v", err)
	}

	fmt.Println("Successfully saved 1,000,000 tags and 1,000,000 media items.")
}
