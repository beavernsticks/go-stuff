package bsgostuff_infrastructure

import (
	"context"
	"errors"
	"fmt"
	"io"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
)

type IStorage interface {
	Close() error
	GetType() bsgostuff_domain.StorageTypeEnum
	GetBaseUrl() string
	GetFullPath(path string) string
	Upload(ctx context.Context, path string, file io.Reader, size int64, contentType string) (string, error)
	Download(ctx context.Context, path string) (io.ReadCloser, error)
}

func NewStorage(config bsgostuff_config.Storage) (IStorage, error) {
	switch config.Type {
	case bsgostuff_domain.StorageTypeEnumLocal:
		return newLocalStorage(config.Local)
	case bsgostuff_domain.StorageTypeEnumS3:
		return newS3Storage(config.S3)
	default:
		return nil, errors.New("unsupported storage type")
	}
}

// MustNewStorage создает адаптер или паникует при ошибке
func MustNewStorage(cfg bsgostuff_config.Storage) IStorage {
	storage, err := NewStorage(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize Storage: %w", err))
	}
	return storage
}
