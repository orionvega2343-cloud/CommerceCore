package domain

import (
	"context"
	"time"
)

type ProductCache interface {
	GetProduct(ctx context.Context, id int) (*Product, error)
	SetProduct(ctx context.Context, product *Product, ttl time.Duration) error
	GetProductList(ctx context.Context, isActive *bool, limit, offset int) ([]*Product, error)
	SetProductList(ctx context.Context, isActive *bool, limit, offset int, products []*Product, ttl time.Duration) error
	DeleteProduct(ctx context.Context, id int) error
}
