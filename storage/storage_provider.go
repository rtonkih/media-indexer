package storage

import "mime/multipart"

type StorageProvider interface {
	UploadFile(fileHeader *multipart.FileHeader, filename string) (string, error)
}
