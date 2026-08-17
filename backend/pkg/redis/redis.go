package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/bees/hindu-ritual-platform/pkg/configs"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func InitRedis(cfg *configs.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRedis(client *redis.Client) error {
	return client.Close()
}

func (r *RedisClient) AddToBlacklist(token string, ttl time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, true, ttl).Err()
}

func (r *RedisClient) IsBlacklisted(token string) (bool, error) {
	exists, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	return exists > 0, err
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Delete(key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
