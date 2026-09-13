package repository

import (
	"CommerceCore/internal/cart/domain"
	"CommerceCore/pkg/querier"
	"CommerceCore/pkg/transaction"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type CartRepoImpl struct {
	q     querier.Querier
	items *CartItemRepoImpl
}

func NewCartRepo(q querier.Querier, items *CartItemRepoImpl) *CartRepoImpl {
	return &CartRepoImpl{q: q, items: items}
}

func (r *CartRepoImpl) CreateOrGet(ctx context.Context, userId string) (*domain.Cart, error) {
	var c domain.Cart
	err := r.q.GetContext(ctx, &c, `SELECT id, user_id, status, created_at, updated_at FROM cart WHERE user_id = $1 AND status = 'active'`, userId)
	if errors.Is(err, sql.ErrNoRows) {
		err = r.q.GetContext(ctx, &c, `INSERT INTO cart(user_id, status) VALUES($1, $2) RETURNING id`, userId, "active")
		if err != nil {
			slog.Error("failed to insert cart into database", "error", err)
			return nil, err
		}
	}
	if err != nil {
		slog.Error("failed to get cart into database", "error", err)
		return nil, err
	}
	return &c, nil
}

func (r *CartRepoImpl) LoadItems(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	return r.items.LoadItems(ctx, cartID)
}

func (r *CartRepoImpl) SetStatus(ctx context.Context, cartID int, status string) error {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	if _, err := q.ExecContext(ctx, `UPDATE cart SET status = $1 WHERE id = $2`, status, cartID); err != nil {
		slog.Error("failed to set cart status", "error", err)
		return err
	}
	return nil
}
