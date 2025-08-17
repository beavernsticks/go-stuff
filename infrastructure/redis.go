package bsgostuff_infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

// RedisAdapter - универсальный клиент для работы с Redis
type RedisAdapter struct {
	client *redis.Client
	prefix string
}

// New создает новый экземпляр адаптера
func NewRedisAdapter(cfg bsgostuff_config.Redis) (*RedisAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	var prefix string
	if cfg.Prefix != "" {
		prefix = cfg.Prefix + ":"
	}

	return &RedisAdapter{
		client: client,
		prefix: prefix, // "prod:user:123"
	}, nil
}

// MustNewRedisAdapter создает адаптер или паникует при ошибке
func MustNewRedisAdapter(cfg bsgostuff_config.Redis) *RedisAdapter {
	adapter, err := NewRedisAdapter(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize Redis adapter: %w", err))
	}
	return adapter
}

// Close закрывает соединение
func (a *RedisAdapter) Close() error {
	return a.client.Close()
}

// GetProto получает и десериализует protobuf-сообщение
func (a *RedisAdapter) GetProto(ctx context.Context, key string, msg proto.Message) error {
	fullKey := a.prefix + key

	data, err := a.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return bsgostuff_domain.ErrNotFound
		}
		return err
	}

	return proto.Unmarshal(data, msg)
}

// SetProto сериализует и сохраняет protobuf-сообщение
func (a *RedisAdapter) SetProto(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error {
	fullKey := a.prefix + key

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return a.client.Set(ctx, fullKey, data, ttl).Err()
}

// GetJSON получает и десериализует JSON
func (a *RedisAdapter) GetJSON(ctx context.Context, key string, dest interface{}) error {
	fullKey := a.prefix + key

	data, err := a.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return bsgostuff_domain.ErrNotFound
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// SetJSON сериализует и сохраняет JSON
func (a *RedisAdapter) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := a.prefix + key

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.client.Set(ctx, fullKey, data, ttl).Err()
}

// Delete удаляет ключ
func (a *RedisAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, a.prefix+key).Err()
}

// WithPrefix создает новый экземпляр с доп. префиксом
func (a *RedisAdapter) WithPrefix(prefix string) *RedisAdapter {
	return &RedisAdapter{
		client: a.client,
		prefix: a.prefix + prefix + ":",
	}
}
