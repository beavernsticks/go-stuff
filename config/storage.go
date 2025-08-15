package bsgostuff_config

import bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"

type StorageTypeLocal struct {
	BasePath string `env:"INFRASTRUCTURE__STORAGE__LOCAL__BASE_PATH" env-default:"uploads"`
	BaseUrl  string `env:"INFRASTRUCTURE__STORAGE__BASE_URL" env-default:""`
}

type StorageTypeS3 struct {
	Endpoint     string `env:"INFRASTRUCTURE__STORAGE__S3__ENDPOINT"`
	Region       string `env:"INFRASTRUCTURE__STORAGE__S3__REGION"`
	PublicDomain string `env:"INFRASTRUCTURE__STORAGE__S3__PUBLIC_DOMAIN"`
	Bucket       string `env:"INFRASTRUCTURE__STORAGE__S3__BUCKET"`
	AccessKey    string `env:"INFRASTRUCTURE__STORAGE__S3__ACCESS_KEY"`
	SecretKey    string `env:"INFRASTRUCTURE__STORAGE__S3__SECRET_KEY"`
}

type Storage struct {
	Type  bsgostuff_domain.StorageTypeEnum `env:"INFRASTRUCTURE__STORAGE__TYPE" env-default:"LOCAL"`
	Local StorageTypeLocal
	S3    StorageTypeS3
}
