package domain

import "context"

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	CreateOrderItem(ctx context.Context, orderItem *OrderItem) (*OrderItem, error)
	GetOrderById(ctx context.Context, orderId int) (*Order, error)
	GetOrderItems(ctx context.Context, orderId int) ([]*OrderItem, error)
	ListAllOrders(ctx context.Context, limit, offset int) ([]*Order, error)
	ListOrderByUserId(ctx context.Context, userId string, limit, offset int) ([]*Order, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}
