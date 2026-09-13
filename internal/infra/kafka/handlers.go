package kafka

import (
	"context"
	"log/slog"
)

func (l *EventLogger) Handle(ctx context.Context, event Event) error {
	slog.Info("logging of data read by the consumer",
		slog.String("event_type", event.EventType),
		slog.String("payload", string(event.Payload)),
		slog.Time("timestamp", event.Timestamp),
	)
	return nil
}
