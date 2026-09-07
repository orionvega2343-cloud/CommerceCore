package domain

import "context"

type PaymentService interface {
	// Create - idempotencyKey приходит от клиента (например, из заголовка
	// Idempotency-Key); PaymentServiceImpl его игнорирует, IdempotencyServiceImpl
	// использует для защиты от дублей.
	Create(ctx context.Context, payment *Payment, idempotencyKey string) (*Payment, error)
	GetPaymentById(ctx context.Context, paymentId int, userId, role string) (*Payment, error)
	ListPayments(ctx context.Context, role, userId string, limit, offset int) ([]*Payment, error)
	TotalAmountPayments(ctx context.Context, role, userId string) (int, error)
}
