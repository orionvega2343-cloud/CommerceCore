package repository

import (
	"CommerceCore/internal/cart/domain"
	"CommerceCore/pkg/querier"
	"CommerceCore/pkg/transaction"
	"context"
	"log/slog"
)

type CartItemRepoImpl struct {
	q querier.Querier
}

func NewCartItemRepo(q querier.Querier) *CartItemRepoImpl {
	return &CartItemRepoImpl{q: q}
}

func (r *CartItemRepoImpl) Create(ctx context.Context, item *domain.CartItem) (*domain.CartItem, error) {
	err := r.q.GetContext(ctx, &item, `INSERT INTO cart_items(cart_id, product_id, quantity, price_snapshot) VALUES($1, $2, $3, $4) RETURNING id`, item.CartId, item.ProductId, item.Quantity, item.PriceSnapshot)
	if err != nil {
		slog.Error("failed insert items", "error", err)
		return nil, err
	}
	return item, nil
}

func (r *CartItemRepoImpl) Update(ctx context.Context, item *domain.CartItem) error {
	_, err := r.q.ExecContext(ctx, `UPDATE cart_items SET quantity = $1 WHERE id = $2`, item.Quantity, item.Id)
	if err != nil {
		slog.Error("failed update items", "error", err)
		return err
	}
	return nil
}

func (r *CartItemRepoImpl) Delete(ctx context.Context, id int) error {
	_, err := r.q.ExecContext(ctx, `DELETE FROM cart_items WHERE id = $1`, id)
	if err != nil {
		slog.Error("failed delete items", "error", err)
		return err
	}
	return nil
}

func (r *CartItemRepoImpl) LoadItems(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	q := r.q
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	var items []domain.CartItem
	if err := q.SelectContext(ctx, &items, `SELECT id, cart_id, product_id, quantity, price_snapshot FROM cart_items WHERE cart_id = $1`, cartID); err != nil {
		slog.Error("failed load items", "error", err)
		return nil, err
	}
	return items, nil

}
