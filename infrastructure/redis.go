package shared_infrastructure

import (
	"context"
	"time"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) Raw() *redis.Client {
	return r.client
}

func NewRedisClient(config bsgostuff_config.Redis) *RedisClient {
	redisHost := config.Host

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: config.Password,
		DB:       config.Database,
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}
