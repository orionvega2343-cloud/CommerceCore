package domain

import "context"

type ProductRepo interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetAllProducts(ctx context.Context, isActive *bool, limit, offset int) ([]*Product, error)
	GetProductById(ctx context.Context, productId int) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, productId int) error
}
