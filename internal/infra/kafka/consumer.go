package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler EventConsumer
}

func NewConsumer(reader *kafka.Reader, handler EventConsumer) *Consumer {
	return &Consumer{reader: reader, handler: handler}
}

// Run - бесконечно читает сообщения из топика, пока не отменят ctx (graceful shutdown)
// или пока чтение не вернёт неустранимую ошибку. На каждое сообщение: распаковывает JSON
// в Event и передаёт его в handler.Handle.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("failed to read message", "error", err)
			break
		}

		var event Event
		if err = json.Unmarshal(m.Value, &event); err != nil {
			slog.Error("failed to unmarshal event", "error", err)
			continue
		}

		if err = c.handler.Handle(ctx, event); err != nil {
			slog.Error("failed to handle event", "error", err)
		}
	}

	if err := c.reader.Close(); err != nil {
		slog.Error("failed to close reader", "error", err)
		return err
	}
	return nil
}
