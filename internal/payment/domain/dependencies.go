package domain

import (
	order "CommerceCore/internal/order/domain"
	"context"
)

// OrderRepo - то, что payment-домену нужно от order для Create.
type OrderRepo interface {
	GetOrderById(ctx context.Context, orderId int) (*order.Order, error)
	// MarkOrderPaid - атомарный переход created -> paid. Недостаточно (заказ уже
	// не created) - errs.OrderNotPayable, вся транзакция Payment откатывается.
	MarkOrderPaid(ctx context.Context, orderId int) error
}

type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload any) error
}
