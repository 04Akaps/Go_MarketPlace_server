package redis

import (
	"context"
	"encoding/json"
	"errors"
	redis "github.com/redis/go-redis/v9"
	"log"
	"net"
	"time"
)

type RedisObject struct {
	RedisClient *redis.Client
	Ctx         context.Context
	Conn        *redis.Conn
	ErrorLog    RedisErrorLog
}

type RedisErrorLog struct {
	RedisErrorChannel chan error
	Logger            *log.Logger
}

type RedisImpl interface {
	SetDataToRedis(string, interface{}) error
	GetDataFromRedis(string) ([]byte, error)
	DeleteDataFromRedis(string) error
	SetRedisConn()
}

func (r *RedisObject) RedisErrorChannelInit() {
	go func() {
		for {
			select {
			case redisErr := <-r.ErrorLog.RedisErrorChannel:
				log.Println("Error : ", redisErr)
				r.ErrorLog.Logger.Println(redisErr)
			}
		}
	}()
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

func (r RedisObject) SetDataToRedis(key string, value interface{}) error {
	if r.Conn == nil {
		r.SetRedisConn()
	}

	byteData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = r.Conn.Set(r.Ctx, key, string(byteData), time.Hour).Err()
	if err != nil {
		return errors.New("redis Set Error")
	}

	return nil
}

func (r *RedisObject) GetDataFromRedis(key string) ([]byte, error) {
	if r.Conn == nil {
		r.SetRedisConn()
	}

	val, err := r.Conn.Get(r.Ctx, key).Bytes()

	err = r.redisErrorHandler(func() ([]byte, error) { return r.GetDataFromRedis(key) }, err)

	if err != nil {
		return nil, err
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

func (r *RedisObject) redisErrorHandler(f func() ([]byte, error), err error) error {
	if err.Error() == "redis: client is closed" {
		r.SetRedisConn()
		f()
		return nil
	}

	if err == redis.TxFailedErr {
		r.ErrorLog.RedisErrorChannel <- err
	}

	if err == redis.Nil {
		r.ErrorLog.RedisErrorChannel <- err
	}

	return err
}
