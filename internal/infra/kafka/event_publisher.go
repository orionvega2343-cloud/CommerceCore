package kafka

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload any) error
}
