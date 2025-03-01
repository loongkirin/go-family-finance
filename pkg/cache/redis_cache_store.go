package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	redisClient *redis.Client
	Prekey      string
	Expiration  time.Duration
	Context     context.Context
}

func NewRedisStore(redisClient *redis.Client, prekey string, defaultExpiration time.Duration) *RedisStore {
	return &RedisStore{
		redisClient: redisClient,
		Prekey:      prekey,
		Expiration:  defaultExpiration,
		Context:     context.Background(),
	}
}

func (rs *RedisStore) UseWithContext(ctx context.Context) *RedisStore {
	rs.Context = ctx
	return rs
}

func (rs *RedisStore) Set(key string, value interface{}, expires time.Duration) error {
	err := rs.redisClient.Set(rs.Context, rs.Prekey+key, value, expires).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (rs *RedisStore) Add(key string, value interface{}, expires time.Duration) error {
	_, err := rs.redisClient.Get(rs.Context, rs.Prekey+key).Result()
	if err == redis.Nil {
		fmt.Println(key, " does not exist")
		return ErrNotStored
	}

	return rs.Set(key, value, expires)
}

func (rs *RedisStore) Replace(key string, value interface{}, expires time.Duration) error {
	return rs.Set(key, value, expires)
}

func (rs *RedisStore) Get(key string) (string, error) {
	value, err := rs.redisClient.Get(rs.Context, rs.Prekey+key).Result()
	if err == redis.Nil {
		fmt.Println(key, " does not exist")
		return "", ErrNotStored
	}
	return value, nil
}

func (rs *RedisStore) Delete(key string) error {
	err := rs.redisClient.Del(rs.Context, rs.Prekey+key).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (rs *RedisStore) Increment(key string, value int64) (int64, error) {
	newValue, err := rs.redisClient.IncrBy(rs.Context, rs.Prekey+key, value).Result()
	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func (rs *RedisStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := rs.redisClient.DecrBy(rs.Context, rs.Prekey+key, value).Result()
	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func (rs *RedisStore) Flush() error {
	err := rs.redisClient.FlushAll(rs.Context).Err()
	return err
}
