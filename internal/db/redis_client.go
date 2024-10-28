package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})
	return &RedisClient{Client: rdb}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return r.Client.SIsMember(ctx, key, member)
}

func (r *RedisClient) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.Client.SAdd(ctx, key, members...)
}
