package main

import (
	"fmt"

	"media-indexer/config"
	"media-indexer/models"
)

func main() {
	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.Media{}, &models.Tag{}, &models.MediaTag{}); err != nil {
		fmt.Printf("Error during migration: %v\n", err)
		return
	}

	fmt.Println("Migration has finished successfully!")
}
