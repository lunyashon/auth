package redis

import (
	"context"
	"log/slog"

	"github.com/lunyashon/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	TokenProvider Token
	Connect       Connect
}

type RedisProvider struct {
	Client *redis.Client
	Config *config.ConfigEnv
	Log    *slog.Logger
}

type Connect interface {
	GetClient() error
	CloseClient() error
	Ping(ctx context.Context) error
}

type Token interface {
	AddToBlackList(
		ctx context.Context,
		userID int,
		token string,
		multitype string,
	) error
	CheckFromBlackList(
		ctx context.Context,
		userID int,
		expiredAt string,
		token string,
	) error
}
