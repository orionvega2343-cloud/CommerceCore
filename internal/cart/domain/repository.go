package domain

import "context"

type CartItemRepo interface {
	Create(ctx context.Context, item *CartItem) (*CartItem, error)
	Update(ctx context.Context, item *CartItem) error
	Delete(ctx context.Context, id int) error
}
