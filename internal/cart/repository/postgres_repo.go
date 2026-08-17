package repository

import (
	"CommerceCore/internal/cart/domain"
	"CommerceCore/pkg/querier"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type CartRepoImpl struct {
	q querier.Querier
}

func NewCartRepo(q querier.Querier) *CartRepoImpl {
	return &CartRepoImpl{q: q}
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
