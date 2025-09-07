package redis

import (
	"context"
	"time"
)

func (r *RedisProvider) CheckTrye(ctx context.Context, login string) (int, error) {
	count, err := r.Client.Get(ctx, login).Int()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisProvider) AddToTrye(ctx context.Context, login string, count int) error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return err
	}
	if err := r.Client.Set(ctx, login, count, 5*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}
