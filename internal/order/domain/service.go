package domain

import "context"

type OrderService interface {
	// Checkout - оформление заказа из активной корзины (Create по ТЗ идёт только отсюда).
	Checkout(ctx context.Context, userId string) (*Order, error)
	// GetOrder - один заказ по его id; admin видит любой, юзер — только свой.
	GetOrder(ctx context.Context, orderId int, userId, role string) (*Order, error)
	// ListOrders - admin получает все заказы, юзер — только свои.
	ListOrders(ctx context.Context, role, userId string, limit, offset int) ([]*Order, error)
	TransitionStatus(ctx context.Context, orderId int, next, userId, role string) (*Order, error)
}
