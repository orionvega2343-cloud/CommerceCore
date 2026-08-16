package repository

import (
	"CommerceCore/internal/catalog/domain"
	"CommerceCore/pkg/querier"
	"context"
	"log/slog"
)

var _ domain.ProductRepo = (*ProductRepoImpl)(nil)

type ProductRepoImpl struct {
	q querier.Querier
}

func NewProductRepoImpl(q querier.Querier) *ProductRepoImpl {
	return &ProductRepoImpl{q: q}
}

func (r *ProductRepoImpl) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	err := r.q.GetContext(ctx, &product, `INSERT INTO products(name, price, stock_quantity, is_active) VALUES($1, $2, $3, $4) RETURNING id`, product.Name, product.Price, product.StockQuantity, product.IsActive)
	if err != nil {
		slog.Error("error inserting product into database", "error", err)
		return nil, err
	}
	return product, nil
}

func (r *ProductRepoImpl) GetProductById(ctx context.Context, productId int) (*domain.Product, error) {
	var p domain.Product
	err := r.q.GetContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products WHERE id = $1`, productId)
	if err != nil {
		slog.Error("error getting product from database", "error", err)
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepoImpl) GetAllProducts(ctx context.Context, isActive *bool, limit, offset int) ([]*domain.Product, error) {
	var p []*domain.Product

	if isActive == nil {
		err := r.q.SelectContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products LIMIT $1 OFFSET $2`, limit, offset)
		if err != nil {
			slog.Error("error getting products from database", "error", err)
			return nil, err
		}
		return p, nil
	}

	err := r.q.SelectContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products WHERE is_active = $1 LIMIT $2 OFFSET $3`, *isActive, limit, offset)
	if err != nil {
		slog.Error("error getting products from database", "error", err)
		return nil, err
	}
	return p, nil
}

func (r *ProductRepoImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.q.ExecContext(ctx, `UPDATE products SET price = $1, stock_quantity = $2 WHERE id = $3`, product.Price, product.StockQuantity, product.Id)
	if err != nil {
		slog.Error("error updating product from database", "error", err)
		return err
	}
	return nil
}

func (r *ProductRepoImpl) DeleteProduct(ctx context.Context, productId int) error {
	_, err := r.q.ExecContext(ctx, `UPDATE products SET is_active = false WHERE id = $1`, productId)
	if err != nil {
		slog.Error("error deleting product from database", "error", err)
		return err
	}
	return nil
}
