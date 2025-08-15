package rabbit

import (
	"context"
	"log/slog"
	"time"

	"github.com/lunyashon/auth/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitService struct {
	Send    RabbitProvider
	Connect ConnectProvider
}

type Rabbit struct {
	log   *slog.Logger
	cfg   *config.ConfigEnv
	conn  *amqp.Connection
	ch    *amqp.Channel
	rbcfg *RabbitConfig
}

type ConnectProvider interface {
	Connect() error
	Channel() error
	CloseConnection() error
	CloseChannel() error
	IsConnected() error
}

type RabbitProvider interface {
	DeclareQueue(
		ctx context.Context,
		queue string,
	) (amqp.Queue, error)
	SendToEmailWithRetry(
		ctx context.Context,
		body []byte,
		queue amqp.Queue,
	)
}

type RabbitConfig struct {
	MaxRetries int           `yaml:"max_retries" env-default:"5"`
	RetryDelay time.Duration `yaml:"retry_delay" env-default:"1s"`
}
