package storage

import (
	"crypto/sha256"
	"fmt"
	"mime/multipart"
)

type S3Storage struct{}

func NewS3Storage() *S3Storage {
	return &S3Storage{}
}

func (s *S3Storage) UploadFile(_fileHeader *multipart.FileHeader, filename string) (string, error) {
	hash := sha256.Sum256([]byte(filename))
	hashString := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("https://s3.amazonaws.com/bucket/%s.jpg", hashString), nil
}
