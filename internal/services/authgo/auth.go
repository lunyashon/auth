package authgo

import (
	"context"
	"log/slog"

	"github.com/lunyashon/auth/internal/config"
	database "github.com/lunyashon/auth/internal/database/psql"
	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"github.com/lunyashon/auth/internal/lib/rabbit"
	"github.com/lunyashon/auth/internal/lib/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

// Интерфейс для реализации методов
type Auth interface {
	LoginUser(
		ctx context.Context,
		data *sso.LoginRequest,
	) (token string, err error)
	RegisterUser(
		ctx context.Context,
		data *sso.RegisterRequest,
	) (userId int64, err error)
	RegisterToken(
		ctx context.Context,
		data *sso.TokenRequest,
	) (result bool, err error)
}

type AuthData struct {
	Auth              Auth
	DB                *database.StructDatabase
	Log               *slog.Logger
	Cfg               *config.ConfigEnv
	Rabbit            *rabbit.RabbitService
	QueueForgotToken  amqp.Queue
	QueueConfirmEmail amqp.Queue
	KeysStore         *jwtsso.KeysStore
	Yaml              *config.ConfigYaml
	Redis             *redis.Redis
}

var errors = map[codes.Code]struct {
	message string
}{
	codes.InvalidArgument:  {},
	codes.PermissionDenied: {},
	codes.Internal:         {},
	codes.NotFound:         {},
	codes.Unauthenticated:  {},
}
