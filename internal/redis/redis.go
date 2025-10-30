package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/wb-go/wbf/redis"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func New(cfg config.RedisConfig) (*Redis, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client := redis.New(addr, cfg.Password, cfg.Database)
	rc := &Redis{
		client: client,
		ttl:    cfg.TTL,
	}
	return rc, nil
}

func (r *Redis) GetValue(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

func (r *Redis) SetValue(ctx context.Context, key string, value interface{}) error {
	return r.client.SetEX(ctx, key, value, r.ttl).Err()
}
