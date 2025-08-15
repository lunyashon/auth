package redis

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lunyashon/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(config *config.ConfigEnv, log *slog.Logger) *Redis {
	var (
		once  sync.Once
		redis *Redis
		err   error
	)

	once.Do(func() {
		redisProvider := &RedisProvider{Config: config, Log: log}
		if err = redisProvider.GetClient(); err != nil {
			panic(err)
		}
		redis = &Redis{
			TokenProvider: redisProvider,
			Connect:       redisProvider,
		}
	})

	return redis
}

func (r *RedisProvider) GetClient() error {

	var err error

	for i := 0; i < 5; i++ {
		r.Client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", r.Config.REDIS_HOST, r.Config.REDIS_PORT),
			Password: r.Config.REDIS_PASSWORD,
			DB:       r.Config.REDIS_NUM_DB,
		})

		if err = r.Client.Ping(context.Background()).Err(); err != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		break
	}

	return err
}

func (r *RedisProvider) CloseClient() error {
	return r.Client.Close()
}

func (r *RedisProvider) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
