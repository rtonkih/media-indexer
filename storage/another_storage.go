package storage

import (
	"crypto/sha256"
	"fmt"
	"mime/multipart"
)

type AnotherStorage struct{}

func NewAnotherStorage() *AnotherStorage {
	return &AnotherStorage{}
}

func (s *AnotherStorage) UploadFile(_fileHeader *multipart.FileHeader, filename string) (string, error) {
	hash := sha256.Sum256([]byte(filename))
	hashString := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("https://another_storage.com/%s.jpg", hashString), nil
}
