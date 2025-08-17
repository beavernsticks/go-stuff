package bsgostuff_config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

// LoadConfig загружает конфигурацию из переменных окружения в указанную структуру
// T - тип структуры конфигурации (должен быть указателем или структурой)
func LoadConfig[T any]() (*T, error) {
	var cfg T
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustLoadConfig загружает конфигурацию или паникует при ошибке
func MustLoadConfig[T any]() *T {
	cfg, err := LoadConfig[T]()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return cfg
}
