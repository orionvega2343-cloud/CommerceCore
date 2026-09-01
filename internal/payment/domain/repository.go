package domain

import "context"

type PaymentRepo interface {
	CreatePayment(ctx context.Context, d *Payment) (*Payment, error)
	GetPaymentById(ctx context.Context, id int) (*Payment, error)
	ListAllPayments(ctx context.Context, limit, offset int) ([]*Payment, error)
	ListByUserId(ctx context.Context, userId string, limit, offset int) ([]*Payment, error)
	TotalAmountAll(ctx context.Context) (int, error)
	TotalAmountByUserId(ctx context.Context, userId string) (int, error)
}
