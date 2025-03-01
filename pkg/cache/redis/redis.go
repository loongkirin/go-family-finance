package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	RedisConfig *RedisConfig
	master      *redis.Client
	slaves      []*redis.Client
	lock        sync.RWMutex
	current     int
}

func NewRedisClient(cfg *RedisConfig) (*RedisClient, error) {
	master := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Master.Host, cfg.Master.Port),
		Password: cfg.Master.Password,
		DB:       cfg.Master.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := master.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to master redis: %w", err)
	}

	var slaves []*redis.Client
	for _, slaveCfg := range cfg.Slaves {
		slave := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", slaveCfg.Host, slaveCfg.Port),
			Password: slaveCfg.Password,
			DB:       slaveCfg.DB,
			PoolSize: cfg.PoolSize,
		})

		if err := slave.Ping(context.Background()).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to slave redis: %w", err)
		}
		slaves = append(slaves, slave)
	}

	return &RedisClient{
		RedisConfig: cfg,
		master:      master,
		slaves:      slaves,
	}, nil
}

func (r *RedisClient) GetMasterDb() *redis.Client {
	return r.master
}

func (r *RedisClient) GetSlaveDb() *redis.Client {
	if len(r.slaves) == 0 {
		return r.master
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.current = (r.current + 1) % len(r.slaves)
	return r.slaves[r.current]
}
