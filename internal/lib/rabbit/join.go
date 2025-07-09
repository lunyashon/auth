package rabbit

import (
	"fmt"
	"log/slog"

	"github.com/lunyashon/auth/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func New(log *slog.Logger, cfg *config.ConfigEnv, rbcfg config.ConfigRabbit) *RabbitService {
	provider := Rabbit{
		log: log,
		cfg: cfg,
		rbcfg: &RabbitConfig{
			MaxRetries: rbcfg.MaxRetries,
			RetryDelay: rbcfg.RetryDelay,
		},
	}
	return &RabbitService{
		Send:    &provider,
		Connect: &provider,
	}
}

func (r *Rabbit) Connect() error {
	var err error

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			r.cfg.RABBIT_NAME,
			r.cfg.RABBIT_PASSWORD,
			r.cfg.RABBIT_HOST,
			r.cfg.RABBIT_PORT,
		),
	)
	if err != nil {
		r.log.Error("failed to dial", "error", err)
		return err
	}
	r.conn = conn
	return nil
}

func (r *Rabbit) Channel() error {
	ch, err := r.conn.Channel()
	if err != nil {
		r.log.Error("failed to create channel", "error", err)
		return err
	}
	r.ch = ch
	return nil
}

func (r *Rabbit) CloseConnection() error {
	return r.conn.Close()
}

func (r *Rabbit) CloseChannel() error {
	return r.ch.Close()
}

func (r *Rabbit) IsConnected() error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("connection is not established")
	}
	if r.ch == nil || r.ch.IsClosed() {
		return fmt.Errorf("channel is not established")
	}
	return nil
}
