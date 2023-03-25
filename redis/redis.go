package redis

import (
	"context"
	"errors"
	redis "github.com/redis/go-redis/v9"
	"net"
	"time"
)

type RedisObject struct {
	RedisClient *redis.Client
	Ctx         context.Context
	Conn        *redis.Conn
}

type RedisImpl interface {
	SetDataToRedis(key, value string) error
	GetDataFromRedis(string) (string, error)
	DeleteDataFromRedis(string) error
	SetRedisConn()
}

func NewRedisClient(option *redis.Options, ctx context.Context) RedisImpl {

	option.PoolFIFO = true

	option.Dialer = func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, 5*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	option.Addr = "redis-11586.c289.us-west-1-2.ec2.cloud.redislabs.com:11586"

	//os.Getenv("redis_addr")
	option.CredentialsProvider = func() (string, string) {

		return "default", "qVp8Gy97gyTJWbK8VMueQi7xEw0fqgmR"
		//return os.Getenv("redis_user_name"), os.Getenv("redis_password")
	}

	return &RedisObject{
		RedisClient: redis.NewClient(option),
		Ctx:         ctx,
	}
}

func (r *RedisObject) SetRedisConn() {
	r.Conn = r.RedisClient.Conn()
}

func (r RedisObject) SetDataToRedis(key, value string) error {
	if r.Conn == nil {
		r.SetRedisConn()
	}

	err := r.Conn.Set(r.Ctx, key, value, time.Hour).Err()
	if err != nil {

		return errors.New("redis Set Error")
	}

	return nil
}

func (r *RedisObject) GetDataFromRedis(key string) (string, error) {
	if r.Conn == nil {
		r.SetRedisConn()
	}

	val, err := r.Conn.Get(r.Ctx, key).Result()

	if err != nil {
		if err == redis.TxFailedErr {
			// Tx가 실패한 경우 입니다.
			return "", err
		}

		if err == redis.Nil {
			// Key 값이 존재 하지 않을 떄
		}

		if err.Error() == "redis: client is closed" {
			//Client가 닫혀 있을 떄
		}

		return "", err
	}

	return val, nil
}

func (r *RedisObject) DeleteDataFromRedis(key string) error {
	if r.Conn == nil {
		r.SetRedisConn()
	}

	_, err := r.Conn.Del(r.Ctx, key).Result()

	if err != nil {
		return err
	}

	return nil
}
