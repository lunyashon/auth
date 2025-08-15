package rabbit

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *Rabbit) DeclareQueue(ctx context.Context, queue string) (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *Rabbit) SendToEmailWithRetry(
	ctx context.Context,
	body []byte,
	queue amqp.Queue,
) {
	for i := 0; i < r.rbcfg.MaxRetries; i++ {
		if err := r.sendToEmail(ctx, body, queue); err != nil {
			r.log.Error("failed to send email", "error", err)
			time.Sleep(r.rbcfg.RetryDelay)
		} else {
			return
		}
	}
}

func (r *Rabbit) sendToEmail(
	ctx context.Context,
	body []byte,
	queue amqp.Queue,
) error {

	if err := r.IsConnected(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.ch.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	); err != nil {
		return err
	}

	return nil
}
