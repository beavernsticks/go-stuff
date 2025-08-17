package bsgostuff_infrastructure

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
)

type localStorage struct {
	basePath string
	baseUrl  string
	mu       sync.RWMutex
}

func newLocalStorage(config bsgostuff_config.StorageTypeLocal) (IStorage, error) {
	// Создаем базовую директорию, если ее нет
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, err
	}

	return &localStorage{
		basePath: config.BasePath,
		baseUrl:  config.BaseUrl,
	}, nil
}

func (s *localStorage) Upload(ctx context.Context, path string, file io.Reader, size int64, contentType string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Создаем полный путь к файлу
	fullPath := filepath.Join(s.basePath, path)

	// Создаем директории, если их нет
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	// Создаем файл
	slog.Info("Upload", "create", fullPath)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Копируем данные
	if _, err := io.Copy(out, file); err != nil {
		slog.Error("Upload error", "err", err)
		// Удаляем частично записанный файл при ошибке
		os.Remove(fullPath)
		return "", err
	}

	slog.Info("Upload", "complete", fullPath)
	return fullPath, nil
}

// Download возвращает reader для чтения файла
func (s *localStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

// Close освобождает ресурсы хранилища
func (s *localStorage) Close() error {
	// Для локального хранилища не требуется освобождение ресурсов
	return nil
}

func (s *localStorage) GetType() bsgostuff_domain.StorageTypeEnum {
	return bsgostuff_domain.StorageTypeEnumLocal
}

func (s *localStorage) GetBaseUrl() string {
	return s.baseUrl
}

func (s *localStorage) GetFullPath(path string) string {
	return filepath.Join(s.basePath, path)
}
