package helper

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/yuugure-aikouka/kyoto-common/handler"
)

func GetMagicLinkKeysByEmail(rdb *redis.Client, email string) ([]string, error) {
	ctx := context.Background()
	var result []string

	magicLinkPattern := fmt.Sprintf("%s:*", handler.MagicLinkCacheKeyPrefix)
	iter := rdb.Scan(ctx, 0, magicLinkPattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := rdb.Get(ctx, key).Result()
		if err == nil && val == email {
			result = append(result, key)
		}
	}

	return result, iter.Err()
}
