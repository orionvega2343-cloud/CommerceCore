package repository

import (
	"CommerceCore/internal/payment/domain"
	"CommerceCore/pkg/querier"
	"CommerceCore/pkg/transaction"
	"context"
	"log/slog"
)

var _ domain.PaymentRepo = (*PaymentRepoImpl)(nil)

type PaymentRepoImpl struct {
	q querier.Querier
}

func NewPaymentRepo(q querier.Querier) *PaymentRepoImpl {
	return &PaymentRepoImpl{q: q}
}

func (r *PaymentRepoImpl) CreatePayment(ctx context.Context, d *domain.Payment) (*domain.Payment, error) {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}

	if err := q.GetContext(ctx, &d, `INSERT INTO payments(order_id, amount, status, method) VALUES($1, $2, $3, $4) RETURNING id`, d.OrderId, d.Amount, d.Status, d.Method); err != nil {
		slog.Error("failed to insert payments", "error", err)
		return nil, err
	}

	return d, nil
}

func (r *PaymentRepoImpl) GetPaymentById(ctx context.Context, id int) (*domain.Payment, error) {
	var d domain.Payment
	if err := r.q.GetContext(ctx, &d, `SELECT id, order_id, amount, status, method FROM payments WHERE id =$1`, id); err != nil {
		slog.Error("failed to get payment by id", "id", id, "error", err)
		return nil, err
	}
	return &d, nil
}

func (r *PaymentRepoImpl) ListAllPayments(ctx context.Context, limit, offset int) ([]*domain.Payment, error) {
	var d []*domain.Payment
	if err := r.q.SelectContext(ctx, &d, `SELECT id, order_id, amount, status, method, created_at FROM payments LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		slog.Error("failed to list payments", "error", err)
		return nil, err
	}
	return d, nil
}

func (r *PaymentRepoImpl) ListByUserId(ctx context.Context, userId string, limit, offset int) ([]*domain.Payment, error) {
	var d []*domain.Payment
	query := `SELECT p.id, p.order_id, p.amount, p.status, p.method, p.created_at FROM payments p JOIN orders o ON p.order_id = o.id WHERE o.user_id = $1 LIMIT $2 OFFSET $3`
	if err := r.q.SelectContext(ctx, &d, query, userId, limit, offset); err != nil {
		slog.Error("failed to list payments", "error", err)
		return nil, err
	}
	return d, nil
}

func (r *PaymentRepoImpl) TotalAmountAll(ctx context.Context) (int, error) {
	var total int
	if err := r.q.GetContext(ctx, &total, `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'succeeded'`); err != nil {
		slog.Error("failed to get total amount", "error", err)
		return 0, err
	}
	return total, nil
}

func (r *PaymentRepoImpl) TotalAmountByUserId(ctx context.Context, userId string) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(p.amount), 0) FROM payments p JOIN orders o ON p.order_id = o.id WHERE o.user_id = $1 AND p.status = 'succeeded'`
	if err := r.q.GetContext(ctx, &total, query, userId); err != nil {
		slog.Error("failed to get total amount", "error", err)
		return 0, err
	}
	return total, nil
}
