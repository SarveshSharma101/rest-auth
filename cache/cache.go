package cache

import (
	"context"
	"fmt"
	"rest-auth/utils"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var CacheClient *redis.Client

func InitCacheClient() error {
	db, err := strconv.Atoi(utils.Getenv(utils.CACHE_DB))
	if err != nil {
		return err
	}
	CacheClient = redis.NewClient(
		&redis.Options{
			Addr:     utils.Getenv(utils.CACHE_ADDR),
			Password: utils.Getenv(utils.CACHE_PASSWORD),
			DB:       db,
		},
	)
	pong, err := CacheClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	fmt.Println("Connected to cache: ", pong)
	return nil
}

func GetValues(key string) (string, error) {
	return CacheClient.Get(context.TODO(), key).Result()
}

func SetValues(key string, value string, expiration time.Duration) error {
	return CacheClient.Set(context.TODO(), key, value, expiration).Err()
}

func DeleteValues(key string) error {
	return CacheClient.Del(context.TODO(), key).Err()
}
