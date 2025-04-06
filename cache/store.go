package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	ClearAll(ctx context.Context) error
}

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) Store {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return s.rdb.Get(ctx, key).Result()
}

func (s *RedisStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.rdb.Del(ctx, key).Err()
}

func (s *RedisStore) ClearAll(ctx context.Context) error {
	return s.rdb.FlushDB(ctx).Err()
}
