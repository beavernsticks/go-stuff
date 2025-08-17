package bsgostuff_infrastructure

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
)

type s3Storage struct {
	client       *s3.Client
	bucket       string
	publicDomain string
}

// NewS3Storage создает новое S3 хранилище
func newS3Storage(config bsgostuff_config.StorageTypeS3) (IStorage, error) {
	client := s3.New(s3.Options{
		Region:       config.Region,
		Credentials:  credentials.NewStaticCredentialsProvider(config.AccessKey, config.SecretKey, "aws_session"),
		BaseEndpoint: aws.String(config.Endpoint),
	})

	return &s3Storage{
		client:       client,
		bucket:       config.Bucket,
		publicDomain: config.PublicDomain,
	}, nil
}

// Upload загружает файл в S3 хранилище
func (s *s3Storage) Upload(ctx context.Context, path string, file io.Reader, size int64, contentType string) (string, error) {
	// Создаём временный файл
	tmpFile, err := os.CreateTemp("", "s3upload-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name()) // Удаляем временный файл после завершения
	defer tmpFile.Close()

	// Копируем данные из pipe во временный файл
	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", err
	}

	// Перемещаем указатель в начало файла
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return "", err
	}

	// Получаем фактический размер файла
	fileInfo, err := tmpFile.Stat()
	if err != nil {
		return "", err
	}

	// Создаем входные параметры для загрузки
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(path),
		Body:          tmpFile,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(contentType),
	}

	// Выполняем загрузку
	_, err = s.client.PutObject(ctx, input)

	return path, err
}

// Download возвращает reader для чтения файла из S3
func (s *s3Storage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	// Создаем входные параметры для скачивания
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	// Выполняем запрос
	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// Close освобождает ресурсы S3 клиента
func (s *s3Storage) Close() error {
	// В текущей реализации AWS SDK не требует явного закрытия клиента
	return nil
}

func (s *s3Storage) GetType() bsgostuff_domain.StorageTypeEnum {
	return bsgostuff_domain.StorageTypeEnumS3
}

func (s *s3Storage) GetBaseUrl() string {
	return s.publicDomain
}

func (s *s3Storage) GetFullPath(path string) string {
	return path
}
