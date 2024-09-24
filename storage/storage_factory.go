package storage

import (
	"os"
)

func StorageFactory() StorageProvider {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "s3" {
		return NewS3Storage()
	} else if storageType == "local" {
		return NewAnotherStorage()
	}
	return NewAnotherStorage()
}
