package kafka

import "context"

// EventConsumer - обрабатывает одно уже прочитанное и распарсенное событие.
// Реализация (например, простой логгер) вызывается из цикла чтения Consumer.Run.
type EventConsumer interface {
	Handle(ctx context.Context, event Event) error
}
