package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type ProducerKafka struct {
	k *kafka.Writer
}

func NewProducer(k *kafka.Writer) *ProducerKafka {
	return &ProducerKafka{k: k}
}

// Publish - метод продюсер, создает конверт для получения данных из Order или payment,
// чтобы не создавать каждый раз заново, конвертирует полученный конверт в JSON,
// отправляет и сохраняет сообщение в топике
func (p *ProducerKafka) Publish(ctx context.Context, eventType string, payload any) error {
	event, err := NewEvent(eventType, payload)
	if err != nil {
		slog.Error("failed to create event", "error", err)
		return err
	}

	b, err := json.Marshal(&event)
	if err != nil {
		slog.Error("failed to marshal event", "error", err)
		return err
	}
	return p.k.WriteMessages(ctx, kafka.Message{Value: b})
}
