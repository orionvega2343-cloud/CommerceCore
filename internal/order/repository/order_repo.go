package repository

import (
	"CommerceCore/internal/order/domain"
	"CommerceCore/internal/order/domain/errs"
	"CommerceCore/pkg/querier"
	"CommerceCore/pkg/transaction"
	"context"
	"database/sql"
	"log/slog"
)

type OrderRepoImpl struct {
	q querier.Querier
}

func NewOrderRepo(q querier.Querier) *OrderRepoImpl {
	return &OrderRepoImpl{q: q}
}

func (r *OrderRepoImpl) CreateOrder(ctx context.Context, m *domain.Order) (*domain.Order, error) {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	err := q.GetContext(ctx, &m, `INSERT INTO orders(user_id, cart_id, status, total_amount) VALUES($1, $2, $3, $4) RETURNING id`, m.UserId, m.CartId, m.Status, m.TotalAmount)
	if err != nil {
		slog.Error("failed to insert order", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *OrderRepoImpl) CreateOrderItem(ctx context.Context, m *domain.OrderItem) (*domain.OrderItem, error) {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	err := q.GetContext(ctx, &m, `INSERT INTO order_items(order_id, product_id, quantity, price_per_unit) VALUES($1, $2, $3, $4) RETURNING id`, m.OrderId, m.ProductId, m.Quantity, m.PricePerUnit)
	if err != nil {
		slog.Error("failed to insert order_item", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *OrderRepoImpl) GetOrderById(ctx context.Context, orderId int) (*domain.Order, error) {
	var m domain.Order
	err := r.q.GetContext(ctx, &m, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id = $1`, orderId)
	if err != nil {
		slog.Error("failed to get order", "error", err)
		return nil, err
	}
	return &m, nil
}

func (r *OrderRepoImpl) GetOrderItems(ctx context.Context, orderId int) ([]*domain.OrderItem, error) {
	var m []*domain.OrderItem
	err := r.q.SelectContext(ctx, &m, `SELECT id, order_id, product_id, quantity, price_per_unit FROM order_items WHERE order_id = $1`, orderId)
	if err != nil {
		slog.Error("failed to get order_items", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *OrderRepoImpl) ListAllOrders(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	var m []*domain.Order
	err := r.q.SelectContext(ctx, &m, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		slog.Error("failed to get all orders", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *OrderRepoImpl) ListOrderByUserId(ctx context.Context, userId string, limit, offset int) ([]*domain.Order, error) {
	var m []*domain.Order
	err := r.q.SelectContext(ctx, &m, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE user_id= $1 LIMIT $2 OFFSET $3`, userId, limit, offset)
	if err != nil {
		slog.Error("failed to get all orders", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *OrderRepoImpl) UpdateStatus(ctx context.Context, id int, status string) error {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	res, err := q.ExecContext(ctx, `UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		slog.Error("failed to update order", "error", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "error", err)
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// MarkOrderPaid - получает доступ к транзакции из домена Payment, через ctx
// делает атомарный SQL запрос, чтобы сменить статус без race condition,
// проверяет количество изменившихся строк, если их 0, выбрасывает sentinel error
func (r *OrderRepoImpl) MarkOrderPaid(ctx context.Context, orderId int) error {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	res, err := q.ExecContext(ctx, `UPDATE orders SET status = 'paid' WHERE id = $1 AND status = 'created'`, orderId)
	if err != nil {
		slog.Error("failed to update order", "error", err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("failed to update order", "error", err)
		return err
	}
	if affected == 0 {
		return errs.OrderNotPayable
	}
	return nil
}
