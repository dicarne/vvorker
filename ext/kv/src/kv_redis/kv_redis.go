package kvredis

import (
	"context"
	"fmt"
	"time"
	"vvorker/conf"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type KVRedis struct {
}

var rdb *redis.Client

func init() {
	if conf.AppConfigInstance.KVProvider == "redis" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%d", conf.AppConfigInstance.ClientRedisPort),
			Password: conf.AppConfigInstance.ServerRedisPassword,
			DB:       0,
		})

		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			logrus.Error("[ERROR] redis ping failed")
		}
	}
}

func (r *KVRedis) Close() {
	if rdb != nil {
		rdb.Close()
	}
}

func (r *KVRedis) Put(bucket string, key string, value []byte, ttl int) (int, error) {
	code := 0
	_, err := rdb.Set(context.Background(), bucket+":"+key, value, time.Second*time.Duration(ttl)).Result()
	if err != nil {
		code = 1
		logrus.Error(err)
	}
	return code, err
}

func (r *KVRedis) PutNX(bucket string, key string, value []byte, ttl int) (int, error) {
	code := 0
	_, err := rdb.SetNX(context.Background(), bucket+":"+key, value, time.Second*time.Duration(ttl)).Result()
	if err != nil {
		code = 1
		logrus.Error(err)
	}
	return code, err
}

func (r *KVRedis) PutXX(bucket string, key string, value []byte, ttl int) (int, error) {
	code := 0
	_, err := rdb.Set(context.Background(), bucket+":"+key, value, time.Second*time.Duration(ttl)).Result()
	if err != nil {
		code = 1
		logrus.Error(err)
	}
	return code, err
}

func (r *KVRedis) Get(bucket string, key string) ([]byte, error) {
	v, err := rdb.Get(context.Background(), bucket+":"+key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(v), nil
}

func (r *KVRedis) Del(bucket string, key string) error {
	return rdb.Del(context.Background(), bucket+":"+key).Err()
}

func (r *KVRedis) Keys(bucket string, prefix string, offset int, size int) ([]string, error) {
	result, err := rdb.Keys(context.Background(), bucket+":"+prefix).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
