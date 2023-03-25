package redis

import (
	"context"
	"errors"
	redis "github.com/redis/go-redis/v9"
	"time"
)

type RedisObject struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

type RedisImpl interface {
	SetDataToRedis(key, value string) error
	GetDataFromRedis(string) (string, error)
	DeleteDataFromRedis(string) error
}

func NewRedisClient(redis *redis.Client, ctx context.Context) RedisImpl {
	return &RedisObject{
		RedisClient: redis,
		Ctx:         ctx,
	}
}

func (r RedisObject) SetDataToRedis(key, value string) error {
	err := r.RedisClient.Set(r.Ctx, key, value, time.Hour).Err()

	if err != nil {
		return errors.New("redis Set Error")
	}

	return nil
}
func (r *RedisObject) GetDataFromRedis(key string) (string, error) {

	val, err := r.RedisClient.Get(r.Ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisObject) DeleteDataFromRedis(key string) error {
	_, err := r.RedisClient.Del(r.Ctx, key).Result()

	if err != nil {
		return err
	}

	return nil
}
